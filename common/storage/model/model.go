package model

import "awesoma31/common/api"

type Point struct {
	ID      uint    `gorm:"primaryKey"`
	X       float64 `gorm:"not null"`
	Y       float64 `gorm:"not null"`
	R       float64 `gorm:"not null"`
	Result  bool    `gorm:"not null"`
	OwnerID uint    `gorm:"column:user_id;not null"`
}

func NewPoint(ID uint, x float64, y float64, r float64, result bool, ownerID uint) *Point {
	return &Point{ID: ID, X: x, Y: y, R: r, Result: result, OwnerID: ownerID}
}

type User struct {
	ID       uint    `gorm:"primaryKey"`
	Username string  `gorm:"not null;uniqueIndex"`
	Password string  `gorm:"not null"`
	Points   []Point `gorm:"foreignKey:user_id"`
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func NewUserFromRequest(r *api.LoginRequest) *User {
	return &User{
		Username: r.Username,
		Password: r.Password,
	}
}
