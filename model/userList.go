package model

import "time"

type UserList struct {
	Name       string
	Online     bool
	AppAuth    string
	ParaAuth   string
	HostName   string
	Expiration time.Time
}
