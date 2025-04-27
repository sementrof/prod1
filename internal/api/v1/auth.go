package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sementrof/prod1/internal/models"
	auth "github.com/sementrof/prod1/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// LoginUserPost godoc
// @Summary login
// @Description login to account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.Userlogin true "User input"
// @Success 201 {string} string "User created with ID"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Server error"
// @Router /login [post]
func (im *ApiImplemented) LoginUserPost(w http.ResponseWriter, r *http.Request) {
	var userLogin models.Userlogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		im.log.Error(err.Error())
		return
	}
	// if user.Email == "" || user.Password == "" {
	// 	http.Error(w, "Missing required fields", http.StatusBadRequest)
	// 	im.log.Error(models.ErrLoginOrPasswordOrEmail)
	// }
	if err := validate.Struct(userLogin); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusUnprocessableEntity)
		im.log.Error(err)
		return
	}
	user := models.User{}
	dbUser, err := im.repo.GetUserByEmail(user, userLogin.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			im.log.Error("User not found")
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		im.log.Error(err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userLogin.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		im.log.Error("Invalid password")
		return
	}

	tokenString, err := auth.GenerateJWTToken(dbUser)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		im.log.Error(err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

}

// func (im *ApiImplemented) logaut(w http.ResponseWriter, r *http.Request) {

// }
