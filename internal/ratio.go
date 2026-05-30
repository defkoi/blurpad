package internal

import (
	"errors"
	"strconv"
	"strings"
)

type ratio [2]float64

func ParseRatio(str string) (rat ratio, err error) {
	err = errors.New("Invalid ratio.")

	var (
		pNum, pDenom float64
		aNum, aDenom float64
	)

	rats := strings.Fields(str)

	parseSingleRatio := func(i int) (float64, float64, bool) {
		rat := strings.Split(rats[i], ":")
		if len(rat) != 2 {
			return 0, 0, false
		}
		num, pErr := strconv.ParseFloat(rat[0], 64)
		if pErr != nil {
			return 0, 0, false
		}
		denom, pErr := strconv.ParseFloat(rat[1], 64)
		if pErr != nil {
			return 0, 0, false
		}
		return num, denom, true
	}

	switch len(rats) {
	case 1:
		num, denom, ok := parseSingleRatio(0)
		if !ok {
			return
		}
		pNum, pDenom, aNum, aDenom =
			num, denom, num, denom
	case 2:
		_pNum, _pDenom, ok := parseSingleRatio(0)
		if !ok {
			return
		}
		_aNum, _aDenom, ok := parseSingleRatio(1)
		if !ok {
			return
		}
		pNum, pDenom, aNum, aDenom =
			_pNum, _pDenom, _aNum, _aDenom
	default:
		return
	}

	if pNum > pDenom {
		pNum, pDenom = pDenom, pNum
	}
	if aNum > aDenom {
		aNum, aDenom = aDenom, aNum
	}

	return ratio{pNum / pDenom, aNum / aDenom}, nil
}
