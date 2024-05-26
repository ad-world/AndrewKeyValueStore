package akv

import (
	"time"
)

type Key string
type Value struct {
	Value string 		     `json:"value"`
	LastUpdated time.Time	 `json:"last_updated"`
}

type DeleteRequest struct {
	Key string
}

type GetRequest struct {
	Key string
}

type PutRequest struct {
	Key string
	Value string
}

type GetLastUpdatedRequest struct {
	Key string
}