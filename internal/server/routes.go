package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(SecurityHeadersMiddleware())
	e.Use(XSSProtectionMiddleware())
	e.Use(ContentTypeOptionsMiddleware())
	e.Use(ContentSecurityPolicyMiddleware())
	e.Use(StrictTransportSecurityMiddleware())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentTypeNosniff: "nosniff",
	}))
	e.Static("/tmp", "./.tmp")
	e.GET("/", s.HelloWorldHandler)
	e.POST("/user", s.AddUser)
	e.GET("/health", s.HealthHandler)

	return e
}

// SecurityHeadersMiddleware sets common security-related headers.
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			return next(c)
		}
	}
}

// XSSProtectionMiddleware sets the X-Xss-Protection header to enable XSS protection.
func XSSProtectionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Xss-Protection", "1; mode=block")
			return next(c)
		}
	}
}

// ContentTypeOptionsMiddleware sets the X-Content-Type-Options header to prevent MIME sniffing.
func ContentTypeOptionsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			return next(c)
		}
	}
}

// ContentSecurityPolicyMiddleware sets the Content-Security-Policy header to restrict resources.
func ContentSecurityPolicyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; font-src 'self'")
			return next(c)
		}
	}
}

// StrictTransportSecurityMiddleware sets the Strict-Transport-Security header to enable HSTS.
func StrictTransportSecurityMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			return next(c)
		}
	}
}