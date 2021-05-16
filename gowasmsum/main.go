package main

import (
	wasmex "gowasmsum/wasm_exporter"
	"syscall/js"
)

func main() {
	// channel is used so we can block on it, otherwise go process will exit when this module is loaded
	// and we want the exported functions still available for later executions
	c := make(chan struct{}, 0)

	// export functions to js
	exportFuncs()

	<-c
}

/// exportFuncs exports go functions to be callable from js
func exportFuncs() {

	js.Global().Set("sumOf", js.FuncOf(wasmex.SumOf))
	js.Global().Set("sumStringOf", js.FuncOf(wasmex.SumStringOf))
	js.Global().Set("initializeWasmMemory", js.FuncOf(wasmex.InitializeWasmMemory))
	js.Global().Set("sumOf_ZeroCopy", js.FuncOf(wasmex.SumOf_ZeroCopy))
}
