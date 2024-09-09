package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Campaign represents the MongoDB document for the campaign
type Campaign struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Image     string             `bson:"image,omitempty"`
	CTA       string             `bson:"cta,omitempty"`
	Status    string             `bson:"status,omitempty"`
	Targeting Targeting          `bson:"targeting,omitempty"`
}
