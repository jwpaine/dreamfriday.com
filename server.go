package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	auth "dreamfriday/auth"
	Database "dreamfriday/database"
	Middleware "dreamfriday/middleware"
	routes "dreamfriday/routes"
)

// Load environment variables
func init() {
	// Load environment variables from .env file
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	// Use the strings directly as raw keys
	Database.ConnStr = os.Getenv("DATABASE_CONNECTION_STRING")
	if Database.ConnStr == "" {
		log.Fatal("DATABASE_CONNECTION_STRING environment variable not set")
	}

}

type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	// Initialize the database connection
	db, err := Database.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	e := echo.New()

	// allow CORS for https://static.cloudflareinsights.com and https://dreamfriday.com:
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://static.cloudflareinsights.com", "https://dreamfriday.com"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	// Add middleware to load site data once
	e.Use(Middleware.LoadSiteDataMiddleware)

	auth.InitSessionStore()

	routes.RegisterRoutes(e)

	e.GET("/admin/:domain", AdminSite) // @TODO: use JSON-based page instead

	e.Static("/static", "static")

	e.GET("/favicon.ico", func(c echo.Context) error {
		// Serve the favicon.ico file from the static directory or a default location
		return c.File("static/favicon.ico")
	})

	listener, err := net.Listen("tcp4", "0.0.0.0:8081")
	if err != nil {
		log.Fatalf("Failed to bind to IPv4: %v", err)
	}
	server := &http.Server{
		Handler: e, // Pass the Echo instance as the handler
	}
	log.Println("Starting server on IPv4 address 0.0.0.0:8081...")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// LoginForm renders a simple login form
func LoginForm(c echo.Context) error {
	// Retrieve session
	session, err := auth.GetSession(c.Request())
	if err != nil {
		log.Println("Failed to retrieve session:", err)
		return c.String(http.StatusInternalServerError, "Session error")
	}

	// Check if user is already logged in
	if session.Values["handle"] != nil {
		log.Println("User already logged in, redirecting to admin panel")
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Render the login page
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"title":   "Login",
		"message": "Ready for login",
	})
}

// /admin/:domain route
// @TODO: use JSON-based page instead
func AdminSite(c echo.Context) error {
	// Retrieve the session
	log.Println("AdminSite")
	session, err := auth.GetSession(c.Request())
	if err != nil {
		log.Println("Failed to get session:", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve session")
	}

	identifier, ok := session.Values["handle"].(string)
	if !ok || identifier == "" {
		log.Println("Unauthorized: Identifier (email or handle) not found in session")
		return c.String(http.StatusUnauthorized, "Unauthorized: No valid identifier found")
	}

	// Retrieve domain from /admin/:domain route
	domain := c.Param("domain")
	log.Println("Pulling preview data for Domain:", domain)

	// Fetch preview data from the database using the identifier
	previewData, status, err := Database.FetchPreviewData(domain, identifier)
	if err != nil {
		log.Println("Failed to fetch preview data for domain:", domain, "Error:", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch preview data for domain: "+domain)
	}

	// Convert previewData (*Models.SiteData) to a formatted JSON string
	previewDataBytes, err := json.MarshalIndent(previewData, "", "    ")
	if err != nil {
		log.Println("Failed to format preview data:", err)
		return c.String(http.StatusInternalServerError, "Failed to format preview data")
	}

	// Convert JSON byte array to string
	previewDataString := string(previewDataBytes)

	// Pass the formatted JSON string to the view
	return c.Render(http.StatusOK, "manage.html", map[string]interface{}{
		"domain":      domain,
		"previewData": previewDataString,
		"status":      status,
		"message":     "",
	})
}

/* place holder password reset for auth0

func PasswordResetForm(c echo.Context) error {
	return HTML(c, Views.PasswordReset())
}

// PasswordReset handles the password reset form submission and calls auth0PasswordReset

func PasswordReset(c echo.Context) error {
	email := c.FormValue("email")
	err := Auth.PasswordReset(email)
	if err != nil {
		return HTML(c, Views.PasswordResetFailed())
	}
	return HTML(c, Views.ConfirmPasswordReset(email))
} */

// RegisterForm renders the registration form

/*
func RegisterForm(c echo.Context) error {
	return RenderTemplate(c, http.StatusOK, Views.Register())
}
/*

Place holder Registeration support for auth0

func Register(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Validate input fields
	if email == "" || password == "" {
		log.Println("Registration failed: Email and password are required")
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"message": "Email and password are required",
		})
	}
	// Ensure authenticator is an Auth0Authenticator
	auth0Auth, ok := authenticator.(*auth.Auth0Authenticator)
	if !ok {
	logPrintln("rror: Autenticatoris not anAuth0 insance")
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	// Register the user via Auth0
	_, err := auth0Auth.Register(email, password)
	if err != nil {
		log.Printf("Registration error for %s: %v", email, err)
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"message": "Registration failed: " + err.Error(),
		})
	}
	// Successfully registered, show confirmation page
	log.Printf("User %s registered successfully", email)
	return c.Render(http.StatusOK, "register_success.html", map[string]interface{}{
		"email": email,
	})
}
*/
