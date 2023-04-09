package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type Models struct {
	User User
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

func getPasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *User) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, password from users where email = $1`

	row := db.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Insert user into table users
func (u *User) InsertUser(email, password string) (int, error) {
	log.Printf("Called func InsertUser, arg$1:%s, arg$2:%s", email, password)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// hashing the password, never store your password in plain text
	hashedPassword, err := getPasswordHash(password)
	log.Println(string(hashedPassword))
	if err != nil {
		return 0, err
	}
	var id int
	sql := `insert into users (email, password) values ($1, $2) returning id`
	err = db.QueryRowContext(ctx, sql, email, hashedPassword).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) GetUserById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	sql := `select id, email from users where id = $1`

	row := db.QueryRowContext(ctx, sql, id)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) UpdatePassword(oldPassword, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	// validate users oldPassword
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword))
	if err != nil {
		return false, errors.New("Password Mismatch")
	}

	// update password
	hashedPassword, err := getPasswordHash(password)
	if err != nil {
		return false, errors.New("Failed in update password, issue during the state of processing the password hash")
	}
	sql := `update users set password = $1 where id = $2`
	row := db.QueryRowContext(ctx, sql, hashedPassword, u.ID)
	err = row.Err()
	if err != nil {
		return false, errors.New("Failed in updating password")
	}
	return true, nil
}

func (u *User) PasswordMatched(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false, errors.New("Password Mismatch")
	} else {
		return true, nil
	}
}

// validate email format
func (u *User) ValidateEmail(email string) (bool, error) {
	pattern := `^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false, err
	}
	return match, nil
}

// validate password complexity
func (u *User) ValidatePassword(password string) (bool, error) {
	type Check struct {
		min bool
	}
	var check Check
	if len(password) < 8 {
		check.min = false
	} else {
		check.min = true
	}

	// TODO: adding more checkings for password complexity
	if check.min {
		return true, nil
	}
	return false, nil
}
