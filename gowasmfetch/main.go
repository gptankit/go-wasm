package main

import (
	"gowasmfetch/fetcher"
	"syscall/js"
)

func main() {
	// channel is used so we can block on it, otherwise go process will exit when this module is loaded
	// and we want the exported functions still available for later executions
	c := make(chan struct{}, 0)

	wasmExportsMap := getWasmExportables()

	// export functions to js
	exportFuncs(wasmExportsMap)

	<-c
}

func getWasmExportables() map[string]js.Func {

	return map[string]js.Func{
		"fetchme": fetcher.NewFetcher().Bind(),
		// add more funcs here
	}
}

/// exportFuncs exports go/wasm functions so as to be callable from js
func exportFuncs(wasmExportsMap map[string]js.Func) {

	for k, v := range wasmExportsMap {
		js.Global().Set(k, v) // set function definition on js 'window' object
	}
}
