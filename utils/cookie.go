package utils

import (
	"net/http"

	"go-auth/config"
)

// SetSessionCookie sets the NextAuth session cookie
func SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "next-auth.session-token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 365 * config.TOKEN_EXPIRY_YEARS,
	}
	http.SetCookie(w, cookie)
}

// ClearSessionCookie clears the session cookie (for logout)
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "next-auth.session-token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete cookie
	}
	http.SetCookie(w, cookie)
}

// GetSessionCookie retrieves the session cookie from request
func GetSessionCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("next-auth.session-token")
}