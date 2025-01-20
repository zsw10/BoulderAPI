package data

import (
	"context"
	"crypto"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/registration"
)

type User struct {
	ID           int               `json:"id"`
	CreatedAt    time.Time         `json:"created_at"`
	Email        string            `json:"email"`
	Key          crypto.PrivateKey `json:"-"`
	Status       string            `json:"status"`
	registration *registration.Resource
}

type UserModel struct {
	DB *sql.DB
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u User) GetRegistration() *registration.Resource {
	return u.registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func (m UserModel) Insert(user *User) error {
	pKeyBytes, err := x509.MarshalPKCS8PrivateKey(user.Key)
	if err != nil {
		return err
	}

	pKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pKeyBytes,
	})

	query := `
	INSERT INTO USERS (id, created_at, email, key, status)
	VALUES ($1, $2, $3, $4, $5)`

	args := []any{user.ID, user.CreatedAt, user.Email, pKeyPEM, user.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) GetByID(id int) (*User, error) {
	query := `
	SELECT id, created_at, email, key, status
	FROM USERS
	WHERE id = $1`

	var user User
	var pemKey string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.CreatedAt, &user.Email, &pemKey, &user.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %v not found", id)
		}
		return nil, err
	}

	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	user.Key = key

	return &user, nil
}
