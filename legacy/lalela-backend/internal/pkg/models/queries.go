package models

//// Queries
type QueryCreateRequest struct {
	Client string `json:"client"`
	Section string `json:"section"`
	Query string `json:"query"`
	Templates map[string]interface{} `json:"templates"`
}

type QueryCreateResponse struct {
	DataPoints interface{} `json:"datapoints"`
}

