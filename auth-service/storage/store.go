package storage

import (
	"awesoma31/common"
	"awesoma31/common/storage/model"
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var ErrUserNotFound = errors.New("user not found")

type UserStore interface {
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	GetPointsByUserID(ctx context.Context, userID uint) ([]model.Point, error)
	FindIdByUsername(ctx context.Context, username string) (int64, error)
}

type userStoreImpl struct {
	db *gorm.DB
}

func NewUserStore() UserStore {
	host := common.GetEnv("DB_HOST", "localhost")
	port := common.GetEnv("DB_PORT", "5432")
	user := common.GetEnv("DB_USER", "awesoma")
	password := common.GetEnv("DB_PASSWORD", "1")
	dbname := common.GetEnv("DB_NAME", "lab4")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	log.Printf("Connected to database: %s:%s:%s\n", host, dbname, port)

	return &userStoreImpl{db}
}

func (r *userStoreImpl) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Points").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userStoreImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userStoreImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userStoreImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userStoreImpl) GetPointsByUserID(ctx context.Context, userID uint) ([]model.Point, error) {
	var points []model.Point
	if err := r.db.WithContext(ctx).Where("owner_id = ?", userID).Find(&points).Error; err != nil {
		return nil, err
	}
	return points, nil
}

func (r *userStoreImpl) FindIdByUsername(ctx context.Context, username string) (int64, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return int64(user.ID), nil
}
