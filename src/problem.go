package main

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ready-steady/numan/basis/linhat"
	"github.com/ready-steady/numan/grid/newcot"
	"github.com/ready-steady/numan/interp/adhier"
	"github.com/ready-steady/persim/power"
	"github.com/ready-steady/persim/system"
	"github.com/ready-steady/persim/time"
	"github.com/ready-steady/prob"
	"github.com/ready-steady/prob/gaussian"
	"github.com/ready-steady/stats/corr"
	"github.com/ready-steady/tempan/expint"
)

const (
	cacheCapacity = 1000
)

type problem struct {
	config Config

	cc uint32 // cores
	tc uint32 // tasks

	uc uint32 // dependent variables
	zc uint32 // independent variables

	marginals []prob.Inverter
	gaussian  prob.Distribution
	trans     []float64

	time  *time.List
	sched *time.Schedule

	sc uint32 // steps

	power  *power.Self
	tempan *expint.Self

	ic uint32 // inputs
	oc uint32 // outputs

	serve  func(p *problem, jobs <-chan job)
	interp *adhier.Self

	cache *cache
}

func (p *problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d, inputs: %d, outputs: %d}",
		p.cc, p.tc, p.uc, p.zc, p.ic, p.oc)
}

func newProblem(config Config) (*problem, error) {
	var err error

	p := &problem{config: config}
	c := &p.config

	plat, app, err := system.Load(c.TGFF)
	if err != nil {
		return nil, err
	}

	p.cc = uint32(len(plat.Cores))
	p.tc = uint32(len(app.Tasks))

	if len(c.CoreIndex) == 0 {
		c.CoreIndex = make([]uint16, p.cc)
		for i := uint16(0); i < uint16(p.cc); i++ {
			c.CoreIndex[i] = i
		}
	}
	if len(c.TaskIndex) == 0 {
		c.TaskIndex = make([]uint16, p.tc)
		for i := uint16(0); i < uint16(p.tc); i++ {
			c.TaskIndex[i] = i
		}
	}

	p.time = time.NewList(plat, app)
	p.sched = p.time.Compute(system.NewProfile(plat, app).Mobility)

	p.uc = uint32(len(c.TaskIndex))

	C := correlate(app, c.TaskIndex, c.ProbModel.CorrLength)
	p.trans, p.zc, err = corr.Decompose(C, p.uc, c.ProbModel.VarThreshold)
	if err != nil {
		return nil, err
	}

	p.marginals = make([]prob.Inverter, p.uc)
	marginalizer := marginalize(c.ProbModel.Marginal)
	if marginalizer == nil {
		return nil, errors.New("invalid marginal distributions")
	}
	for i, tid := range c.TaskIndex {
		delay := c.ProbModel.MaxDelay * plat.Cores[p.sched.Mapping[tid]].Time[app.Tasks[tid].Type]
		p.marginals[i] = marginalizer(delay)
	}

	p.gaussian = gaussian.New(0, 1)

	switch c.Target {
	case "end-to-end-delay":
		p.ic = p.zc
		p.oc = 1

		p.serve = serveEndToEndDelay

	case "temperature-profile":
		p.sc = uint32(p.sched.Span / c.TempAnalysis.TimeStep)

		p.power = power.New(plat, app, c.TempAnalysis.TimeStep)
		p.tempan, err = expint.New(expint.Config(c.TempAnalysis))
		if err != nil {
			return nil, err
		}

		p.ic = 1 + p.zc // +1 for time
		p.oc = uint32(len(c.CoreIndex))

		p.serve = serveTemperatureProfile

	default:
		return nil, errors.New("the target is unknown")
	}

	p.interp = adhier.New(newcot.NewOpen(uint16(p.ic)), linhat.NewOpen(uint16(p.ic)),
		adhier.Config(c.Interpolation), uint16(p.oc))

	p.cache = newCache(p.zc, cacheCapacity)

	return p, nil
}

func (p *problem) solve() *adhier.Surrogate {
	ic, oc, Δ := p.ic, p.oc, p.ic-p.zc // Δ for time
	cache := p.cache

	jobs := p.spawnWorkers()

	NC, EC := uint32(0), uint32(0)

	if p.config.Verbose {
		fmt.Printf("%12s %12s (%6s) %12s %12s (%6s)\n",
			"new nodes", "new evals", "%", "total nodes", "total evals", "%")
	}

	surrogate := p.interp.Compute(func(nodes []float64, index []uint64) []float64 {
		nc, ec := uint32(len(nodes))/ic, uint32(0)

		NC += nc
		if p.config.Verbose {
			fmt.Printf("%12d", nc)
		}

		done := make(chan result, nc)
		values := make([]float64, oc*nc)

		for i := uint32(0); i < nc; i++ {
			key := cache.key(index[Δ+i*ic:])

			data := cache.get(key)
			if data == nil {
				ec++
			}

			jobs <- job{
				key:   key,
				data:  data,
				node:  nodes[i*ic:],
				value: values[i*oc:],
				done:  done,
			}
		}

		for i := uint32(0); i < nc; i++ {
			result := <-done
			cache.set(result.key, result.data)
		}

		EC += ec
		if p.config.Verbose {
			fmt.Printf(" %12d (%6.2f) %12d %12d (%6.2f)\n",
				ec, float64(ec)/float64(nc)*100,
				NC, EC, float64(EC)/float64(NC)*100)
		}

		return values
	})

	close(jobs)

	return surrogate
}

func (p *problem) compute(nodes []float64) []float64 {
	ic, oc := p.ic, p.oc

	jobs := p.spawnWorkers()

	nc := uint32(len(nodes)) / ic

	done := make(chan result, nc)
	values := make([]float64, p.oc*nc)

	jc, rc := uint32(0), uint32(0)
	nextJob := job{
		node:  nodes[jc*ic:],
		value: values[jc*oc:],
		done:  done,
	}

	for jc < nc || rc < nc {
		select {
		case jobs <- nextJob:
			jc++

			if jc >= nc {
				close(jobs)
				jobs = nil
				continue
			}

			nextJob = job{
				node:  nodes[jc*ic:],
				value: values[jc*oc:],
				done:  done,
			}
		case <-done:
			rc++
		}
	}

	return values
}

func (p *problem) evaluate(s *adhier.Surrogate, points []float64) []float64 {
	return p.interp.Evaluate(s, points)
}

func (p *problem) spawnWorkers() chan job {
	wc := int(p.config.Workers)
	if wc <= 0 {
		wc = runtime.NumCPU()
	}

	if p.config.Verbose {
		fmt.Printf("Using %d workers...\n", wc)
	}

	runtime.GOMAXPROCS(wc)

	jobs := make(chan job)
	for i := 0; i < wc; i++ {
		go p.serve(p, jobs)
	}

	return jobs
}
