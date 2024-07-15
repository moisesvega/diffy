package client

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
)

type arcRc struct {
	Hosts map[string]struct {
		Token string `json:"token"`
	} `json:"hosts"`
}

func getArcToken() (string, error) {
	// TODO: make it more easy to test this.
	cUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get current user: %w", err)
	}

	filePath := path.Join(cUser.HomeDir, ".arcrc")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("unable to find .arrcrc file")
	}
	f, err := os.Open(filePath)
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

	cfg, ok := arcConfig.Hosts[cUser.HomeDir]
	if !ok {
		return "", fmt.Errorf("unable to find .arrcrc file")
	}
	return cfg.Token, nil
}
