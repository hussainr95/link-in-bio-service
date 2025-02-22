package entity

import "time"

// Link represents the data model for a bio link.
type Link struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	URL       string    `json:"url" bson:"url"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
	Clicks    int       `json:"clicks" bson:"clicks"`
}
