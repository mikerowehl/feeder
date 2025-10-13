/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package repository

import (
	"github.com/mikerowehl/feeder/internal/rss"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(filename string) (*FeedRepository, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&rss.Feed{}, &rss.Item{})
	return &FeedRepository{db: db}, nil
}

func (r *FeedRepository) Save(feed *rss.Feed) error {
	err := r.db.Save(feed).Error
	if err != nil {
		return err
	}
	for _, item := range feed.Items {
		if item.ID != 0 {
			err := r.db.Save(&item).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *FeedRepository) All() ([]rss.Feed, error) {
	var feeds []rss.Feed
	err := r.db.Preload("Items").Find(&feeds).Error
	return feeds, err
}

func (r *FeedRepository) Close() error {
	sqliteDb, err := r.db.DB()
	if err != nil {
		return err
	}
	sqliteDb.Close()
	return nil
}
