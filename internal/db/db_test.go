package db

import (
	"os"
	"path/filepath"
	"testing"
)

func createTmpFilePath() (string, string, error) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		return "", "", err
	}

	path := filepath.Join(dir, "test.json")

	return dir, path, nil
}

func TestOpen(t *testing.T) {
	dir, path, err := createTmpFilePath()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	db, err := Open(path)
	if err != nil {
		t.Fatalf("couldn't open db at path: %s", path)
	}

	if len(db.Users()) != 0 {
		t.Fatalf("db at path: %s already has users", path)
	}

	_, err = os.Stat(path)
	if err == nil {
		t.Fatal(err)
	}
}

func TestUser(t *testing.T) {
	const eName = "foo"
	const ePw = "bar"
	const eSk = "baz"

	dir, path, err := createTmpFilePath()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	db, err := Open(path)
	if err != nil {
		t.Fatalf("couldn't open db at path: %s", path)
	}

	err = db.AddUser(eName, ePw, eSk)
	if err != nil {
		t.Fatalf("couldn't add user %s", err)
	}

	if len(db.users) != 1 {
		t.Errorf("db does not have 1 user after adding the first")
	}

	user := db.Users()[eName]
	if user.Username != eName {
		t.Errorf("username wrong. expected: %s, actual: %s", eName, user.Username)
	}

	if !user.VerifyPassword(ePw) {
		t.Error("couldn't verify password")
	}

	if !user.VerifyStreamKey(eSk) {
		t.Error("couldn't verify streamKey")
	}

	const newPw = "FOOBARBAZ"
	err = db.ChangePassword(eName, ePw, newPw)
	if err != nil {
		t.Error("error changing password")
	}

	user = db.Users()[eName]
	if !user.VerifyPassword(newPw) {
		t.Error("couldn't verify new password")
	}

	db.RemoveUser(eName)

	if len(db.users) != 0 {
		t.Errorf("db does not have 0 users after removing the only one")
	}
}

func verifyUser(t *testing.T, user User, password, streamKey string) {
	if !user.VerifyPassword(password) {
		t.Errorf("couldn't verify password for %s", user.Username)
	}

	if !user.VerifyStreamKey(streamKey) {
		t.Errorf("couldn't verify streamkey for %s", user.Username)
	}

}

func TestExistingDb(t *testing.T) {
	db, _ := Open("existing_db_test.json")

	if len(db.users) != 4 {
		t.Errorf("db does not have 4 users")
	}

	users := db.Users()

	verifyUser(t, users["John"], "foo", "bar")
	verifyUser(t, users["Jane"], "lorem", "ipsum")
	verifyUser(t, users["Angelo"], "zonk", "bonk")
	verifyUser(t, users["Angela"], "dum", "dum")
}
