package handlers

import (
	"net/http"
	"time"

	"lab2-terrylneal/internal/jsonhelper"
)

// home handler
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := "My name is Terryl Neal\n"
	response += "My semester project is a Medical Appointment Scheduling System\n"
	response += "I chose this as my project as it seemed interesting and I believe i can already start mapping what I need to do logistically for it to work.\n\n"
	w.Write([]byte(response))
}

// about handler
func About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	response := "About Page\n"
	response += "My name is Terryl Neal and I am completing my Bachelor's in IT at UB.\n"
	response += "It's been a while since Ive handled some of the higher level classes as Ive been thorugh nothihng but general core classes in the last year and half due to my transfer.\n"
	response += "A lot of these concepts I barely remember and I still have to catch up in some respects due to that reason. I am trying my best.\n\n"
	w.Write([]byte(response))
}

// contact handler
func Contact(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
		http.NotFound(w, r)
		return
	}
	response := "Contact Page\n"
	response += "Email: 2024161285@ub.edu.bz\n"
	response += "Phone #: 614-1850"
	w.Write([]byte(response))
}

// hobby handler
func Hobby(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hobby" {
		http.NotFound(w, r)
		return
	}
	response := "Hobby Page\n"
	response += "My favourite hobby is fishing. Learned how to do it at 5 with my grandma and it's been a great destressor ever since. Gives me an appreciation for the world around me when it's just you, fish and a line.\n\n"
	w.Write([]byte(response))
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func APIInfo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/info" {
		jsonhelper.WriteJSONError(w, http.StatusNotFound, "Endpoint not found")
		return
	}

	info := map[string]interface{}{
		"app_name":    "Medical Appointment Scheduling System",
		"version":     "1.0.0",
		"author":      "Terryl Neal",
		"endpoints":   []string{"/", "/about", "/contact", "/hobby", "/api/info", "/api/users"},
		"description": "A system for scheduling medical appointments",
	}

	jsonhelper.WriteJSON(w, http.StatusOK, info)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/users" {
		jsonhelper.WriteJSONError(w, http.StatusNotFound, "Endpoint not found")
		return
	}

	if r.Method != http.MethodGet {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	users := []User{
		{ID: 1, Username: "john_doe", Email: "john@example.com", CreatedAt: time.Now()},
		{ID: 2, Username: "jane_smith", Email: "jane@example.com", CreatedAt: time.Now()},
	}

	jsonhelper.WriteJSON(w, http.StatusOK, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/users" {
		jsonhelper.WriteJSONError(w, http.StatusNotFound, "Endpoint not found")
		return
	}

	if r.Method != http.MethodPost {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newUser User
	err := jsonhelper.ReadJSON(w, r, &newUser)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	// Validation of needed fields
	if newUser.Username == "" || newUser.Email == "" {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Username and email are required")
		return
	}

	newUser.ID = 3
	newUser.CreatedAt = time.Now()

	jsonhelper.WriteJSON(w, http.StatusCreated, newUser)
}
