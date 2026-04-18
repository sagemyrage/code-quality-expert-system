package config

type AppConfig struct {
	Port string
}

type PostgresConfig struct {
	DB       string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

type RedisConfig struct {
	DB       int
	Password string
	Host     string
	Port     string
}

type SessionConfig struct {
	Secret string
}

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
}
