package phabricator

import (
	"encoding/json"
	"fmt"
	"os"
)

type arcRc struct {
	Hosts map[string]struct {
		Token string `json:"token"`
	} `json:"hosts"`
}

func getArcToken(filepath string) (string, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("unable to find .arrcrc file")
	}
	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("unable to open .arrcrc file: %w", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	arcConfig := arcRc{}
	err = decoder.Decode(&arcConfig)
	if err != nil {
		return "", fmt.Errorf("unable to decode .arrcrc file: %w", err)
	}

	// Will return the first one
	for _, conf := range arcConfig.Hosts {
		return conf.Token, nil
	}
	return "", fmt.Errorf("unable to find token")
}
