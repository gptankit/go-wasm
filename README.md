# Understanding WebAssembly and interoperability with Go and Javascript

_(This is a gentle introduction to the world of WebAssembly from the lens of Go and Javascript. The purpose of this primer is to introduce WebAssembly to people who are already familiar with Go and want to use their understanding to build fast programs for the web and other environments outside the world of Go)_

Traditionally Javascript has been the language of choice for doing any kind of logic on browsers. Because of its interpreted nature, it is usually slow in execution. Not to mention that writing complex logics in Javascript quickly becomes cumbersome and bloated. A new system was needed to allow for complex and faster code execution in the browser-land.

## Introducing WebAssembly

WebAssembly (or wasm) is a new binary code format capable of executing in browser (and some non-browser) environments. It - 

- Is portable, fast and safe
- Is primarily made for running in the browser alongside Javascript (same VM can execute both JS and Wasm)
- Can also be executed by other runtimes like _Node.js_ and wasm specific runtimes (like [wasmer](https://wasmer.io/) and [wasmtime](https://wasmtime.dev/))
- Runs in a sandbox thus providing no direct access to host environment
- Has a binary representation (_.wasm_) and a text representation (_.wat_ - uses [S-expressions](https://en.wikipedia.org/wiki/S-expression))
- Is a compilation target for C, C++, Rust, Go and many other languages

Most languages use [emscripten](https://emscripten.org/) to compile to wasm (like C, C++, Rust, but not Go). Most languages also consider wasm module as a library but Go considers it as an application running alonside Javascript.

## Grammar

WebAssembly has two representational formats - **binary** and **text**. Both of them are generated from a common abstract syntax. This abstract syntax follows from the same grammar rules but is represented differently in the final output. What it means is that the binary output (_wasm_) is a _linear encoding_ of abstract syntax to be executed by the virtual machine, while the text output (_wat_) is _S-expressions_ rendering of abstract syntax useful for analyzing and understanding the underlying logic.

## Internals

#### Types

WebAssembly supports **i32**, **i64**, **f32** and **f64** fundamental types only (also called *numtypes*). So, it is the responsibility of the glue code (or embedder) to convert other fundamental/composite types to _numtypes_ so wasm can understand them.

Apart from these, wasm internally maintains other non-fundamental types - *reference types*, *value types*, *result types*, *function types*, *limits*, *memory types*, *table types*, *global types* and *external types* which are used in the wasm execution life cycle.

#### Values

While wasm operates only on numeric values (integers and floating points), the program can represent other entities in terms of *bytes* and *names* (strings) as well. These are then converted to integers while executing.

#### Memories

Wasm memory is *linear* - just a contiguous block of untyped bytes. This memory is exported and both wasm and Javascript can share data using this memory space.

#### Tables

Wasm table is a vector of reference types - references to functions, global or memory addresses etc.

#### Module

Wasm module is the result of final compilation of a wasm binary by the browser. It can define a set of exports and imports to interact with external environments.

#### Instance

A wasm instance is a runtime representation of a wasm module. It is treated as running in a sandbox environment seperate from host environment and other wasm instances and can only interact with them via well defined APIs.

## Inspecting

Inspecting wasm by hand is hard as it is a binary format. Primary tool that we can use for inspecting its internal structure is [wabt](https://github.com/WebAssembly/wabt). It can be used to do a multitude of things with wasm including conversion to text format (and vice versa), printing info on the binary, validation etc. Most important programs included in *wabt* project are - 

- **wat2wasm**: translate from WebAssembly text format to the WebAssembly binary format
- **wasm2wat**: the inverse of wat2wasm, translate from the binary format back to the text format (also known as a .wat)
- **wasm-objdump**: print information about a wasm binary. Similiar to objdump.
- **wasm-strip**: remove sections of a WebAssembly binary file
- **wasm-validate**: validate a file in the WebAssembly binary format

## Go and WebAssembly

WebAssembly is primarily built to be a compilation target to most high level programming languages (you can also code in WebAssembly text format and have it compiled to the binary though its not recommended). Below we'll limit ourselves to generating wasm binary from Go.

a) In order for Javascript to call into Go/wasm, the Go functions must follow below signature - 

`fn func(this js.Value, args []js.Value) interface{}`

here, *this* represents a Javascript context and *args* array represent the parameters to be passed to this function. The return value can be any value understood by both Go and Javascript lifted to an interface{} as the return type.

This *fn* is to be wrapped in a *js.Func* object by passing fn to a *js.FuncOf* call - 

`fn_exported js.FuncOf(fn) js.Func`

This *fn_exported* is then set on the Javascript *'window'* object using [syscall/js](https://golang.org/pkg/syscall/js/) package - 

`js.Global().Set("fn_exported_name", fn_exported)`

b) Wasm cannot access DOM directly, therefore Go provides the package *syscall/js* which can be used to interact with DOM. Functions like *js.Global()* return Javascript *'window'* object that Go can use to get or set DOM elements.

c) Wasm allows for single return value only. The allowed return types are *js.Value*, *js.Func*, *nil*, *bool*, *integers*, *floats*, *strings*, *[]interface{}* and *map[string]interface{}* which are then converted to corresponding Javascript types while crossing function boundaries.

d) Go code can be compiled to wasm by setting **GOOS=js** and **GOARCH=wasm** directives in *go build* -

`GOOS=js GOARCH=wasm go build -o $(BINPATH)`

The generated *.wasm* file must be instantiated in the Javascript so the browser can compile it and make it ready for use. Go provides the instantiation mechanism through a special script called *wasm_exec.js* (found in *$GOROOT/misc/wasm/wasm_exec.js* location) which must also be copied to the web server and included in the html file. Apart from this, *wasm_exec.js* also enables mechanisms through which Go/wasm code can interact with Javascript APIs.

## Usage from Javascript

a) One virtual machine is capable of executing both Javascript and wasm binary\
b) VM is responsible to compile wasm binary (either AOT and/or JIT) depending on the host environment\
c) Javascript has access to exported functions and memories of a wasm binary while wasm has access to the DOM via Javascript APIs\
d) Data can be shared between Javascript and wasm environment through function calls or linear memory but given that they differ in the data types, some glue code is usually needed to convert values from one environment to another.

#### Instantiating wasm 

In order to instantiate a *.wasm* file, the Javascript file must import *wasm_exec.js* script.

`<script src="wasm_exec.js"></script>`

Then, we can use a little glue code provided by _wasm_exec.js_ to instantiate the wasm file (say _main.wasm_).

```
(async function loadAndRunGoWasm() {
    const go = new Go(); 
    const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
    go.run(result.instance)
})()
```

The _WebAssembly.instantiateStreaming_ call accepts two arguments - first, a _response_ object which contains the *.wasm* file fetched from server, and second, an _importObject_ which we get from the go runtime in `const go = new Go()` call. The *instantiateStreaming* call returns a result whose instance property can be passed to `go.run()` to start this wasm instance.

#### Calling into wasm

As Go exports functions to the Javascript _window_ object, we can call wasm functions normally as we would do for any Javascript functions - by using the function name `fn_exported_name` set in Go exports and passing in the required parameters.

## Working examples

- **gowasmsum** (example demonstrating use of memory copying techniques between wasm and Javascript): https://github.com/gptankit/go-wasm/tree/main/gowasmsum
- **gowasmeval** (example demonstrating use of external Go packages and suggested code structure): https://github.com/gptankit/go-wasm/tree/main/gowasmeval
