package main

import (
	"html/template"
	"net/http"
	"time"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	err := templates.ExecuteTemplate(w, template, data)
	if err != nil {
		internalServerError(w, err)
	}
}

// / GET
func home(w http.ResponseWriter, r *http.Request) {
	claims, err := checkClaims(r)
	if err != nil {
		renderTemplate(w, "index.html", nil)
		return
	}
	auth := &Auth{Username: claims.Username, Groups: claims.Groups}
	renderTemplate(w, "index.html", auth)
}

// /register GET
func registerPage(w http.ResponseWriter, r *http.Request) {
	_, err := checkClaims(r)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "register.html", nil)
}

// /register POST
func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		badRequest(w, err)
		return
	}

	username := r.PostForm.Get("username")
	user, err := getUserByName(username)
	if err == nil {
		templateError := TemplateError{Msg: "username already exists"}
		renderTemplate(w, "register.html", templateError)
		return
	}
	user.Username = username

	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("confirm-password")
	if password != confirmPassword {
		templateError := TemplateError{Msg: "passwords do not match", Data: user}
		renderTemplate(w, "register.html", templateError)
		return
	}

	user.Password, err = generateHash(password)
	if err != nil {
		internalServerError(w, err)
		return
	}

	user.Created = time.Now()
	user, err = createUser(user)
	if err != nil {
		internalServerError(w, err)
		return
	}

	success := &Success{
		Title:    "Registration Successful",
		Route:    "/login",
		RouteMsg: "to login",
	}
	renderTemplate(w, "success.html", success)
}

// /login GET
func loginPage(w http.ResponseWriter, r *http.Request) {
	_, err := checkClaims(r)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "login.html", nil)
}

// /login POST
func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		badRequest(w, err)
		return
	}
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	user, err := getUserByName(username)
	if err != nil {
		templateError := TemplateError{Msg: "user does not exist"}
		renderTemplate(w, "login.html", templateError)
		return
	}

	err = checkHash(user.Password, password)
	if err != nil {
		templateError := TemplateError{Msg: "invalid credentials"}
		renderTemplate(w, "login.html", templateError)
		return
	}

	groups, err := getUserGroups(user.ID)
	if err != nil {
		internalServerError(w, err)
		return
	}

	tokenString, err := createJWT(user.Username, groups)
	if err != nil {
		internalServerError(w, err)
		return
	}

	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(config.JWTExpiration * time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, jwtCookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// /password GET
func passwordPage(w http.ResponseWriter, r *http.Request) {
	_, err := checkClaims(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "password.html", nil)
}

// /password POST
func password(w http.ResponseWriter, r *http.Request) {
	claims, err := checkClaims(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		badRequest(w, err)
		return
	}
	currentPassword := r.PostForm.Get("current-password")
	newPassword := r.PostForm.Get("new-password")
	confirmPassword := r.PostForm.Get("confirm-password")

	if newPassword == currentPassword {
		templateError := TemplateError{Msg: "new password is the same as current password"}
		renderTemplate(w, "password.html", templateError)
		return
	}

	if newPassword != confirmPassword {
		templateError := TemplateError{Msg: "passwords do not match"}
		renderTemplate(w, "password.html", templateError)
		return
	}

	user, err := getUserByName(claims.Username)
	if err != nil {
		internalServerError(w, err)
		return
	}

	err = checkHash(user.Password, currentPassword)
	if err != nil {
		templateError := TemplateError{Msg: "current password incorrect"}
		renderTemplate(w, "password.html", templateError)
		return
	}

	password, err := generateHash(newPassword)
	if err != nil {
		internalServerError(w, err)
		return
	}

	err = updatePassword(user.ID, password)
	if err != nil {
		internalServerError(w, err)
		return
	}

	success := &Success{
		Title:    "Password Changed",
		Route:    "/",
		RouteMsg: "to return home",
	}
	renderTemplate(w, "success.html", success)
}

// /logout GET
func logout(w http.ResponseWriter, r *http.Request) {
	_, err := checkClaims(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, jwtCookie)

	success := &Success{
		Title:    "Logged out",
		Route:    "/",
		RouteMsg: "to return home",
	}
	renderTemplate(w, "success.html", success)
}
