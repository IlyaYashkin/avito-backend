package usersegmentlog

import "time"

type UserSegmentLog struct {
	UserId             int32
	SegmentName        string
	Operation          string
	OperationTimestamp time.Time
}

const LOG_OPERATION_ADD = "added"
const LOG_OPERATION_ADD_PERCENTAGE = "added by percentage"
const LOG_OPERATION_DELETE = "deleted"
const LOG_OPERATION_DELETE_TTL = "deleted by ttl"
