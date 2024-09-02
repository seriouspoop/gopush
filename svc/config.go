package svc

import (
	"os"
	"strings"

	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
	"github.com/seriouspoop/gopush/utils"
)

const (
	configFile    = ".gopush_config.toml"
	DefaultRemote = "origin"
)

func (s *Svc) LoadConfig() error {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if !s.bash.FileExists(configFile, userDir) {
		return ErrFileNotFound
	}

	s.cfg, err = config.Read(configFile, userDir)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) SetUserPreference() error {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if !s.bash.FileExists(configFile, userDir) {
		_, err := s.bash.CreateFile(configFile, userDir)
		if err != nil {
			return err
		}
	}
	// cfg == nil if file is empty
	cfg, err := config.Read(configFile, userDir)
	if err != nil {
		return err
	}
	utils.Logger(utils.LOG_INFO, "Gathering default settings...")
	if cfg.DefaultRemote == "" {
		remoteName, err := utils.Prompt("remote (default=origin)", false)
		if err != nil {
			return err
		}
		remoteName = strings.TrimSpace(remoteName)
		if remoteName == "" {
			remoteName = DefaultRemote
		}
		cfg.DefaultRemote = remoteName
	}
	if cfg.BranchPrefix == "" {
		branchPrefix, err := utils.Prompt("branch prefix (default=empty)", false)
		if err != nil {
			return err
		}
		branchPrefix = strings.TrimSpace(branchPrefix)
		cfg.BranchPrefix = branchPrefix
	}
	return cfg.Write(configFile, userDir)
}

func (s *Svc) SetRemoteAuth() error {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if !s.bash.FileExists(configFile, userDir) {
		_, err := s.bash.CreateFile(configFile, userDir)
		if err != nil {
			return err
		}
	}
	// cfg == nil if file is empty
	cfg, err := config.Read(configFile, userDir)
	if err != nil {
		return err
	}
	remoteDetails, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}

	utils.Logger(utils.LOG_INFO, "Gathering auth details...")
	provider := remoteDetails.Provider()
	if cfg.ProviderAuth(provider) == nil {
		utils.Logger(utils.LOG_FAILURE, "auth credentials not found")
		username, token, err := s.authInput(provider.String())
		if err != nil {
			return err
		}
		cred := &config.Credentials{
			Username: username,
			Token:    token,
		}
		switch provider {
		case model.ProviderGITHUB:
			cfg.Auth.GitHub = cred
		case model.ProviderBITBUCKET:
			cfg.Auth.BitBucket = cred
		case model.ProviderGITLAB:
			cfg.Auth.GitLab = cred
		}
		utils.Logger(utils.LOG_SUCCESS, "auth generated")
	} else {
		utils.Logger(utils.LOG_SUCCESS, "auth found")
	}
	return cfg.Write(configFile, userDir)
}

func (s *Svc) authInput(provider string) (string, string, error) {
	var username, token string
	username, err := utils.Prompt("%s username", false, provider)
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	token, err = utils.Prompt("%s token", true, provider)
	if err != nil {
		return "", "", err
	}
	token = strings.TrimSpace(token)

	return username, token, nil
}
