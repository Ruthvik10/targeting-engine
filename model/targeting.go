package model

type Targeting struct {
	IncludeApp     []string `bson:"includeApp,omitempty"`
	ExcludeApp     []string `bson:"excludeApp,omitempty"`
	IncludeOS      []string `bson:"includeOS,omitempty"`
	ExcludeOS      []string `bson:"excludeOS,omitempty"`
	IncludeCountry []string `bson:"includeCountry,omitempty"`
	ExcludeCountry []string `bson:"excludeCountry,omitempty"`
}
