package config

import (
	"fmt"
	"github.com/prometheus/common/log"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	Clusters map[string]ClusterConfig `yaml:"clusters"`
}

type SafeConfig struct {
	sync.RWMutex
	C *Config
}

type ClusterConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (sc *SafeConfig) ReloadConfig(configFile string) error {
	var c = &Config{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf("Error reading config file: %s", err)
		return err
	}
	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		log.Errorf("Error parsing config file: %s", err)
		return err
	}

	sc.Lock()
	sc.C = c
	sc.Unlock()

	log.Infoln("Loaded config file")
	return nil
}

func (sc *SafeConfig) ClusterConfigForTarget(target string) (*ClusterConfig, error) {
	sc.Lock()
	defer sc.Unlock()
	if clusterConfig, ok := sc.C.Clusters[target]; ok {
		return &ClusterConfig{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		}, nil
	}
	if clusterConfig, ok := sc.C.Clusters["default"]; ok {
		return &ClusterConfig{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		}, nil
	}
	return nil, fmt.Errorf("no credentials found for target %s", target)
}
