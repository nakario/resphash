package resphash_test

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nakario/resphash"
)

func TestResphash(t *testing.T) {
	respBody := bytes.NewBufferString("body")
	hash := md5.Sum(respBody.Bytes())
	hstr := base64.RawURLEncoding.EncodeToString(hash[:])

	m := http.NewServeMux()
	m.Handle("/", resphash.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(respBody.Bytes())
	})))

	e := echo.New()
	e.Use(resphash.EchoMiddleware)
	e.GET("/", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/plain", respBody.Bytes())
	})

	handlers := []http.Handler{
		m,
		e.Server.Handler,
	}
	for i, handler := range handlers {
		req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", respBody)
		got := httptest.NewRecorder()
		handler.ServeHTTP(got, req)
		if got.Body.String() != respBody.String() {
			t.Errorf("handler %d: body shouldn't be modified", i)
		}
		header := got.Header().Get("Resp-Hash")
		if header != hstr {
			t.Errorf("hander %d: unexpected header: expected: %s, actual: %s", i, hstr, header)
		}
	}
}
