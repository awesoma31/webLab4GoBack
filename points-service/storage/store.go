package storage

import (
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

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type ConfigOption func(*Config)

func NewConfig(opts ...ConfigOption) *Config {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}

func WithHost(host string) ConfigOption {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port string) ConfigOption {
	return func(c *Config) {
		c.Port = port
	}
}

func WithUsername(username string) ConfigOption {
	return func(c *Config) {
		c.Username = username
	}
}

func WithPassword(password string) ConfigOption {
	return func(c *Config) {
		c.Password = password
	}
}

func WithDBName(dbName string) ConfigOption {
	return func(c *Config) {
		c.DBName = dbName
	}
}

func defaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "awesoma",
		Password: "1",
		DBName:   "lab4",
	}
}

func NewStore(opts ...ConfigOption) PointsStore {
	conf := NewConfig(opts...)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		conf.Host, conf.Username, conf.Password, conf.DBName, conf.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	log.Printf("Connected to database: %s:%s:%s\n", conf.Host, conf.DBName, conf.Port)

	return &pointsStoreImpl{db: db}
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
