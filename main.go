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
	
	
	db.DB.AutoMigrate(
		&models.User{}, 
		&models.Token{},
		&models.ViewGroup{},
		&models.ViewGroupAudit{},
		&models.CustomMap{},
		&models.CameraPosition{},
	)

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

	//view group routes
	http.HandleFunc("/view-groups", handleViewGroups)  
	http.HandleFunc("/view-groups/", handleSingleViewGroup)

	//custom Map routes
	http.HandleFunc("/custom-maps", handleCustomMaps)
	http.HandleFunc("/custom-maps/", handleSingleCustomMap)
	
	//log.Println("Server started at: http://localhost:8080")
	//log.Println("Available endpoints:")
	//log.Println("POST   /auth              - Login with username/password")
	//log.Println("POST   /users             - Create new user")
	//log.Println("GET    /users             - Get all users")
	//log.Println("PUT    /users/{id}        - Update user")
	//log.Println("DELETE /users/{id}        - Delete user")
	//log.Println("POST   /tokens/generation   - Generate login token")
	//log.Println("POST   /tokens/verify     - Verify token (API)")
	//log.Println("GET    /verify            - Verify token and login (Browser)")
	//log.Println("POST   /view-groups       - Create view group") 
	//log.Println("GET    /view-groups       - Get all view groups")  
	//log.Println("PUT    /view-groups/{id}  - Update view group") 
	//log.Println("DELETE /view-groups/{id}  - Delete view group")     
	
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

func handleViewGroups(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlers.CreateViewGroupHandler(w, r)
	case "GET":
		handlers.GetViewGroupsHandler(w, r)
	case "OPTIONS":
		handlers.CreateViewGroupHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSingleViewGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		handlers.UpdateViewGroupHandler(w, r)
	case "DELETE":
		handlers.DeleteViewGroupHandler(w,r)
	case "OPTIONS":
		handlers.UpdateViewGroupHandler(w, r)
	default:
		http.Error(w, "Method not allowed",http.StatusMethodNotAllowed)
	}
}

func handleCustomMaps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlers.CreateCustomMapHandler(w,r)
	case "GET":
		handlers.GetCustomMapsHandler(w,r)
	case "OPTIONS":
		handlers.CreateCustomMapHandler(w,r)
	default:
		http.Error(w,"Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSingleCustomMap(w http.ResponseWriter,r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetCustomMapHandler(w,r)
	case "PUT":
		handlers.UpdateCustomMapHandler(w,r)
	case "DELETE":
		handlers.DeleteCustomMapHandler(w,r)
	case "OPTIONS":
		handlers.UpdateCustomMapHandler(w,r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}