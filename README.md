# whisper-wasi
A prototype trying to get whisper work in wasi

Quick start:

1. Create a directory called testdata in the root dir.
2. Copy your model file there. You can select any file from: https://ggml.ggerganov.com
3. Run it with: `go run . -model=/name-of-file.bin -audio=/path/to/file.wav`
