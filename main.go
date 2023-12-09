package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"wasmtest/generated"

	"github.com/go-audio/wav"
	embind "github.com/jerbob92/wazero-emscripten-embind"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed binary.wasm
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

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	comp, err := r.CompileModule(ctx, wasm)
	if err != nil {
		log.Fatal(err)
	}

	envBuilder := r.NewHostModuleBuilder("env")

	engine := embind.CreateEngine(embind.NewConfig())
	err = engine.NewFunctionExporterForModule(comp).ExportFunctions(envBuilder)
	if err != nil {
		log.Fatal(err)
	}

	_, err = envBuilder.Instantiate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_initialize").
		WithName("").
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithFSConfig(wazero.NewFSConfig().WithDirMount("./testdata", "/"))

	ctx = engine.Attach(ctx)
	_, err = r.InstantiateModule(ctx, comp, moduleConfig)
	if err != nil {
		log.Fatal(err)
	}

	data, err := processAudio(audio)
	if err != nil {
		log.Fatal(err)
	}

	index, err := generated.Init(engine, ctx, model)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("index returned: %v", index)

	defer func() {
		err = generated.Free(engine, ctx, index)
		if err != nil {
			log.Fatal(err)
		}
	}()

	now := time.Now()
	ret, err := generated.Full_default(engine, ctx, index, NewFloat32Array(data), "en", 1, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Processing returned: %d. Time Taken %s", ret, time.Since(now))
}

type float32Array struct {
	BackingBuffer []uint8 `embind_arg:"0"`
	ByteOffset    uint32  `embind_arg:"1"`
	Length        int32   `embind_arg:"2"`
	Constructor   *float32Array
	Buffer        []float32
}

func NewFloat32Array(in []float32) *float32Array {
	return &float32Array{
		Buffer:      in,
		Constructor: &float32Array{},
		Length:      int32(len(in)),
	}
}

func (fa *float32Array) Set(in *float32Array) {
	reslice := fa.BackingBuffer[fa.ByteOffset : fa.ByteOffset+(uint32(fa.Length)*4)]
	conv := unsafe.Slice((*uint8)(unsafe.Pointer(&in.Buffer[0])), len(in.Buffer)*4)
	copy(reslice, conv)
	fa.Buffer = in.Buffer
	fa.Length = in.Length
}

func processAudio(path string) ([]float32, error) {
	var data []float32
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	// Decode the WAV file - load the full buffer
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return nil, err
	} else if dec.SampleRate != 0x3e80 {
		return nil, fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return nil, fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}

	return data, nil
}
