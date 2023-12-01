package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"

	"wasmtest/generated"

	embind "github.com/jerbob92/wazero-emscripten-embind"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed libmain4.wasm
var wasm []byte

func main() {
	var model, audio string
	flag.StringVar(&model, "model", "", "Path to the model file.")
	flag.StringVar(&audio, "audio", "", "Path to the audio file.")
	flag.Parse()

	log.SetFlags(log.Llongfile | log.LstdFlags)
	ctx := context.Background()
	runtimeConfig := wazero.NewRuntimeConfig()
	r := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)
	defer r.Close(ctx)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Fatal(err)
	}

	comp, err := r.CompileModule(ctx, wasm)
	if err != nil {
		log.Fatal(err)
	}

	envBuilder := r.NewHostModuleBuilder("env")

	// emscriptenExporter, err := emscripten.NewFunctionExporterForModule(comp)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// emscriptenExporter.ExportFunctions(envBuilder)

	engine := embind.CreateEngine(embind.NewConfig())
	err = engine.NewFunctionExporterForModule(comp).ExportFunctions(envBuilder)
	if err != nil {
		log.Fatal(err)
	}
	// fixEnvImports(envBuilder)

	// snpBuilder := r.NewHostModuleBuilder("wasi_snapshot_preview1")
	// wasi_snapshot_preview1.NewFunctionExporter().ExportFunctions(snpBuilder)

	// snpBuilder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3, p4, p5 int32) int32 {
	// 	log.Printf("fd_seek called-: %d %d %d %d %d", p1, p2, p3, p4, p5)
	// 	return 0
	// }).Export("fd_seek")

	// _, err = snpBuilder.Instantiate(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = envBuilder.Instantiate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_initialize").
		WithName("hellomod")

	ctx = engine.Attach(ctx)
	mod, err := r.InstantiateModule(ctx, comp, moduleConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("module name is ----", mod.Name())

	// log.Println("Exported functions---")
	// for _, v := range comp.ExportedFunctions() {
	// 	// log.Printf("%s--", k)
	// 	log.Printf("name: %q debugname: %q, ParamTypes: %v, ParamNames: %v\n", v.Name(), v.DebugName(), v.ParamTypes(), v.ParamNames())
	// }

	for _, sym := range engine.GetSymbols() {
		log.Printf("symbol: %v", sym.Symbol())
	}

	index, err := generated.Init(engine, ctx, "ggml-model-whisper-tiny.bin")

	// index, err := engine.CallPublicSymbol(ctx, "init", "ggml-model-whisper-tiny.bin")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("index returned: %v", index)

	// generated.Full_default(engine, ctx, arg0 uint32, arg1 any, arg2 string, arg3 int32, arg4 bool)

	err = generated.Free(engine, ctx, index)
	if err != nil {
		log.Fatal(err)
	}

	// _, err = engine.CallPublicSymbol(ctx, "free", index)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Imported functions---")
	// for _, fn := range comp.ImportedFunctions() {
	// 	log.Printf("name: %q debugname: %q, ParamTypes: %v, ParamNames: %v\n", fn.Name(), fn.DebugName(), fn.ParamTypes(), fn.ParamNames())
	// }
}

func fixEnvImports(builder wazero.HostModuleBuilder) {
	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3 int32) int32 {
	// 	log.Printf("_embind_register_function called-: %d %d %d %d %d %d", p1, p2, p3)
	// 	return 0
	// }).Export("_emval_get_method_caller")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3, p4, p5, p6 int32) {
	// 	log.Printf("_embind_register_function called-: %d %d %d %d %d %d", p1, p2, p3, p4, p5, p6)
	// }).Export("_embind_register_function")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1 int32) int32 {
	// 	log.Printf("__cxa_allocate_exception called-: %d", p1)
	// 	return 0
	// }).Export("__cxa_allocate_exception")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3 int32) {
	// 	log.Printf("__cxa_throw called-: %d %d %d", p1, p2, p3)
	// }).Export("__cxa_throw")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2 int32) int32 {
	// 	log.Printf("clock_gettime called-: %d %d", p1, p2)
	// 	return 0
	// }).Export("clock_gettime")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1 int32) {
	// 	log.Printf("exit called-: %d", p1)
	// }).Export("exit")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context) {
	// 	log.Printf("abort called-:")
	// }).Export("abort")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3 int32) int32 {
	// 	log.Printf("emscripten_memcpy_big called-: %d %d %d", p1, p2, p3)
	// 	return 0
	// }).Export("emscripten_memcpy_big")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3 int32) int32 {
	// 	log.Printf("__syscall_open called-: %d %d %d", p1, p2, p3)
	// 	return 0
	// }).Export("__syscall_open")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3 int32) int32 {
	// 	log.Printf("__syscall_fcntl64 called-: %d %d %d", p1, p2, p3)
	// 	return 0
	// }).Export("__syscall_fcntl64")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3 int32) int32 {
	// 	log.Printf("__syscall_ioctl called-: %d %d %d", p1, p2, p3)
	// 	return 0
	// }).Export("__syscall_ioctl")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1 int32) int32 {
	// 	log.Printf("emscripten_resize_heap called-: %d", p1)
	// 	return 0
	// }).Export("emscripten_resize_heap")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3, p4, p5 int32) int32 {
	// 	log.Printf("strftime_l called-: %d %d %d %d %d", p1, p2, p3, p4, p5)
	// 	return 0
	// }).Export("strftime_l")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1 int32) {
	// 	log.Printf("setTempRet0 called-: %d", p1)
	// }).Export("setTempRet0")

	// builder.NewFunctionBuilder().WithFunc(func(ctx context.Context, p1, p2, p3, p4, p5, p6, p7 int32) {
	// 	log.Printf("_embind_register_bigint called-: %d %d %d %d %d %d %d", p1, p2, p3, p4, p5, p6, p7)
	// }).Export("_embind_register_bigint")
}
