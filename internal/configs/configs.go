package configs

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Server   Server   `yaml:"server"`
	Redis    Redis    `yaml:"redis"`
	Mailer   Mailer   `yaml:"mailer"`
	JWT      JWT      `yaml:"token"`
	Logger   Logger   `yaml:"log"`
}

type (
	JWT struct {
		SecretAuth  string `yaml:"secret"`
		SecretReset string `yaml:"secret_reset"`
	}

	Logger struct {
		LogLevel string `yaml:"level"`
		LogFile  string `yaml:"file"`
	}

	Mailer struct {
		Name              string `yaml:"name"`
		FromEmailAddress  string `yaml:"from_email_address"`
		FromEmailPassword string `yaml:"from_email_password"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int64  `yaml:"DB"`
	}
)

func New() *Config {
	return &Config{
		Server:   Server{},
		Postgres: Postgres{},
		Redis:    Redis{},
		JWT:      JWT{},
		Logger:   Logger{},
	}
}

func MustConfig(p *string) *Config {
	var path string
	if p == nil {
		path = fetchConfigPath()
	} else {
		path = *p
	}

	if _, ok := os.Stat(path); os.IsNotExist(ok) {
		panic("Config file does not exist: " + path)
	}

	cfg := New()

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "./config.yaml", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
