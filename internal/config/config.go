package config

// APIConfig holds global configuration values for the API handlers
type APIConfig struct {
	JWTSecret string
	PolkaKey  string
	Platform  string
}
