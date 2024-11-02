package entities

type HealthRes struct {
	Status       string `json:"status"`
	IsProduction bool   `json:"is_production"`
	Region       string `json:"region"`
}
