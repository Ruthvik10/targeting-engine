package model

type Targeting struct {
	IncludeApp     []string `bson:"includeApp"`
	ExcludeApp     []string `bson:"excludeApp"`
	IncludeOS      []string `bson:"includeOS"`
	ExcludeOS      []string `bson:"excludeOS"`
	IncludeCountry []string `bson:"includeCountry"`
	ExcludeCountry []string `bson:"excludeCountry"`
}
