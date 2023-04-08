package data

import (
	"context"
	"database/sql"
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

func (u *User) DebugBackdoor(sql string) error {
	return nil
}

// Insert user into table users
func (u *User) InsertUser(email, password string) (int, error) {
	log.Printf("Called func InsertUser, arg$1:%s, arg$2:%s", email, password)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// hashing the password, never store your password in plain text
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}
	var id int
	sql := `insert into users (email, password) values ($1, $2) returning id`
	log.Printf("Called func InsertUser, sql:%s", sql)
	err = db.QueryRowContext(ctx, sql, email, hashedPassword).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) ValidateEmail(email string) (bool, error) {
	pattern := `^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false, err
	}
	return match, nil
}

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
