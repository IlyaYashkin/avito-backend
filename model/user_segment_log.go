package model

import "time"

type UserSegmentLog struct {
	Id                 int32
	UserId             int32
	SegmentName        string
	Operation          string
	OperationTimestamp time.Time
}
