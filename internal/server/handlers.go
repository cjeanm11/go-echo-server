
package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type userData struct {
	Username string `json:"username" validate:"required"`    
	Email    string `json:"email" validate:"required,email"` 
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{"message": "Hello World"}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) AddUser(c echo.Context) error {
	var user userData
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body format"})
	}
	userResponse := s.db.AddUser(user.Username, user.Email, "")
	if userResponse["error"] != "" {
		return c.JSON(http.StatusInternalServerError, userResponse)
	}
	return c.JSON(http.StatusCreated, userResponse)
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}