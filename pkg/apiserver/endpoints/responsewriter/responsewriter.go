package responsewriter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"k8s.io/apiserver/pkg/endpoints/responsewriter"
	"k8s.io/klog/v2"
)

var _ responsewriter.CloseNotifierFlusher = (*ResponseAdapter)(nil)
var _ http.ResponseWriter = (*ResponseAdapter)(nil)
var _ io.ReadCloser = (*ResponseAdapter)(nil)

// WrapHandler wraps an http.Handler to return a function that can be used as a [http.RoundTripper].
// This is used to directly connect the LoopbackConfig [http.RoundTripper]
// with the apiserver's [http.Handler], which avoids the need to start a listener
// for internal clients that use the LoopbackConfig.
// All other requests should not use this wrapper, and should be handled by the
// Grafana HTTP server to ensure that signedInUser middleware is applied.
func WrapHandler(handler http.Handler) func(req *http.Request) (*http.Response, error) {
	// ignore the lint error because the response is passed directly to the client,
	// so the client will be responsible for closing the response body.
	//nolint:bodyclose
	return func(req *http.Request) (*http.Response, error) {
		w := NewAdapter(req)
		go func() {
			handler.ServeHTTP(w, req)
			if err := w.CloseWriter(); err != nil {
				klog.Errorf("error closing writer: %v", err)
			}
		}()

		return w.Response()
	}
}

// ResponseAdapter is an implementation of [http.ResponseWriter] that allows conversion to a [http.Response].
type ResponseAdapter struct {
	req         *http.Request
	res         http.Response
	reader      io.ReadCloser
	writer      io.WriteCloser
	buffered    *bufio.ReadWriter
	ready       chan struct{}
	wroteHeader int32
}

// NewAdapter returns an initialized [ResponseAdapter].
func NewAdapter(req *http.Request) *ResponseAdapter {
	r, w := io.Pipe()
	writer := bufio.NewWriter(w)
	reader := bufio.NewReader(r)
	buffered := bufio.NewReadWriter(reader, writer)
	return &ResponseAdapter{
		req: req,
		res: http.Response{
			Proto:      req.Proto,
			ProtoMajor: req.ProtoMajor,
			ProtoMinor: req.ProtoMinor,
			Header:     make(http.Header),
		},
		reader:   r,
		writer:   w,
		buffered: buffered,
		ready:    make(chan struct{}),
	}
}

// Header implements [http.ResponseWriter].
// It returns the response headers to mutate within a handler.
func (ra *ResponseAdapter) Header() http.Header {
	return ra.res.Header
}

// Write implements [http.ResponseWriter].
func (ra *ResponseAdapter) Write(buf []byte) (int, error) {
	// via https://pkg.go.dev/net/http#ResponseWriter.Write
	// If WriteHeader is not called explicitly, the first call to Write will trigger an implicit WriteHeader(http.StatusOK).
	ra.WriteHeader(http.StatusOK)
	return ra.buffered.Write(buf)
}

// Read implements [io.Reader].
func (ra *ResponseAdapter) Read(buf []byte) (int, error) {
	return ra.buffered.Read(buf)
}

// WriteHeader implements [http.ResponseWriter].
func (ra *ResponseAdapter) WriteHeader(code int) {
	if atomic.CompareAndSwapInt32(&ra.wroteHeader, 0, 1) {
		ra.res.StatusCode = code
		ra.res.Status = fmt.Sprintf("%03d %s", code, http.StatusText(code))
		close(ra.ready)
	}
}

// FlushError implements [http.Flusher].
func (ra *ResponseAdapter) Flush() {
	// We discard io.ErrClosedPipe. This is because as we return the response as
	// soon as we have the first write or the status set, the client side with
	// the response could potentially call Close on the response body, which
	// would cause the reader side of the io.Pipe to be closed. This would cause
	// a subsequent call to Write or Flush/FlushError (that have data to write
	// to the pipe) to fail with this error. This is expected and legit, and
	// this error should be checked by the handler side by either validating the
	// error in Write or the one in FlushError. This means it is a
	// responsibility of the handler author(s) to handle this error. In other
	// cases, we log the error, as it could be potentially not easy to check
	// otherwise.
	if err := ra.FlushError(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
		klog.Error("Error flushing response buffer: ", "error", err)
	}
}

// FlushError implements an alternative Flush that returns an error. This is
// internally used in net/http and in some standard library utilities.
func (ra *ResponseAdapter) FlushError() error {
	if ra.buffered.Writer.Buffered() == 0 {
		return nil
	}

	return ra.buffered.Writer.Flush()
}

// Response returns the [http.Response] generated by the [http.Handler].
func (ra *ResponseAdapter) Response() (*http.Response, error) {
	ctx := ra.req.Context()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case <-ra.ready:
		res := ra.res
		res.Body = ra

		return &res, nil
	}
}

// Decorate implements [responsewriter.UserProvidedDecorator].
func (ra *ResponseAdapter) Unwrap() http.ResponseWriter {
	return ra
}

// CloseNotify implements [http.CloseNotifier].
func (ra *ResponseAdapter) CloseNotify() <-chan bool {
	ch := make(chan bool)
	go func() {
		<-ra.req.Context().Done()
		ch <- true
	}()
	return ch
}

// Close implements [io.Closer].
func (ra *ResponseAdapter) Close() error {
	return ra.reader.Close()
}

// CloseWriter should be called after the http.Handler has returned.
func (ra *ResponseAdapter) CloseWriter() error {
	ra.Flush()
	return ra.writer.Close()
}
