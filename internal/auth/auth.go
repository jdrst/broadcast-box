package auth

import (
	"errors"
	"net/http"

	"github.com/glimesh/broadcast-box/internal/database"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	authedValue   = "authenticated"
	sessionName   = "broadcast-box"
	usernameValue = "username"
	passwordValue = "password"
)

type AuthContext struct {
	Db    *database.Queries
	store sessions.Store
}

func NewContext(db *database.Queries, sessionKey []byte) AuthContext {
	return AuthContext{
		Db:    db,
		store: sessions.NewCookieStore(sessionKey),
	}
}

func (ctx *AuthContext) AuthHandler(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := ctx.store.Get(r, sessionName)
		authenticated := session.Values[authedValue]

		if authenticated != nil && authenticated != false {
			next(w, r)
			return
		}

		http.Error(w, errors.New("Please log in.").Error(), http.StatusUnauthorized)
		return
	}
}

func (ctx *AuthContext) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := ctx.store.Get(r, sessionName)
	username := session.Values[usernameValue].(string)
	_, err := w.Write([]byte(username))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ctx *AuthContext) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.PostForm.Get(usernameValue)
	pwd := r.PostForm.Get(passwordValue)

	user, err := GetUser(r.Context(), ctx.Db, username)
	if err != nil {
		http.Error(w, errors.New("Invalid login data.").Error(), http.StatusUnauthorized)
		return
	}

	if user.VerifyPassword(pwd) {
		ctx.setAuthed(username, true, w, r)
		_, err := w.Write([]byte(username))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, errors.New("Invalid login data.").Error(), http.StatusUnauthorized)
}

func (ctx *AuthContext) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx.setAuthed("", false, w, r)
}

func (ctx *AuthContext) setAuthed(username string, isAuthed bool, w http.ResponseWriter, r *http.Request) {
	session, _ := ctx.store.Get(r, sessionName)
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
