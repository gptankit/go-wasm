## go-wasm-eval

Evaluates expressions. This example demonstrates usage of external Go packages and how to structure Go/wasm code.

![image](https://user-images.githubusercontent.com/16796393/118148639-64594780-b42e-11eb-8271-3d0304232ca0.png)

External Go packages can be imported and used in the same way as you would do for any normal Go program. Since we can only export functions of the type _js.Func_, one should ensure that the input/return parameters are loosely coupled with the core Go functionality (provided by either external package or layers above). Decoupling the two will enable us to extend/change both modules easily in the light of future contract changes (Go/wasm or external package).

```
func (e *NormalEvaluator) Bind() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		expression := args[0].String()
		result, err := e.EvaluateFn(expression)
		if isErr(err) {
			return false
		}
		return js.ValueOf(result)
	})
}
```

To run, cd to _gowasmeval/site_ and then *go run server.go*