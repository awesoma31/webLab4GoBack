package docs

// PointsPage represents the paginated list of points (for Swagger docs)
type PointsPage struct {
	Content       []Point `json:"content"`
	PageNumber    int32   `json:"pageNumber"`
	PageSize      int32   `json:"pageSize"`
	TotalElements int64   `json:"totalElements"`
	TotalPages    int32   `json:"totalPages"`
}

// Point represents a single point (for Swagger docs)
type Point struct {
	ID     int64   `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	R      float64 `json:"r"`
	Result bool    `json:"result"`
}
