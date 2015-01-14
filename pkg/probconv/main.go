package probconv

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ready-steady/prob"
	"github.com/ready-steady/prob/beta"
)

type family uint

const (
	unknownFamily family = iota
	betaFamily
)

func ParseInverter(line string) func(min, max float64) prob.Inverter {
	name, params := parse(line)

	switch name {
	case betaFamily:
		if len(params) != 2 || params[0] <= 0 || params[1] <= 0 {
			return nil
		}

		return func(min, max float64) prob.Inverter {
			return beta.New(params[0], params[1], min, max)
		}

	default:
		return nil
	}
}

func parse(line string) (family, []float64) {
	pattern := regexp.MustCompile("^(.+)\\((.+)\\)$")

	chunks := pattern.FindStringSubmatch(line)
	if chunks == nil {
		return unknownFamily, nil
	}

	var name family

	switch strings.ToLower(trim(chunks[1])) {
	case "beta":
		name = betaFamily
	default:
		return unknownFamily, nil
	}

	chunks = strings.Split(chunks[2], ",")

	params := make([]float64, len(chunks))
	for i := range chunks {
		value, err := strconv.ParseFloat(trim(chunks[i]), 64)
		if err != nil {
			return unknownFamily, nil
		}
		params[i] = value
	}

	return name, params
}

func trim(line string) string {
	return strings.Trim(line, " \t")
}
