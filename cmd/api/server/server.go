package server

import (
	"net/http"
	"time"

	"github.com/artistomin/friend4me/cmd/api/config"
	"github.com/artistomin/friend4me/cmd/api/mw"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New instantates new Echo server
func New() *echo.Echo {
	e := echo.New()
	mw.Add(e, middleware.Logger(), middleware.Recover(),
		mw.CORS(), mw.SecureHeaders())
	e.GET("/", healthCheck)
	e.Validator = &CustomValidator{V: validator.New()}
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler
	e.Binder = &CustomBinder{}
	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

// Start starts echo server
func Start(e *echo.Echo, cfg *config.Server) {
	e.Server.Addr = cfg.Port
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Minute
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Minute
	e.Debug = cfg.Debug
	e.Logger.Fatal(e.Start(cfg.Port))
}
