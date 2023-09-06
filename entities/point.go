package entities

type PointKey struct {
	PK string `dynamo:",hash" json:"pk"`  // Hash key, a.k.a. partition key
	SK string `dynamo:",range" json:"sk"` // Range key, a.k.a. sort key
}

// Use struct tags much like the standard JSON library,
// you can embed anonymous structs too!
type Point struct {
	PointKey
	Coordinates []string `dynamo:",set" json:"coordinates"`
	SubRegion   string   `json:"subRegion"`
	Region      string   `json:"region"`
	Address     string   `json:"address"`
	ResourceUrl string   `json:"resourceUrl,omitempty"`
	Count       int      `json:"count"`
}
