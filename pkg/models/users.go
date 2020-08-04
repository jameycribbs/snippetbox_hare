package models

import (
	"errors"
	"time"

	"github.com/jameycribbs/hare"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"hashed_password"`
	Created        time.Time `json:"created"`
	Active         bool      `json:"active"`
}

func (user *User) GetID() int {
	return user.ID
}

func (user *User) SetID(id int) {
	user.ID = id
}

func (user *User) AfterFind() {
	*user = User(*user)
}

type Users struct {
	*hare.Table
}

func NewUsers(db *hare.Database) (*Users, error) {
	tbl, err := db.GetTable("users")
	if err != nil {
		return nil, err
	}

	return &Users{Table: tbl}, nil
}

func (users *Users) Query(queryFn func(user User) bool, limit int) ([]User, error) {
	var results []User
	var err error

	for _, id := range users.IDs() {
		user := User{}

		if err = users.Find(id, &user); err != nil {
			return nil, err
		}

		if queryFn(user) {
			results = append(results, user)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}

func (users *Users) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	results, err := users.Query(func(r User) bool {
		return r.Email == email
	}, 1)
	if err != nil {
		return (err)
	}

	if len(results) > 0 {
		return ErrDuplicateEmail
	}

	user := User{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
		Created:        time.Now(),
		Active:         true,
	}

	_, err = users.Create(&user)
	if err != nil {
		return err
	}

	return nil
}

func (users *Users) Authenticate(email, password string) (int, error) {
	results, err := users.Query(func(r User) bool {
		return r.Email == email && r.Active
	}, 1)
	if err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(results[0].HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return results[0].ID, nil
}
