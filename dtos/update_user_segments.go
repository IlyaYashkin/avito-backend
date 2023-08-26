package dtos

type UpdateUserSegments struct {
	UserId         int32    `json:"user_id"`
	AddSegments    []string `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
}
