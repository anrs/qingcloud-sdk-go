package conn

import (
	"errors"
	"fmt"
	"reflect"
)

type Keys []string

type Condition struct {
	Required Keys
	Integers Keys
	Lists    Keys
}

func (c Condition) CheckIntegers(params *Dict) error {
	if len(c.Integers) < 1 {
		return nil
	}

	for _, k := range c.Integers {
		v, ok := (*params)[k]
		if !ok {
			continue
		}

		if _, ok := v.(int); !ok {
			return errors.New(fmt.Sprintf("%s value is not int", k))
		}
	}

	return nil
}

func (c Condition) CheckRequired(params *Dict) error {
	if len(c.Required) < 1 {
		return nil
	}

	for _, k := range c.Required {
		if _, ok := (*params)[k]; !ok {
			return errors.New(fmt.Sprintf("%s not found", k))
		}
	}

	return nil
}

func (c Condition) CheckLists(params *Dict) error {
	if len(c.Lists) < 1 {
		return nil
	}

	for _, k := range c.Lists {
		v, ok := (*params)[k]
		if !ok {
			continue
		}

		kind := reflect.ValueOf(v).Kind()
		if kind != reflect.Slice && kind != reflect.Array {
			return errors.New(fmt.Sprintf("%s value is not list", k))
		}
	}

	return nil
}

func (c Condition) Check(params *Dict) error {
	if err := c.CheckRequired(params); err != nil {
		return err
	}

	if err := c.CheckIntegers(params); err != nil {
		return err
	}

	if err := c.CheckLists(params); err != nil {
		return err
	}

	return nil
}
