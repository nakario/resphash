package resphash

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HashWriter adds the 'Resp-Hash' HTTP header to investigate the distribution
// of the response body. It implements `http.ResponseWriter` interface.
type HashWriter struct {
	w          http.ResponseWriter
	s          int
	onlyHeader bool
}

func (hw *HashWriter) Header() http.Header {
	return hw.w.Header()
}

func (hw *HashWriter) Write(b []byte) (n int, err error) {
	hash := md5.Sum(b)
	hstr := base64.RawURLEncoding.EncodeToString(hash[:])
	hw.Header().Add("Resp-Hash", hstr)
	if hw.s != 0 {
		hw.w.WriteHeader(hw.s)
	}
	hw.onlyHeader = false
	return hw.w.Write(b)
}

func (hw *HashWriter) WriteHeader(s int) {
	hw.s = s
	hw.onlyHeader = true
}

// Middleware is a utility function to wrap HTTP handlers.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hw := &HashWriter{w: w}
		next.ServeHTTP(hw, r)
		if hw.onlyHeader {
			hw.w.WriteHeader(hw.s)
		}
	})
}

// EchoMiddleware is similar to `Middleware`, but for echo package.
func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		hw := &HashWriter{w: c.Response().Writer}
		c.Response().Writer = hw
		if err := next(c); err != nil {
			return err
		}
		if hw.onlyHeader {
			hw.w.WriteHeader(hw.s)
		}
		return nil
	}
}
