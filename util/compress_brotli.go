// slightly modified from https://codecov.io/gh/labstack/echo/pull/1587/diff#diff-bWlkZGxld2FyZS9jb21wcmVzc19icm90bGkuZ28=
// see https://github.com/labstack/echo/pull/1587 for code updates

package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// BrotliConfig defines the config for Brotli middleware
	BrotliConfig struct {
		// Skipper defines a function to skip middleware
		Skipper middleware.Skipper
		// Brotli compression level, optional default = -1
		Level int `yaml:"level"`
	}
	brotliResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

const (
	brotliScheme = "br"
)

// DefaultBrotliConfig is the default Brotli middleware config
var (
	DefaultBrotliConfig = BrotliConfig{
		Skipper: middleware.DefaultSkipper,
		Level:   brotli.DefaultCompression,
	}
)

// Brotli returns a middleware which compresses HTTP response using brotli compression
func Brotli() echo.MiddlewareFunc {
	return BrotliWithConfig(DefaultBrotliConfig)
}

// BrotliWithConfig return Brotli middleware with config
func BrotliWithConfig(config BrotliConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultBrotliConfig.Skipper
	}
	if config.Level == 0 {
		config.Level = DefaultBrotliConfig.Level
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			res := c.Response()
			res.Header().Add(echo.HeaderVary, echo.HeaderAcceptEncoding)
			if strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), brotliScheme) {
				res.Header().Set(echo.HeaderContentEncoding, brotliScheme) // Issue #806
				rw := res.Writer
				w := brotli.NewWriterOptions(rw, brotli.WriterOptions{Quality: config.Level})
				defer func() {
					if res.Size == 0 {
						if res.Header().Get(echo.HeaderContentEncoding) == brotliScheme {
							res.Header().Del(echo.HeaderContentEncoding)
						}
						// Reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						res.Writer = rw
						w.Reset(ioutil.Discard)
					}
					w.Close()
				}()
				grw := &brotliResponseWriter{Writer: w, ResponseWriter: rw}
				res.Writer = grw
			}
			return next(c)
		}
	}
}

func (w *brotliResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent { // Issue #489
		w.ResponseWriter.Header().Del(echo.HeaderContentEncoding)
	}
	w.Header().Del(echo.HeaderContentLength) // Issue #444
	w.ResponseWriter.WriteHeader(code)
}

func (w *brotliResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(echo.HeaderContentType) == "" {
		w.Header().Set(echo.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *brotliResponseWriter) Flush() {
	w.Writer.(*brotli.Writer).Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *brotliResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()

}
