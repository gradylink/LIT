package main

import (
	"strconv"

	"github.com/dave/jennifer/jen"
)

func ParseInput(i Input, toNumber bool) (jen.Code, error) {
	if i.Type == 1 {
		switch (*i.Value).(type) {
		case float64:
			return jen.Lit(*i.Value), nil
		case string:
			if toNumber {
				num, err := strconv.ParseFloat((*i.Value).(string), 64)
				if err != nil {
					return nil, err
				}

				return jen.Lit(num), nil
			}
			return jen.Lit(*i.Value), nil
		}
	}
	if i.BroadcastID != nil {
		return jen.Lit(*i.BroadcastID), nil
	}

	return nil, nil
}
