package jshttp

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// ToResponse converts JavaScript sides Response to *http.Response.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func ToResponse(res js.Value) (*http.Response, error) {
	status := res.Get("status").Int()
	promise := res.Call("text")
	body, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}
	header := ToHeader(res.Get("headers"))
	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)

	return &http.Response{
		Status:        strconv.Itoa(status) + " " + res.Get("statusText").String(),
		StatusCode:    status,
		Header:        header,
		Body:          io.NopCloser(strings.NewReader(body.String())),
		ContentLength: contentLength,
	}, nil
}

// ToJSResponse converts *http.Response to JavaScript sides Response.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func ToJSResponse(w *ResponseWriterBuffer) (js.Value, error) {
	<-w.ReadyCh // wait until ready
	status := w.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	respInit := jsutil.NewObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(w.Header()))
	respInit.Set("webSocket", w.WebSocket)
	if status == 101 || status == 204 || status == 205 || status == 304 {
		return jsutil.ResponseClass.New(respInit), nil
	}
	readableStream := jsutil.ConvertReaderToReadableStream(w.Reader)
	return jsutil.ResponseClass.New(readableStream, respInit), nil
}
