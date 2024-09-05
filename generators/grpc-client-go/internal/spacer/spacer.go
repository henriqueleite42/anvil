package spacer

import (
	"errors"
)

func Space[T any](
	toSpace []T,
	getKeyValue func(T) ([]string, error),
	formatString func([]string, int) (string, error),
) ([]string, error) {

	valuesArr := [][]string{}
	for _, v := range toSpace {
		r, err := getKeyValue(v)
		if err != nil {
			return nil, err
		}
		valuesArr = append(valuesArr, r)
	}

	if len(valuesArr) == 0 {
		return nil, errors.New("empty array to space")
	}

	lenOfBiggest := 0
	for _, v := range valuesArr {
		curLen := len(v[0])
		if curLen > lenOfBiggest {
			lenOfBiggest = curLen
		}
	}

	values := []string{}
	for _, v := range valuesArr {
		r, err := formatString(v, lenOfBiggest)
		if err != nil {
			return nil, err
		}
		values = append(values, r)
	}

	return values, nil
}
