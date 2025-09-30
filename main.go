package main

import (
	"go-auth/db"
	"go-auth/handlers"
	"go-auth/models"
	"log"
	"net/http"
)

func main() {
	db.Connect()
	
	
	db.DB.AutoMigrate(&models.User{}, &models.Token{})

	// Authentication routes
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/login", handlers.LoginFormHandler)     // Browser login (form + redirect)
	http.HandleFunc("/logout", handlers.LogoutHandler)       // Logout
	http.HandleFunc("/session", handlers.SessionHandler)     // Get current session
	
	// User management routes
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleSingleUser)

	// Token routes
	http.HandleFunc("/tokens/generation", handlers.GenerateTokenHandler)
	http.HandleFunc("/tokens/verify", handlers.VerifyTokenHandler)
	http.HandleFunc("/verify", handlers.VerifyAndLoginHandler)
	
	log.Println("Server started at: http://localhost:8080")
	log.Println("Available endpoints:")
	log.Println("POST   /auth              - Login with username/password")
	log.Println("POST   /users             - Create new user")
	log.Println("GET    /users             - Get all users")
	log.Println("PUT    /users/{id}        - Update user")
	log.Println("DELETE /users/{id}        - Delete user")
	log.Println("POST   /tokens/generation   - Generate login token")
	log.Println("POST   /tokens/verify     - Verify token (API)")
	log.Println("GET    /verify            - Verify token and login (Browser)")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlers.CreateUserHandler(w, r)
	case "GET":
		handlers.GetUsersHandler(w, r)
	case "OPTIONS":
		handlers.CreateUserHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSingleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		handlers.UpdateUserHandler(w, r)
	case "DELETE":
		handlers.DeleteUserHandler(w, r)
	case "OPTIONS":
		handlers.UpdateUserHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}