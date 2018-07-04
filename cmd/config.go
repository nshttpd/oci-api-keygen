package cmd

import (
	"encoding/json"
	"io/ioutil"

	"os"

	"strings"

	"fmt"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Version int      `json:"version,omitempty"`
	KeyPath string   `json:"keyPath"`
	ApiKeys []ApiKey `json:"apiKeys"`
}

type ApiKey struct {
	Tenancy     string `json:"tenancy"`
	Fingerprint string `json:"fingerprint"`
}

func loadConfig() *Config {
	c := &Config{}
	if cf, err := ioutil.ReadFile(cfgFile); err != nil {
		if os.IsNotExist(err) {
			p := strings.Split(cfgFile, string(os.PathSeparator))
			c.KeyPath = strings.Join(p[0:len(p)-1], string(os.PathSeparator))
			c.ApiKeys = make([]ApiKey, 0)
			// make sure the directory exists and if not create it
			if _, err := ioutil.ReadDir(c.KeyPath); err != nil {
				err := os.MkdirAll(c.KeyPath, 0755)
				if err != nil {
					log.WithFields(log.Fields{"keyPath": c.KeyPath, "error": err}).Fatal("error creating key dir")
				}
			}
		} else {
			log.Error(err)
		}
	} else {
		err = json.Unmarshal(cf, c)
	}

	return c

}

func (c *Config) SaveConfig() {
	sl := log.WithField("cfgfile", cfgFile)
	if cf, err := json.Marshal(c); err != nil {
		sl.WithField("error", err).Error("error marshalling config file")
	} else {
		if err := ioutil.WriteFile(cfgFile, cf, 0644); err != nil {
			sl.WithField("error", err).Error("error saving config file")
		} else {
			sl.Debug("saved config file")
		}
	}
}

func (c *Config) GetTenancy(tenancy string) (ApiKey, error) {

	for _, a := range c.ApiKeys {
		if a.Tenancy == tenancy {
			return a, nil
		}
	}

	err := fmt.Errorf("invalid tenancy specified")

	return ApiKey{}, err
}
