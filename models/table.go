package models

import "time"

type ConfigureStruct struct {
	ConfigName string    `json:"configName"`
	Docs       string    `json:"docs"`
	Version    string    `bson:"version"`
	Content    string    `json:"content"`
	Created    time.Time `json:"created"`
	Deleted    time.Time `json:"deleted"`
}
