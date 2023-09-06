package schemas

type CreatePointSchema struct {
	Coordinates        []float64 `json:"coordinates"`
	Base64PotholeImage string    `json:"base64PotholeImage"`
	Count              int       `json:"count"`
}

type UpdatePointSchema struct {
	Coordinates        []float64 `json:"coordinates"`
	Base64PotholeImage string    `json:"base64PotholeImage"`
	Count              int       `json:"count"`
}
