package config

type Config struct {
	Env string    ``
	Net *Network  `yaml:"network"`
	Db  *Database `yaml:"database"`
}

type Network struct {
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

type Database struct {
	Driver   string `yaml:"driver" env:"DB_DRIVER" env-default:"postgres"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
}

const DEFAULT_CONFIG_PATH string = "~/lesson-manager/config.json"
