package dto

type UpdateSegment struct {
	Name string `json:"name" binding:"required"`
}
