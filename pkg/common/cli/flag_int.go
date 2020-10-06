package cli

import "reflect"

func (f *Flag) GetInt(ctx *Context) (int32, error) {
	parsedInt, err := f.GetInt64(ctx)
	if err != nil {
		return 0, err
	}
	return int32(parsedInt), nil
}

func (f *Flag) GetUint(ctx *Context) (uint32, error) {
	parsedInt, err := f.GetUint64(ctx)
	if err != nil {
		return 0, err
	}
	return uint32(parsedInt), nil
}

func (f *Flag) GetInt64(ctx *Context) (int64, error) {
	if !goodToExtract(f.Name, ctx) {
		return 0, ErrFlagNotReady
	}
	rVal, err := parse(ctx.mappedFlags[f.Name][0], reflect.Int64)
	if err != nil {
		return 0, err
	}
	return rVal.Interface().(int64), nil
}

func (f *Flag) GetUint64(ctx *Context) (uint64, error) {
	if !goodToExtract(f.Name, ctx) {
		return 0, ErrFlagNotReady
	}
	rVal, err := parse(ctx.mappedFlags[f.Name][0], reflect.Uint64)
	if err != nil {
		return 0, err
	}
	return rVal.Interface().(uint64), nil
}

func (f *Flag) GetIntSlice(ctx *Context) ([]int32, error) {
	int64Slice, err := f.GetInt64Slice(ctx)
	if err != nil {
		return nil, err
	}
	ret := []int32{}
	for _, i := range int64Slice {
		ret = append(ret, int32(i))
	}
	return ret, nil
}

func (f *Flag) GetUintSlice(ctx *Context) ([]uint32, error) {
	int64Slice, err := f.GetUint64Slice(ctx)
	if err != nil {
		return nil, err
	}
	ret := []uint32{}
	for _, i := range int64Slice {
		ret = append(ret, uint32(i))
	}
	return ret, nil
}

func (f *Flag) GetInt64Slice(ctx *Context) ([]int64, error) {
	if !goodToExtract(f.Name, ctx) {
		return nil, ErrFlagNotReady
	}
	rVals, err := parseMultiple(ctx.mappedFlags[f.Name], reflect.Int64)
	if err != nil {
		return nil, err
	}
	ret := []int64{}
	for _, r := range rVals {
		ret = append(ret, r.Interface().(int64))
	}
	return ret, nil
}

func (f *Flag) GetUint64Slice(ctx *Context) ([]uint64, error) {
	if !goodToExtract(f.Name, ctx) {
		return nil, ErrFlagNotReady
	}
	rVals, err := parseMultiple(ctx.mappedFlags[f.Name], reflect.Uint64)
	if err != nil {
		return nil, err
	}
	ret := []uint64{}
	for _, r := range rVals {
		ret = append(ret, r.Interface().(uint64))
	}
	return ret, nil
}
