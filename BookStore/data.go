package main

import (
	"context"
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"time"
)

const (
	defaultShelfSize = 10
)

// 使用gorm

func NewDB(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	db.AutoMigrate(&Shelf{}, &Book{})
	return db, nil
}

// Shelf 书架
type Shelf struct {
	ID       int64 `gorm:"primaryKey"`
	Theme    string
	Size     int64
	CreateAt time.Time
	UpdataAt time.Time
}

type Book struct {
	ID       int64 `gorm:"primaryKey"`
	Author   string
	Title    string
	ShelfID  int64
	CreateAt time.Time
	UpdataAt time.Time
}

// 数据库操作
type bookstore struct {
	db *gorm.DB
}

// CreateShelf 创建书架
func (b *bookstore) CreateShelf(ctx context.Context, data Shelf) (*Shelf, error) {
	if len(data.Theme) <= 0 {
		return nil, errors.New("invalid theme")
	}
	size := data.Size
	if size <= 0 {
		size = defaultShelfSize
	}
	v := Shelf{Theme: data.Theme, Size: size, CreateAt: time.Now(), UpdataAt: time.Now()}
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}

func (b *bookstore) GetShelf(ctx context.Context, id int64) (*Shelf, error) {
	v := Shelf{}
	err := b.db.WithContext(ctx).First(&v, id).Error
	return &v, err
}

func (b *bookstore) ListShelves(ctx context.Context) ([]*Shelf, error) {
	var vl []*Shelf
	err := b.db.WithContext(ctx).Find(&vl).Error
	return vl, err
}

func (b *bookstore) DeleteShelf(ctx context.Context, id int64) error {
	return b.db.WithContext(ctx).Delete(&Shelf{}, id).Error
}

// ListBooks 根据书架id查询图书
func (b *bookstore) ListBooks(ctx context.Context, sid int64, cursor string, pageSize int) ([]*Book, error) {
	var bl []*Book
	err := b.db.Debug().WithContext(ctx).Where("shelf_id = ? and id > ?", sid, cursor).Order("id asc").Limit(pageSize).Find(&bl).Error
	return bl, err
}

func (b *bookstore) CreateBook(ctx context.Context, data Book) (*Book, error) {
	v := Book{Author: data.Author, Title: data.Title, ShelfID: data.ShelfID}
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}

func (b *bookstore) DeleteBook(ctx context.Context, sid int64, id int64) error {
	return b.db.WithContext(ctx).Where("shelf_id = ?", sid).Delete(&Book{}, id).Error
}
