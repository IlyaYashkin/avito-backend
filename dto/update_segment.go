package dto

type UpdateSegment struct {
	Name           string  `json:"name" binding:"required"`
	UserPercentage float32 `json:"user_percentage"`
}
