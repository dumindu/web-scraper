package repository

import (
	"github.com/google/uuid"

	"web-scraper.dev/internal/model"
)

func (db *Db) ListKeywordsByUserId(userID uuid.UUID) (model.Keywords, error) {
	keywords := make([]*model.Keyword, 0)
	if err := db.Where("user_id = ?", userID).Order("created_at desc, updated_at desc").Find(&keywords).Error; err != nil {
		return nil, err
	}
	return keywords, nil
}

func (db *Db) ReadKeywordByIdAndUserId(id int64, userId uuid.UUID) (*model.Keyword, error) {
	keyword := &model.Keyword{}
	if err := db.Where("id = ? AND user_id = ?", id, userId).First(keyword).Error; err != nil {
		return nil, err
	}
	return keyword, nil
}
