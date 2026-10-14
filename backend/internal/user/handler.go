package user

import (
	"encoding/json"
	"errors"
	"evolve-server/internal/util"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

type UserHandler struct {
	user *UserService
}

func NewUserHandler(user *UserService) *UserHandler {
	return &UserHandler{user: user}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if !req.validate() {
		util.SendError(w, http.StatusBadRequest, "not all fields valid")
		return
	}

	token, err := h.user.SignUp(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			util.SendError(w, http.StatusConflict, ErrEmailTaken.Error())
			return
		}

		util.SendError(w, http.StatusInternalServerError, "internal server error")
		log.Println(err)
		return
	}

	util.SendSuccess(w, http.StatusCreated, authResponse{Token: token})
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if !req.validate() {
		util.SendError(w, http.StatusBadRequest, "not all fields valid")
		return
	}

	token, err := h.user.SignIn(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			util.SendError(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
			return
		}

		util.SendError(w, http.StatusInternalServerError, "internal server error")
		log.Println(err)
		return
	}

	util.SendSuccess(w, http.StatusCreated, authResponse{Token: token})
}

func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	// obtain userID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	user, err := h.user.Me(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserDoesNotExist) {
			util.SendError(w, http.StatusBadRequest, ErrUserDoesNotExist.Error())
			return
		}
		util.SendError(w, http.StatusInternalServerError, "internal server error")
		log.Println(err)
		return
	}

	util.SendSuccess(w, http.StatusOK, user)
}

type signUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *signUpRequest) validate() bool {
	// check if email is valid
	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return false
	}

	// check if username is valid
	if utf8.RuneCountInString(r.Username) > 15 || utf8.RuneCountInString(r.Username) < 4 {
		return false
	}

	// check if password is valid
	if !isValidPassword(r.Password) {
		return false
	}

	return true
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *signInRequest) validate() bool {
	// check if email is valid
	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return false
	}

	return true
}

type authResponse struct {
	Token string `json:"token"`
}

var (
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasDigit   = regexp.MustCompile(`[0-9]`)
	hasSpecial = regexp.MustCompile(`[#?!@$%^&*-]`)
	isASCII    = regexp.MustCompile(`^[\x00-\x7F]+$`)
)

func isValidPassword(password string) bool {
	return len(password) >= 8 &&
		len(password) <= 32 &&
		isASCII.MatchString(password) &&
		hasUpper.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasDigit.MatchString(password) &&
		hasSpecial.MatchString(password)
}
