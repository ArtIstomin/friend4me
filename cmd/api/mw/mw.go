package mw

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/secure"
)

// Add adds middlewares to gin engine
func Add(r *echo.Echo, m ...echo.MiddlewareFunc) {
	for _, v := range m {
		r.Use(v)
	}
}

// SecureHeaders adds general security headers for basic security measures
func SecureHeaders() echo.MiddlewareFunc {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		ForceSTSHeader:       true,
		STSSeconds:           5184000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
	})

	return echo.WrapMiddleware(secureMiddleware.Handler)
}

// CORS adds Cross-Origin Resource Sharing support
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		MaxAge:           86400,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
