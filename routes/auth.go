package routes

import (
	"dreamfriday/auth"
	"dreamfriday/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes registers login/logout routes
func RegisterAuthRoutes(e *echo.Echo) {
	authenticator := auth.GetAuthenticator() // Get the instance of Auth0Authenticator
	authHandler := handlers.NewAuthHandler(authenticator)

	e.POST("/login", authHandler.Login)
	e.GET("/logout", authHandler.Logout)
}
