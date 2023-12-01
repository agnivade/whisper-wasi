module wasmtest

go 1.20

require (
	github.com/go-audio/wav v1.1.0
	github.com/jerbob92/wazero-emscripten-embind v1.3.1
	github.com/tetratelabs/wazero v1.5.0
)

require (
	github.com/go-audio/audio v1.0.0 // indirect
	github.com/go-audio/riff v1.0.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/jerbob92/wazero-emscripten-embind v1.3.1 => ../wazero-emscripten-embind
