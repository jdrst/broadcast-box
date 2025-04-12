package auth

import (
	"errors"
	"net/http"

	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const authedValue = "authenticated"
const sessionName = "broadcast-box"
const usernameValue = "username"

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var hash, _ = bcrypt.GenerateFromPassword([]byte(os.Getenv("PASSWORD")), 14)
var passwordStore = map[string][]byte{os.Getenv("USERNAME"): hash}

func AuthHandler(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, sessionName)
		authenticated := session.Values[authedValue]

		if authenticated != nil && authenticated != false {
			next(w, r)
			return
		}

		http.Error(w, errors.New("Please log in.").Error(), http.StatusUnauthorized)
		return
	}
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	username := session.Values[usernameValue].(string)
	_, err := w.Write([]byte(username))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := r.PostForm.Get("username")
	pwd := r.PostForm.Get("password")

	if verifyPassword(pwd, passwordStore[user]) {
		setAuthed(user, true, w, r)
		_, err := w.Write([]byte(user))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, errors.New("Invalid login data.").Error(), http.StatusUnauthorized)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	setAuthed("", false, w, r)
}

func setAuthed(username string, isAuthed bool, w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	session.Values[authedValue] = isAuthed
	session.Values[usernameValue] = username
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func verifyPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
