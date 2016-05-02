package distribution

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/ready-steady/probability/distribution"
)

type family uint

const (
	unknownFamily family = iota
	betaFamily
	uniformFamily
)

func Parse(line string) (func(float64, float64) distribution.Continuous, error) {
	family, params := parse(line)

	switch family {
	case betaFamily:
		return func(min, max float64) distribution.Continuous {
			return distribution.NewBeta(params[0], params[1], min, max)
		}, nil
	case uniformFamily:
		return func(min, max float64) distribution.Continuous {
			return distribution.NewUniform(min, max)
		}, nil
	default:
		return nil, errors.New("the marginal distribution is unknown")
	}
}

func parse(line string) (family, []float64) {
	pattern := regexp.MustCompile("^(.+)\\((.*)\\)$")

	match := pattern.FindStringSubmatch(line)
	if match == nil {
		return unknownFamily, nil
	}

	name, rest := strings.ToLower(trim(match[1])), trim(match[2])

	params := make([]float64, 0)
	if len(rest) > 0 {
		chunks := strings.Split(rest, ",")
		for i := range chunks {
			if value, err := strconv.ParseFloat(trim(chunks[i]), 64); err == nil {
				params = append(params, value)
			} else {
				return unknownFamily, nil
			}
		}
	}

	switch name {
	case "beta":
		if len(params) == 2 && params[0] > 0.0 && params[1] > 0.0 {
			return betaFamily, params
		}
	case "uniform":
		if len(params) == 0 {
			return uniformFamily, params
		}
	}

	return unknownFamily, nil
}

func trim(line string) string {
	return strings.Trim(line, " \t")
}
