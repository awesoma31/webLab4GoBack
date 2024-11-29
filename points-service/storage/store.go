package storage

import (
	"awesoma31/common"
	"awesoma31/common/storage/model"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PointsStore interface {
	Create(ctx context.Context, point *model.Point) (*model.Point, error)
	GetPointsByUserIDWithPagination(ctx context.Context, userID int64, pageSize int, pageNumber int) ([]model.Point, error)
	GetTotalPointsByUserID(ctx context.Context, userID int64) (int, error)
}

type pointsStoreImpl struct {
	db *gorm.DB
}

func NewStore() PointsStore {
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
		log.Fatalf("Failed to connect to database: %s", err)
	}
	log.Printf("Connected to database: %s:%s:%s\n", host, dbname, port)

	return &pointsStoreImpl{db}
}

func (s *pointsStoreImpl) Create(ctx context.Context, point *model.Point) (*model.Point, error) {
	if err := s.db.WithContext(ctx).Create(point).Error; err != nil {
		return nil, err
	}
	return point, nil
}

func (s *pointsStoreImpl) GetPointsByUserIDWithPagination(ctx context.Context, userID int64, pageSize int, pageNumber int) ([]model.Point, error) {
	var points []model.Point
	offset := pageNumber * pageSize
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Limit(pageSize).Offset(offset).Find(&points).Error; err != nil {
		return nil, err
	}
	return points, nil
}

func (s *pointsStoreImpl) GetTotalPointsByUserID(ctx context.Context, userID int64) (int, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Point{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
