package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"log"

	"github.com/glimesh/broadcast-box/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const rounds = 13

type User struct {
	user *database.User
}

func NewUser(ctx context.Context, queries *database.Queries, username, password, streamKey string) (User, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return User{}, err
	}

	pw, err := hash(password)
	if err != nil {
		return User{}, err
	}

	sk, err := hash(streamKey)
	if err != nil {
		return User{}, err
	}

	u, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        id.String(),
		Name:      username,
		Password:  pw,
		Streamkey: sk,
	})

	return User{user: &u}, nil
}

func RemoveUser(ctx context.Context, queries *database.Queries, username string) error {
	return queries.DeleteUser(ctx, username)
}

func GetUser(ctx context.Context, queries *database.Queries, username string) (*User, error) {
	u, err := queries.GetUser(ctx, username)
	if err != nil {
		return &User{}, err
	}

	return &User{user: &u}, nil
}

func ListUsers(ctx context.Context, queries *database.Queries) (map[string]*User, error) {
	users, err := queries.ListUsers(ctx)
	if err != nil {
		return map[string]*User{}, err
	}

	res := make(map[string]*User, len(users))
	for _, u := range users {
		res[u.Name] = &User{user: &u}
	}

	return res, nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.user.Password, []byte(password))
	return err == nil
}

func (u *User) VerifyStreamKey(streamkey string) bool {
	err := bcrypt.CompareHashAndPassword(u.user.Streamkey, []byte(streamkey))
	return err == nil
}

func (u *User) ChangePassword(ctx context.Context, queries *database.Queries, current, new string) error {
	if !u.VerifyPassword(current) {
		return errors.New("passwords did not match")
	}

	hash, err := hash(new)
	if err != nil {
		return err
	}

	err = queries.UpdateUser(ctx, database.UpdateUserParams{
		Name:      u.user.Name,
		Password:  hash,
		Streamkey: u.user.Streamkey,
		ID:        u.user.ID,
	})
	if err != nil {
		return err
	}

	u.user.Password = hash
	return nil
}

func (u *User) ChangeStreamkey(ctx context.Context, queries *database.Queries, new string) error {
	hash, err := hash(new)
	if err != nil {
		return err
	}

	err = queries.UpdateUser(ctx, database.UpdateUserParams{
		Name:      u.user.Name,
		Password:  u.user.Password,
		Streamkey: hash,
		ID:        u.user.ID,
	})
	if err != nil {
		return err
	}

	u.user.Streamkey = hash
	return nil
}

func GenerateStreamKey() string {
	return rand.Text()
}

func GetUsernameForStreamkey(ctx context.Context, queries *database.Queries, streamKey string) (bool, string) {
	log.Printf("streamkey: (%s)", streamKey)
	hash, err := hash(streamKey)
	//TODO return actual error
	if err != nil {
		return false, ""
	}
	name, err := queries.GetUsernameForStreamKey(ctx, hash)
	if err != nil {
		log.Print(err)
		log.Printf("hash was (%s)", hash)
		users, _ := queries.ListUsers(ctx)
		for _, u := range users {
			log.Printf("User: %s, SK: %s", u.Name, u.Streamkey)
		}
	}

	//TODO actually we want to check for ErrNoRows and return the error otherwise
	return err == nil, name
}

func (u User) String() string {
	return u.user.ID + " " + u.user.Name
}

func hash(input string) ([]byte, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(input), rounds)
	if err != nil {
		return []byte{}, err
	}

	return h, nil
}
