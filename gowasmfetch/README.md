## go-wasm-fetch

Fetches data from a url. This example demonstrates I/O handling in Go/wasm.

![image](https://user-images.githubusercontent.com/16796393/119686733-40e9c000-be64-11eb-9f17-edc8ebcb8ae5.png)

a) I/O in WebAssembly is limited by I/O capabilities of JS itself. When you do a http request from Go/wasm, it gets translated to _fetch_ call in JS. As http calls in Go are blocking in nature and WebAssembly does not allow blocking (on I/O) from JS, all such calls must be done in a goroutine. Problem comes when we have to return a response back to JS. As the function cannot block, we cannot use Go channels to communicate the response to JS. The response, in this case, must be set on the JS _Response_ object itself, which JS can then stream.

b) To enable communication between wasm and JS, we need to return a _Promise_ object instead of the actual Go function doing I/O. This is needed so we can have _resolve_ and _reject_ functions passed to our actual Go function and _Response_ streamed in JS.

```
promiseFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

  resolve := args[0]
  reject := args[1]

  // calling fetcher in a goroutine which will do http request and accordingly call resolve/reject callbacks
  go f.FetcherFn(url, resolve, reject)
  return nil
})

// return 'Promise' object
jsPromise := js.Global().Get("Promise")
return jsPromise.New(promiseFunc)
```

```
res, resErr := httpClient.Do(req)

if resErr != nil {
  reject.Invoke(newJSError(resErr))
  return
} else if res != nil && res.Body != nil && res.StatusCode == http.StatusOK {
  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)
  resolve.Invoke(newJSResponse(body))
  return
}
```

c) [Same-origin policy](https://en.wikipedia.org/wiki/Same-origin_policy) is applicable to all http requests made from wasm (same as JS). If you try to enter a url from another domain (given that this experiment is running locally on _localhost:8181_), you will get a '_Access blocked by CORS policy_' error. For the sake of this experiment, I added a _responder_ test server in the project (which listens on a different port _9191_ than the main site) with CORS handled. You need to cd to _gowasmfetch/responder_ and then *./start.sh* to start this test server. Once done, a url such as 'http://localhost:9191/fetchme' must return a successful response.

To run, cd to _gowasmfetch/site_ and then *./start.sh*
