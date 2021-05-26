package wasm_exporter

import (
	"reflect"
	"syscall/js"
	"unsafe"
)

// SumOf returns sum of uint8 numbers
func SumOf(this js.Value, args []js.Value) interface{} {

	jsArray := args[0]
	jsArrayLen := jsArray.Get("byteLength").Int() // we can also use jsArray.Length() provided by syscall/js

	// one way of copying memory (given array content is uin8)
	// if array content is not uint8 but uint16 or uint32 or uint64, consider using encoding/binary package to flatten array to []byte and then encode to uint8
	goArray := make([]uint8, jsArrayLen)

	js.CopyBytesToGo(goArray, jsArray) // array value types have to be uint8

	var sum uint8
	for i := 0; i < len(goArray); i++ {
		sum += goArray[i]
	}

	return js.ValueOf(sum)
}

// SumStringOf returns concatenation of strings
func SumStringOf(this js.Value, args []js.Value) interface{} {

	jsArray := args[0]
	jsArrayLen := jsArray.Get("length").Int()

	// one way of copying memory (given array content is string)
	goArray := make([]string, jsArrayLen)

	for i := 0; i < jsArrayLen; i++ {
		goArray[i] = jsArray.Index(i).String()
	}

	var sum string
	for i := 0; i < len(goArray); i++ {
		sum += goArray[i]
	}

	return js.ValueOf(sum)
}

// Below functions are used as a zero copy alternative to SumOf function

// InitializeWasmMemory initializes wasm memory of passed length and returns a pointer
func InitializeWasmMemory(this js.Value, args []js.Value) interface{} {

	var ptr *[]uint8
	goArrayLen := args[0].Int()

	goArray := make([]uint8, goArrayLen)
	ptr = &goArray

	boxedPtr := unsafe.Pointer(ptr)
	boxedPtrMap := map[string]interface{}{
		"internalptr": boxedPtr,
	}
	return js.ValueOf(boxedPtrMap)
}

// SumOf_ZeroCopy loads the array populated at the pointer and returns sum
func SumOf_ZeroCopy(this js.Value, args []js.Value) interface{} {

	var len = args[1].Int()

	sliceHeader := &reflect.SliceHeader{
		Data: uintptr(args[0].Int()),
		Len:  len,
		Cap:  len,
	}

	var ptr = (*[]uint8)(unsafe.Pointer(sliceHeader))

	var sum uint8
	for i := 0; i < len; i++ {
		sum += uint8((*ptr)[i])
	}

	return js.ValueOf(sum)
}
