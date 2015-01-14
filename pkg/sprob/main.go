package sprob

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ready-steady/prob"
	"github.com/ready-steady/prob/beta"
)

func Parse(line string) func(min, max float64) prob.Inverter {
	re := regexp.MustCompile("^(.+)\\((.+)\\)$")

	chunks := re.FindStringSubmatch(line)
	if chunks == nil {
		return nil
	}

	name := strings.ToLower(trim(chunks[1]))
	args := strings.Split(chunks[2], ",")

	switch strings.ToLower(name) {
	case "beta":
		if len(args) != 2 {
			return nil
		}

		α, err := strconv.ParseFloat(trim(args[0]), 64)
		if err != nil || α <= 0 {
			return nil
		}

		β, err := strconv.ParseFloat(trim(args[1]), 64)
		if err != nil || β <= 0 {
			return nil
		}

		return func(min, max float64) prob.Inverter {
			return beta.New(α, β, min, max)
		}
	}

	return nil
}

func trim(line string) string {
	return strings.Trim(line, " \t")
}
