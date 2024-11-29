package service

import (
	"awesoma31/common/api"
	"awesoma31/common/storage/model"
	"context"
	"fmt"
	"github.com/awesoma31/points-service/storage"
	"math"
)

type PointsService interface {
	AddPoint(ctx context.Context, point *api.PointData, id int64) (*model.Point, error)
	GetPointsPageByID(ctx context.Context, pageParam string, size int32, id int64) (*api.PointsPage, error)
	GetTotalPointsByID(ctx context.Context, id int64) (int, error)
}

type pointsService struct {
	store storage.PointsStore
}

func NewPointsService(store storage.PointsStore) PointsService {
	return &pointsService{store: store}
}

func (s *pointsService) AddPoint(ctx context.Context, pointData *api.PointData, id int64) (*model.Point, error) {
	point := &model.Point{
		X:       pointData.X,
		Y:       pointData.Y,
		R:       pointData.R,
		Result:  getResult(pointData),
		OwnerID: id,
	}

	savedPoint, err := s.store.Create(ctx, point)
	if err != nil {
		return nil, fmt.Errorf("failed to add pointData: %w", err)
	}

	return savedPoint, nil
}

func (s *pointsService) GetPointsPageByID(ctx context.Context, pageParam string, size int32, userId int64) (*api.PointsPage, error) {
	totalPoints, err := s.store.GetTotalPointsByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get total points: %w", err)
	}

	totalPages := int32(math.Ceil(float64(totalPoints) / float64(size)))
	if totalPages == 0 {
		totalPages = 1
	}

	var pageNumber int32
	if pageParam == "last" {
		pageNumber = totalPages - 1
	} else {
		parsedPage, err := parsePageParam(pageParam)
		if err != nil {
			return nil, fmt.Errorf("invalid page parameter: %w", err)
		}
		pageNumber = clamp(parsedPage-1, 0, totalPages-1)
	}

	// Get points for the requested page
	points, err := s.store.GetPointsByUserIDWithPagination(ctx, userId, int(size), int(pageNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get points: %w", err)
	}

	// Convert the database points to Protobuf points
	var pbPoints []*api.Point
	for _, p := range points {
		pbPoints = append(pbPoints, &api.Point{
			Id:     p.ID,
			X:      p.X,
			Y:      p.Y,
			R:      p.R,
			Result: p.Result,
		})
	}

	// Return the Protobuf PointsPage
	return &api.PointsPage{
		Content:       pbPoints,
		PageNumber:    pageNumber + 1,
		PageSize:      size,
		TotalElements: int64(totalPoints),
		TotalPages:    totalPages,
	}, nil
}

func (s *pointsService) GetTotalPointsByID(ctx context.Context, id int64) (int, error) {
	totalPoints, err := s.store.GetTotalPointsByUserID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to get total points: %w", err)
	}
	return totalPoints, nil
}

// Utility to parse page parameter
func parsePageParam(pageParam string) (int32, error) {
	var page int32
	_, err := fmt.Sscanf(pageParam, "%d", &page)
	return page, err
}

// Utility to clamp a value within a range
func clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func getResult(p *api.PointData) bool {
	x := p.X
	y := p.Y
	r := p.R

	if y > 0 && x < 0 {
		return false
	}
	if y < 0 && x < 0 {
		if (y*y + x*x) > (r/2)*(r/2) {
			return false
		}
	}
	if y > 0 && x > 0 {
		if y > ((-1.0/2.0)*x + r/2) {
			return false
		}
	}
	if x > 0 && y < 0 {
		return !(x > r) && !(y < -r/2)
	}
	return true
}
