package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Campaign represents the MongoDB document for the campaign
type Campaign struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Image     string             `bson:"image"`
	CTA       string             `bson:"cta"`
	Status    string             `bson:"status"`
	Targeting Targeting          `bson:"targeting"`
}
