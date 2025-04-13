package db

import (
	"crypto/rand"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const rounds = 13

type User struct {
	Id            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	PasswordHash  []byte    `json:"password_hash"`
	StreamKeyHash []byte    `json:"streamkey_hash"`
}

func NewUser(username, password, streamKey string) (User, error) {
	user := User{Username: username}

	id, err := uuid.NewV7()
	if err != nil {
		return User{}, err
	}

	user.Id = id
	err = user.setPassword(password)
	if err != nil {
		return User{}, err
	}

	err = user.setStreamKey(streamKey)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	return err == nil
}

func (u *User) VerifyStreamKey(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.StreamKeyHash, []byte(password))
	return err == nil
}

func (u *User) setPassword(password string) error {
	pHash, err := bcrypt.GenerateFromPassword([]byte(password), rounds)
	if err != nil {
		return err
	}

	u.PasswordHash = pHash

	return nil
}

func (u *User) setStreamKey(streamKey string) error {
	skHash, err := bcrypt.GenerateFromPassword([]byte(streamKey), rounds)
	if err != nil {
		return err
	}

	u.StreamKeyHash = skHash

	return nil
}

func (u *User) ChangePassword(currentPassword, newPassword string) error {
	if !u.VerifyPassword(currentPassword) {
		return errors.New("passwords did not match")
	}

	return u.setPassword(newPassword)
}

func GenerateStreamKey() string {
	return rand.Text()
}
