// Code generated by wazero-emscripten-embind, DO NOT EDIT.
package generated

import (
	"context"

	"github.com/jerbob92/wazero-emscripten-embind"
)

func Free(e embind.Engine, ctx context.Context, arg0 uint32) error {
	_, err := e.CallPublicSymbol(ctx, "free", arg0)
	return err
}

func Full_default(e embind.Engine, ctx context.Context, arg0 uint32, arg1 any, arg2 string, arg3 int32, arg4 bool) (int32, error) {
	res, err := e.CallPublicSymbol(ctx, "full_default", arg0, arg1, arg2, arg3, arg4)
	if err != nil {
		return int32(0), err
	}
	if res == nil {
		return int32(0), nil
	}
	return res.(int32), nil
}

func Init(e embind.Engine, ctx context.Context, arg0 string) (uint32, error) {
	res, err := e.CallPublicSymbol(ctx, "init", arg0)
	if err != nil {
		return uint32(0), err
	}
	if res == nil {
		return uint32(0), nil
	}
	return res.(uint32), nil
}