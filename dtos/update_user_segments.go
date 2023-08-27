package dtos

type UpdateUserSegments struct {
	UserId         int32    `json:"user_id" binding:"required"`
	AddSegments    []any    `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
}
