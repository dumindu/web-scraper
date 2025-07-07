package repository

import "web-scraper.dev/internal/model"

func (db *Db) CreateUser(u *model.User) error {
	return db.Create(u).Error
}

func (db *Db) ReadUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Db) ReadUserWithActivationTokenByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := db.Preload("ActivationToken").Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
