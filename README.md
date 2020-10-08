compile with

	GOOS=js GOARCH=wasm go build -o main.wasm 

Serve the directory (at least index.html, main.wasm, and wasm_exec.js)

then run chrome or another browser that supports wasm

Test with something like:

	regextimestamp('Gravwell Webassembly (?P<timestamp>.+)', 'timestamp', 'An entry with some timestamp in it via Gravwell Webassembly Oct 7 18:43:00')

In the console.


