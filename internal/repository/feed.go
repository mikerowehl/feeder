/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package repository

import (
	"errors"

	"github.com/mikerowehl/feeder/internal/rss"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	_ "modernc.org/libc"
	_ "modernc.org/sqlite"
)

type FeedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(filename string) (*FeedRepository, error) {
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        filename,
	}, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&rss.Feed{}, &rss.Item{})
	return &FeedRepository{db: db}, nil
}

func (r *FeedRepository) Save(feed *rss.Feed) error {
	err := r.db.Save(feed).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FeedRepository) Delete(id uint) error {
	err := r.db.Unscoped().Select(clause.Associations).Delete(&rss.Feed{}, id).Error
	return err
}

func (r *FeedRepository) All() ([]rss.Feed, error) {
	var feeds []rss.Feed
	err := r.db.Preload("Items").Find(&feeds).Error
	return feeds, err
}

func (r *FeedRepository) AllFeeds() ([]rss.Feed, error) {
	var feeds []rss.Feed
	err := r.db.Find(&feeds).Error
	return feeds, err
}

func (r *FeedRepository) AllItems() ([]rss.Item, error) {
	var items []rss.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *FeedRepository) Unread() ([]rss.Feed, error) {
	var feeds []rss.Feed
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.
			Where("read = ?", false).
			Order("published DESC")
	}).Find(&feeds).Error
	return feeds, err
}

func (r *FeedRepository) MarkAll() error {
	result := r.db.Model(&rss.Item{}).Where("read = ?", false).Update("read", true)
	return result.Error
}

func (r *FeedRepository) TrimItems(feedId uint, count int) error {
	var cutoffID uint
	err := r.db.Model(&rss.Item{}).
		Where("feed_id = ?", feedId).
		Order("published DESC").
		Offset(count).
		Limit(1).
		Pluck("id", &cutoffID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// There are fewer than count recs, nothing to trim
			return nil
		}
		return err
	}

	return r.db.Where("feed_id = ? AND id <= ?", feedId, cutoffID).
		Delete(&rss.Item{}).
		Error
}

func (r *FeedRepository) Close() error {
	sqliteDb, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqliteDb.Close()
}
