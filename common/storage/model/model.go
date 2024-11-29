package model

import "awesoma31/common/api"

type Point struct {
	ID      int64   `gorm:"primaryKey"`
	X       float64 `gorm:"not null"`
	Y       float64 `gorm:"not null"`
	R       float64 `gorm:"not null"`
	Result  bool    `gorm:"not null"`
	OwnerID int64   `gorm:"column:user_id;not null"`
}

func NewPoint(ID int64, x float64, y float64, r float64, result bool, ownerID int64) *Point {
	return &Point{ID: ID, X: x, Y: y, R: r, Result: result, OwnerID: ownerID}
}

type User struct {
	ID       int64   `gorm:"primaryKey"`
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

func (p *Point) ToProto() *api.Point {
	return &api.Point{
		Id:     p.ID,
		X:      p.X,
		Y:      p.Y,
		R:      p.R,
		Result: p.Result,
	}
}

func FromProto(proto *api.Point) *Point {
	return &Point{
		ID:      proto.Id,
		X:       proto.X,
		Y:       proto.Y,
		R:       proto.R,
		Result:  proto.Result,
		OwnerID: 0, // Set OwnerID separately if needed
	}
}
