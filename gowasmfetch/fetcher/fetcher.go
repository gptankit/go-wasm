package fetcher

import (
	gofetch "gowasmfetch/fetcher/gofetch"
	"syscall/js"
)

type GenericFetcher interface {
	Bind() js.Func
}

type NormalFetcher struct {
	FetcherFn func(string, js.Value, js.Value)
}

// NewFetcher returns fether object
func NewFetcher() GenericFetcher {

	fetcher := &NormalFetcher{
		FetcherFn: gofetch.InitHttpReq(),
	}

	return fetcher
}

// Bind returns a JS promise wrapper that can be set on JS window
func (f *NormalFetcher) Bind() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := args[0].String()
		// return a promise that can be called from JS to fetch url
		return newJSPromise(f, url)
	})
}

// newJSPromise returns a new JS promise
func newJSPromise(f *NormalFetcher, url string) js.Value {

	promiseFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		resolve := args[0]
		reject := args[1]

		// calling fetcher in a goroutine which will do http request and accordingly call resolve/reject callbacks
		go f.FetcherFn(url, resolve, reject)
		return nil
	})

	// return Promise object
	jsPromise := js.Global().Get("Promise")
	return jsPromise.New(promiseFunc)
}
