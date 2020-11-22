package config

//Config struct
type Config struct {
	DB *DBConfig
}

//DBConfig struct
type DBConfig struct {
	Username string
	Password string
	Database string
	Port     int
	Host     string
}

//GetConfig returns a db configuration
func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Username: "postgres",
			Password: "cristiano1994",
			Database: "spataro_visitas",
			Port:     5432,
			Host:     "localhost",
		},
	}
}
