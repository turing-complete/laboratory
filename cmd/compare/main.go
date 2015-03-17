package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/statistics"
	"github.com/ready-steady/statistics/test"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, predict *mat.File, observe *mat.File) error {
	if predict == nil || observe == nil {
		return errors.New("two data files are required")
	}

	solution := new(internal.Solution)
	if err := predict.Get("solution", solution); err != nil {
		return err
	}

	observations := []float64{}
	if err := observe.Get("values", &observations); err != nil {
		return err
	}

	predictions := []float64{}
	if err := predict.Get("values", &predictions); err != nil {
		return err
	}

	ns := int(config.Assessment.Samples)
	no := len(observations) / ns
	nm := no / 2

	cut := func(data []float64, i int) []float64 {
		piece := make([]float64, ns)
		for j := 0; j < ns; j++ {
			piece[j] = data[j*no+i]
		}
		return piece
	}

	fmt.Println(solution)

	μo := make([]float64, nm)
	vo := make([]float64, nm)

	μp := make([]float64, nm)
	vp := make([]float64, nm)

	εμ := make([]float64, nm)
	εv := make([]float64, nm)
	εp := make([]float64, nm)

	analytic := len(solution.Expectation) == no

	// Compute errors across all outputs.
	for i := 0; i < nm; i++ {
		j := i * 2

		observations := cut(observations, j)
		predictions := cut(predictions, j)

		μo[i] = statistics.Mean(observations)
		vo[i] = statistics.Variance(observations)

		if analytic {
			μp[i] = solution.Expectation[j]
			vp[i] = solution.Expectation[j+1] - μp[i]*μp[i]
		} else {
			μp[i] = statistics.Mean(predictions)
			vp[i] = statistics.Variance(predictions)
		}

		εμ[i] = math.Abs(μo[i] - μp[i])
		εv[i] = math.Abs(vo[i] - vp[i])

		_, _, εp[i] = test.KolmogorovSmirnov(observations, predictions, 0)
	}

	if nm == 1 {
		fmt.Printf("Result: μ %.2e ± %.2e (%.2e), v %.2e ± %.2e (%.2e), p %.2e\n",
			μo[0], εμ[0], εμ[0]/μo[0], vo[0], εv[0], εv[0]/vo[0], εp[0])
		return nil
	}

	if config.Verbose {
		for i := 0; i < nm; i++ {
			fmt.Printf("%7d: μ %.2e ± %.2e (%.2e), v %.2e ± %.2e (%.2e), p %.2e\n",
				i, μo[i], εμ[i], εμ[i]/μo[i], vo[i], εv[i], εv[i]/vo[i], εp[i])
		}
	}

	μμo, μεμ := statistics.Mean(μo), statistics.Mean(εμ)
	μvo, μεv := statistics.Mean(vo), statistics.Mean(εv)
	μεp := statistics.Mean(εp)

	fmt.Printf("Average: μ %.2e ± %.2e (%.2e), v %.2e ± %.2e (%.2e), p %.2e\n",
		μμo, μεμ, μεμ/μμo, μvo, μεv, μεv/μvo, μεp)

	mεμ, kμ := max(εμ)
	mεv, kv := max(εv)
	mεp, _ := max(εp)

	fmt.Printf("Maximal: μ %.2e ± %.2e (%.2e), v %.2e ± %.2e (%.2e), p %.2e\n",
		μo[kμ], mεμ, mεμ/μo[kμ], vo[kv], mεv, mεv/vo[kv], mεp)

	return nil
}

func max(data []float64) (float64, int) {
	value, k := math.Inf(-1), -1

	for i, x := range data {
		if x > value {
			value, k = x, i
		}
	}

	return value, k
}
