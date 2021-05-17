package gofetch

import (
	"io/ioutil"
	"net/http"
	"syscall/js"
	"time"
)

// InitHttpReq sets up a JS func that can fetch from a url and accordingly call resolve/reject callbacks
func InitHttpReq() func(string, js.Value, js.Value) {

	return func(url string, resolve js.Value, reject js.Value) {

		httpClient := &http.Client{
			Timeout: 2 * time.Second,
		}

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			reject.Invoke(newJSError(reqErr))
			return
		}

		res, resErr := httpClient.Do(req)
		if resErr != nil {
			reject.Invoke(newJSError(resErr))
			return
		}

		if res != nil && res.Body != nil && res.StatusCode == http.StatusOK {
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			resolve.Invoke(newJSResponse(body))
			return
		}
	}
}

// newJSError returns a new JS error object
func newJSError(err error) interface{} {

	jsError := js.Global().Get("Error")
	return jsError.New(err.Error())
}

// newJSResponse returns a new JS response object with response body set
func newJSResponse(body []byte) interface{} {

	// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
	jsArray := js.Global().Get("Uint8Array")
	jsBody := jsArray.New(len(body))
	js.CopyBytesToJS(jsBody, body)

	jsResponse := js.Global().Get("Response")
	return jsResponse.New(jsBody)
}
