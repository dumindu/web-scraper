package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"web-scraper.dev/internal/model"
)

type Db struct {
	*gorm.DB
}

func New(db *gorm.DB) *Db {
	return &Db{db}
}

func (db *Db) TxBegin() *Db {
	tx := db.DB.Begin()
	return &Db{tx}
}

func (db *Db) Commit() {
	db.DB.Commit()
}

func (db *Db) Rollback() {
	db.DB.Rollback()
}

type DB interface {
	TxBegin() *Db
	Commit()
	Rollback()

	CreateUser(u *model.User) error
	ReadUserByEmail(email string) (*model.User, error)
	ReadUserWithActivationTokenByEmail(email string) (*model.User, error)

	CreateOrUpdateUserActivationTokenByUserId(uat *model.UserActivationToken) error
	DeleteUserActivationTokenByUserId(userId uuid.UUID) error
}
