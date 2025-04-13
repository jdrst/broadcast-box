package db

import (
	"encoding/json"
	"os"
	"sync"
)

type Db struct {
	users map[string]User
	path  string
	mutex sync.RWMutex
}

// Opens a DB at the given path. If no DB exists it will be created.
func Open(path string) (*Db, error) {
	content, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		//TODO: we add an "admin" in case of no users
		db := Db{path: path, users: make(map[string]User)}
		db.AddUser("admin", "admin", "admin")
		return &db, nil
	} else if err != nil {
		return &Db{}, err
	}

	users := []User{}

	err = json.Unmarshal(content, &users)
	if err != nil {
		return &Db{}, err
	}

	userMap := make(map[string]User, len(users))
	for _, u := range users {
		userMap[u.Username] = u
	}

	return &Db{
		users: userMap,
		path:  path,
	}, nil
}

func (db *Db) AddUser(username, password, streamKey string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	user, err := NewUser(username, password, streamKey)
	if err != nil {
		return err
	}

	db.users[user.Username] = user

	return db.writeToDisk()
}

func (db *Db) RemoveUser(username string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	delete(db.users, username)

	return db.writeToDisk()
}

func (db *Db) ChangePassword(username, currentPassword, newPassword string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	user := db.users[username]

	err := user.ChangePassword(currentPassword, newPassword)
	if err != nil {
		return err
	}

	db.users[username] = user
	return nil
}

func (db *Db) Users() map[string]User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.users
}

func (db *Db) writeToDisk() error {
	users := make([]User, 0, len(db.users))
	for _, u := range db.users {
		users = append(users, u)
	}

	json, err := json.Marshal(users)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json, 0600)

	if err != nil {
		return err
	}

	return nil
}
