package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/statistics/metric"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(problem *internal.Problem, fi *mat.File, fo *mat.File) error {
	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	interpolator, err := internal.NewInterpolator(problem, target)
	if err != nil {
		return err
	}

	problem.Println(problem)
	problem.Println(target)

	surrogate := new(adhier.Surrogate)
	if fi == nil {
		return errors.New("an input file is required")
	}
	if err := fi.Get("surrogate", surrogate); err != nil {
		return err
	}

	problem.Println(surrogate)

	sc := problem.Config.Samples
	if sc == 0 {
		return errors.New("the number of samples is zero")
	}

	ic, oc := target.InputsOutputs()

	if problem.Config.Seed > 0 {
		rand.Seed(problem.Config.Seed)
	} else {
		rand.Seed(time.Now().Unix())
	}
	points := probability.Sample(uniform.New(0, 1), sc*ic)

	problem.Println("Evaluating the original model...")
	values := invoke(target, points)

	problem.Println("Evaluating the surrogate model...")
	approximations := interpolator.Evaluate(surrogate, points)

	fmt.Printf("Inputs: %d, Outputs: %d, Nodes: %d, NRMSE: %.2e\n",
		surrogate.Inputs, surrogate.Outputs, surrogate.Nodes,
		metric.NRMSE(approximations, values))

	if fo == nil {
		return nil
	}

	if err := fo.PutMatrix("points", points, ic, sc); err != nil {
		return err
	}
	if err := fo.PutMatrix("values", values, oc, sc); err != nil {
		return err
	}
	if err := fo.PutMatrix("approximations", approximations, oc, sc); err != nil {
		return err
	}

	return nil
}

func invoke(target internal.Target, points []float64) []float64 {
	wc := uint32(runtime.GOMAXPROCS(0))
	ic, oc := target.InputsOutputs()
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	jobs := make(chan uint32, pc)
	done := make(chan bool, pc)

	for i := uint32(0); i < wc; i++ {
		go func() {
			for j := range jobs {
				target.Evaluate(points[j*ic:(j+1)*ic], values[j*oc:(j+1)*oc], nil)
				done <- true
			}
		}()
	}

	for i := uint32(0); i < pc; i++ {
		jobs <- i
	}
	for i := uint32(0); i < pc; i++ {
		<-done
	}

	close(jobs)

	return values
}
