package cli

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

var ErrFlagNotReady = errors.New("Flag is not presented")

func (f *Flag) GetBool(ctx *Context) bool {
	if ctx.FlagsParsed {
		_, ok := ctx.mappedFlags[f.Name]
		return ok
	}
	return false
}

func (f *Flag) GetBoolStrict(ctx *Context) (bool, error) {
	if goodToExtract(f.Name, ctx) {
		val := ctx.mappedFlags[f.Name][0]
		if val == "true" {
			return true, nil
		} else if val == "false" {
			return false, nil
		}
	}
	return false, ErrFlagNotReady
}

func (f *Flag) GetString(ctx *Context) (string, error) {
	if goodToExtract(f.Name, ctx) {
		return ctx.mappedFlags[f.Name][0], nil
	}
	return "", ErrFlagNotReady
}

func (f *Flag) GetStringSlice(ctx *Context) ([]string, error) {
	if goodToExtract(f.Name, ctx) {
		return ctx.mappedFlags[f.Name], nil
	}
	return nil, ErrFlagNotReady
}

func goodToExtract(name string, ctx *Context) bool {
	if ctx.FlagsParsed {
		if s, ok := ctx.mappedFlags[name]; ok {
			return len(s) > 0
		}
	}
	return false
}

func parse(val string, t reflect.Kind) (reflect.Value, error) {
	var err error
	var ret interface{}
	switch t {
	case reflect.Int64:
		ret, err = strconv.ParseInt(val, 10, 64)
	case reflect.Uint64:
		ret, err = strconv.ParseUint(val, 10, 64)
	case reflect.Float64:
		ret, err = strconv.ParseFloat(val, 64)
	}
	if err != nil {
		return reflect.ValueOf(nil), errors.Wrapf(err, "Failed to parse value: %s", val)
	}
	return reflect.ValueOf(ret), nil
}

func parseMultiple(vals []string, t reflect.Kind) ([]reflect.Value, error) {
	ret := []reflect.Value{}
	for _, val := range vals {
		_ret, err := parse(val, t)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse multiple values")
		}
		ret = append(ret, _ret)
	}
	return ret, nil
}
