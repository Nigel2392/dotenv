package dotenv_test

import (
	"testing"

	"github.com/Nigel2392/dotenv"
)

type TestConfig struct {
	Host    string       `env:"Host"`
	Port    int          `env:"Port"`
	TEST    string       `env:"TEST"`
	BOOLEAN bool         `env:"BOOLEAN"`
	NAMES   []string     `env:"NAMES"`
	Inner   *InnerConfig `env:"INNER"`
}

type InnerConfig struct {
	DEFAULT  string `env:"DEFAULT"`
	TELEGRAM string `env:"TELEGRAM"`
	//Test     *TestConfig `env:"TESTCONFIG"`
}

func TestLoad(t *testing.T) {
	var env = (&dotenv.Env{Variables: map[string][]string{}})
	env.LoadString(`
		TESTCONFIG.Host = "localhost"
		TESTCONFIG.Port = 8080
		TESTCONFIG.TEST = "test"
		TESTCONFIG.BOOLEAN = true
		TESTCONFIG.NAMES = "John", "Doe"
		INNERCONFIG.DEFAULT = "default"
		INNERCONFIG.TELEGRAM = "telegram"
	`)
	var config TestConfig
	var inner InnerConfig
	env.Unmarshal(&config, &inner)
	if config.Host != "localhost" {
		t.Error("Host should be localhost")
	}
	if config.Port != 8080 {
		t.Error("Port should be 8080")
	}
	if config.TEST != "test" {
		t.Error("TEST should be test")
	}
	if config.BOOLEAN != true {
		t.Error("BOOLEAN should be true")
	}
	if len(config.NAMES) != 2 {
		t.Error("NAMES should be an array of 2")
	}
	if config.NAMES[0] != "John" {
		t.Error("NAMES[0] should be John")
	}
	if config.NAMES[1] != "Doe" {
		t.Error("NAMES[1] should be Doe")
	}
	if config.Inner != nil {
		t.Error("Inner should be nil")
	}
	if inner.DEFAULT != "default" {
		t.Error("DEFAULT should be default")
	}
	if inner.TELEGRAM != "telegram" {
		t.Error("TELEGRAM should be telegram")
	}
	t.Logf("%+v", config)
	t.Logf("%+v", inner)
}
