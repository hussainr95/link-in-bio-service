package entity

import "time"

// Visit represents a single visit for analytics.
type Visit struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	LinkID    string    `json:"linkId" bson:"linkId"`
	VisitedAt time.Time `json:"visitedAt" bson:"visitedAt"`
}
