package config
type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername any    `json:"current_user_name"`
}
