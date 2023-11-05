# resphash
Add hash of response body to header

## Usage

```go
import (
    "http"

    "github.com/labstack/echo/v4"
    "github.com/nakario/resphash"
)

func main() {
    http.Handle("/", resphash.Middleware(yourHandler))

    // If you use echo
    e := echo.New()
    e.Use(resphash.EchoMiddleware)
}
```
