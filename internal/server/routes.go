package server

import (
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"server-template/pkg/util"
)

var (
	sessionSecret = util.GetEnvOrDefault("SESSION_SECRET_KEY","")
    sessionStore  = sessions.NewCookieStore([]byte(sessionSecret)) 
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(SecurityHeadersMiddleware())
	e.Use(CORSConfigMiddleware())
	//	e.Use(CSRFConfigMiddleware())
	e.Static("/tmp", "./.tmp")
	e.GET("/", s.HelloWorldHandler)
	e.POST("/user", s.AddUser)
	e.GET("/health", s.HealthHandler)

	return e
}

func SessionMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
	        sess, _ := sessionStore.Get(c.Request(), "go-server") 
	        c.Set("session", sess)                               
	        return next(c)
    	}
    }
}
	
func CSRFConfigMiddleware() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
				TokenLength:  32,
				TokenLookup:  "header:X-CSRF-Token",
				CookieName:   "_csrf",
				CookieSecure: true,
			    CookieHTTPOnly: true,
			})
}

func CORSConfigMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{ "GET", "OPTIONS"}, 
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, 
		ExposeHeaders:    []string{"Set-Cookie"}, 
		AllowCredentials: true,
		MaxAge:           86400, 
	})
}

func SecurityHeadersMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Response().Header().Set("X-Frame-Options", "DENY")
            c.Response().Header().Set("X-Content-Type-Options", "nosniff")
            c.Response().Header().Set("X-Xss-Protection", "1; mode=block")
            c.Response().Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; font-src 'self'")
            c.Response().Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
            return next(c)
        }
    }
}