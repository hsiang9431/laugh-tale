package cli

import "reflect"

func (f *Flag) GetFloat32(ctx *Context) (float32, error) {
	parsed, err := f.GetFloat64(ctx)
	if err != nil {
		return 0, err
	}
	return float32(parsed), nil
}

func (f *Flag) GetFloat64(ctx *Context) (float64, error) {
	if !goodToExtract(f.Name, ctx) {
		return 0, ErrFlagNotReady
	}
	rVal, err := parse(ctx.mappedFlags[f.Name][0], reflect.Float64)
	if err != nil {
		return 0, err
	}
	return rVal.Interface().(float64), nil
}

func (f *Flag) GetFloat32Slice(ctx *Context) ([]float32, error) {
	float64Slice, err := f.GetFloat64Slice(ctx)
	if err != nil {
		return nil, err
	}
	ret := []float32{}
	for _, i := range float64Slice {
		ret = append(ret, float32(i))
	}
	return ret, nil
}

func (f *Flag) GetFloat64Slice(ctx *Context) ([]float64, error) {
	if !goodToExtract(f.Name, ctx) {
		return nil, ErrFlagNotReady
	}
	rVals, err := parseMultiple(ctx.mappedFlags[f.Name], reflect.Float64)
	if err != nil {
		return nil, err
	}
	ret := []float64{}
	for _, r := range rVals {
		ret = append(ret, r.Interface().(float64))
	}
	return ret, nil
}
