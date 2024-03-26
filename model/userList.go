package model

import "time"

type UserList struct {
	Name           string
	Online         string
	AppAuth        string
	ParaAuth       string
	HostName       string
	Expiration     time.Time
	SystemStatus   string
	FCCurr         int
	FCCoolOutT     int
	FirstErrCode   int
	StartCount     int
	RunHours       int
	ECUSWName      string
	LabviewVersion int
	LabviewBranch  int
}
