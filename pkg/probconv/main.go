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
	family, params := parse(line)

	switch family {
	case betaFamily:
		return func(min, max float64) prob.Inverter {
			return beta.New(params[0], params[1], min, max)
		}
	}

	return nil
}

func parse(line string) (family, []float64) {
	pattern := regexp.MustCompile("^(.+)\\((.+)\\)$")

	chunks := pattern.FindStringSubmatch(line)
	if chunks == nil {
		return unknownFamily, nil
	}

	name := strings.ToLower(trim(chunks[1]))
	chunks = strings.Split(chunks[2], ",")
	params := make([]float64, len(chunks))
	for i := range chunks {
		value, err := strconv.ParseFloat(trim(chunks[i]), 64)
		if err != nil {
			return unknownFamily, nil
		}
		params[i] = value
	}

	switch name {
	case "beta":
		if len(params) == 2 && params[0] > 0 && params[1] > 0 {
			return betaFamily, params
		}
	}

	return unknownFamily, nil
}

func trim(line string) string {
	return strings.Trim(line, " \t")
}
