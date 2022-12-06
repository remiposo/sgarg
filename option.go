package sgarg

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidOptName = errors.New("option's name is invalid")
)

type opt struct {
	name   optName
	values optValues
}

func newOpt(name string, values optValues) (*opt, error) {
	n, err := newOptName(name)
	if err != nil {
		return nil, err
	}
	o := &opt{
		name:   n,
		values: values,
	}
	return o, nil
}

func (o *opt) withArg() bool {
	_, ok := o.values.(*boolOptValues)
	return !ok
}

func (o *opt) appendValue(value string) error {
	return o.values.Append(value)
}

func (o *opt) abbreviatable(name string) bool {
	return o.name.isLong() && strings.HasPrefix(string(o.name), name)
}

type optName string

func newOptName(name string) (optName, error) {
	rg := regexp.MustCompile(`^[[:alnum:]]+(?:\-[[:alnum:]]+)*$`)
	if !rg.MatchString(name) {
		return "", ErrInvalidOptName
	}

	return optName(name), nil
}

func (on optName) isLong() bool {
	return len(on) != 1
}

type optValues interface {
	Append(rawValue string) error
}

type boolOptValues struct {
	values *[]bool
}

var _ optValues = &boolOptValues{}

func (bv *boolOptValues) Append(value string) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*bv.values = append(*bv.values, v)
	return nil
}

type stringOptValues struct {
	values *[]string
}

var _ optValues = &stringOptValues{}

func (sv *stringOptValues) Append(value string) error {
	*sv.values = append(*sv.values, value)
	return nil
}
