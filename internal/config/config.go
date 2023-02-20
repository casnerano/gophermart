package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

const (
	DefaultYAMLConfigPath = "./configs/gophermart.yml"
)

type AppEnv int

const (
	AppEnvDev AppEnv = iota
	AppEnvTest
	AppEnvStage
	AppEnvProd
)

func (ae *AppEnv) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "prod":
		*ae = AppEnvProd
	case "stage":
		*ae = AppEnvStage
	case "test":
		*ae = AppEnvTest
	case "dev":
		*ae = AppEnvDev
	default:
		return fmt.Errorf("invalid value \"%s\" for type config.AppEnv", text)
	}
	return nil
}

func (ae *AppEnv) UnmarshalYAML(value *yaml.Node) error {
	var strEnvVal string
	if err := value.Decode(&strEnvVal); err != nil {
		return err
	}
	return ae.UnmarshalText([]byte(strEnvVal))
}

type Configuration struct {
	App struct {
		ENV    AppEnv `yaml:"env" env:"APP_ENV"`
		Secret string `yaml:"secret" env:"APP_SECRET"`
	} `yaml:"app"`
	Server struct {
		Address string `yaml:"address" env:"RUN_ADDRESS"`
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn" env:"DATABASE_URI"`
	} `yaml:"database"`
	Accrual struct {
		Service struct {
			Address string `yaml:"address" env:"ACCRUAL_SYSTEM_ADDRESS"`
		} `yaml:"service"`
		Queue struct {
			DSN string `yaml:"dsn"`
		} `yaml:"queue"`
		PoolInterval int `yaml:"pool_interval"`
	} `yaml:"accrual"`
}

func New() (*Configuration, error) {
	config := Configuration{}
	if err := config.loadFromYAML(DefaultYAMLConfigPath); err != nil {
		return nil, err
	}

	if err := config.loadFromEnv(); err != nil {
		return nil, err
	}

	if err := config.loadFromFlags(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Загрузка из файла в формате YAML
func (c *Configuration) loadFromYAML(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		return err
	}

	return nil
}

// Загрузка из переменных окружения
func (c *Configuration) loadFromEnv() error {
	return env.Parse(c)
}

// Загрузка из флагов запуска приложения
func (c *Configuration) loadFromFlags() error {
	flag.StringVar(&c.Server.Address, "a", c.Server.Address, "Server address")
	flag.StringVar(&c.Database.DSN, "d", c.Database.DSN, "Database data source name")
	flag.StringVar(&c.Accrual.Service.Address, "r", c.Accrual.Service.Address, "Accrual service address")
	flag.StringVar(&c.Accrual.Queue.DSN, "q", c.Accrual.Queue.DSN, "Accrual queue data source name")

	flag.Parse()

	return nil
}
