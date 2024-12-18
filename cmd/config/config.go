package config

import (
	"fmt"
	_ "fmt"
	"log"
	"os"
	"sync"

	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

var once sync.Once

type App struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Train struct {
	SeatPerSection int `yaml:"seat_per_section"`
	TotalSeats     int `yaml:"total_seats"`
}

type Config struct {
	App   App   `yaml:"app"`
	Train Train `yaml:"train"`
}

func GetConfig() (*Config, error) {
	var config Config
	var e error
	_, b, _, _ := runtime.Caller(0)
	basepath   := filepath.Dir(b)
	fmt.Println(basepath)
	once.Do(func() {

		configFile, err := os.ReadFile(basepath + "/config.yaml")
		if err != nil {
			log.Printf("Config file not found %v : ", err)
			e = err
			return
		}
		er := yaml.Unmarshal(configFile, &config)

		if er != nil {
			log.Printf("Config file not able to unmarshall %v", er)
			e = er
			return
		}
	})

	if e != nil {
		return nil, e
	}
	return &config, nil
}
