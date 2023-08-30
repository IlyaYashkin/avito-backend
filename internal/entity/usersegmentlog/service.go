package usersegmentlog

import "time"

type UserSegmentLog struct {
	UserId             int32
	SegmentName        string
	Operation          string
	OperationTimestamp time.Time
}
