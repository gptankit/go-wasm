package gofetch

import (
	"errors"
	"io/ioutil"
	"net/http"
	"syscall/js"
	"time"
)

// InitHttpReq sets up a js func that can fetch from a url and accordingly call resolve/reject callbacks
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
		} else if res != nil && res.Body != nil && res.StatusCode == http.StatusOK {
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			resolve.Invoke(newJSResponse(body))
			return
		}

		reject.Invoke(newJSError(errors.New("Something went wrong!")))
	}
}

// newJSError returns a new js error object
func newJSError(err error) interface{} {

	// setting error to the error object
	jsError := js.Global().Get("Error")
	return jsError.New(err.Error())
}

// newJSResponse returns a new js response object with response body set
func newJSResponse(body []byte) interface{} {

	// converting "data" to a js Uint8Array object
	jsArray := js.Global().Get("Uint8Array")
	jsBody := jsArray.New(len(body))
	js.CopyBytesToJS(jsBody, body)

	// setting js array to the response object
	jsResponse := js.Global().Get("Response")
	return jsResponse.New(jsBody)
}
