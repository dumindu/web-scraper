package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"web-scraper.dev/internal/model"
)

func (db *Db) CreateOrUpdateUserActivationTokenByUserId(uat *model.UserActivationToken) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{"token": uat.Token, "token_expired_at": uat.TokenExpiredAt}),
	}).Create(&uat).Error
}

func (db *Db) DeleteUserActivationTokenByUserId(userId uuid.UUID) error {
	return db.Where("user_id = ?", userId).Delete(&model.UserActivationToken{}).Error
}
