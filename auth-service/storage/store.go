package storage

import (
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

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type ConfigOption func(*Config)

func defaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "awesoma",
		Password: "1",
		DBName:   "lab4",
	}
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

func NewConfig(opts ...ConfigOption) *Config {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}

func NewUserStore(opts ...ConfigOption) UserStore {
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
