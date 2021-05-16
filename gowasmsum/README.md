## go-wasm-sum


Calculates sum of an array of values. This example demonstrates three ways to share data between JS and WebAssembly.

![image](https://user-images.githubusercontent.com/16796393/118147990-c1a0c900-b42d-11eb-9095-93a767be392d.png)

**Sum1: Considering value as uint8 and using in-built function for copying memory.**

Create a new go array of the same len as js array - 

```
goArray := make([]uint8, jsArrayLen)
```

And then copy using `js.CopyBytesToGo()` - 

```
js.CopyBytesToGo(goArray, jsArray)
```

**Sum2: Considering value as string and manually copying memory.**

Create a new go array of the same len as js array - 

```
goArray := make([]string, jsArrayLen)
```

And then copy contents by looping over js array - 

```
for i := 0; i < jsArrayLen; i++ {
  goArray[i] = jsArray.Index(i).String()
}
```

**Sum3: Considering value as uint8 and using pointer to share and populate memory.**

Have two functions - one for returning a pointer to linear memory (that JS can populate data with) and another to read the populated data from that pointer location.

Crucial thing to note while returning pointer is that Go provides no mapping for a pointer to a JS object, so I converted it to an _unsafe.Pointer_ and set it in a map (so Go can convert the `map[string]interface{}` to a JS _Object_ with the pointer value casted to a _uintptr_).

```
goArray := make([]uint8, goArrayLen)
ptr = &goArray

boxedPtr := unsafe.Pointer(ptr)
boxedPtrMap := map[string]interface{}{
  "internalptr": boxedPtr,
}
return js.ValueOf(boxedPtrMap)
```

Another thing to note while receiving pointer is that the pointer value - _uintptr_ - is just a number with no pointer semantics, so I converted it to a _*[]uint8_  which can then be used to iterate over.

```
sliceHeader := &reflect.SliceHeader{
  Data: uintptr(args[0].Int()),
  Len:  len,
  Cap:  len,
}

var ptr = (*[]uint8)(unsafe.Pointer(sliceHeader))
```

To run, cd to _gowasmsum/site_ and then *go run server.go*
