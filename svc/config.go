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
		remoteName, err := utils.Prompt(false, "remote (default=origin)")
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
		branchPrefix, err := utils.Prompt(false, "branch prefix (default=empty)")
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

	if remoteDetails.AuthMode() == model.AuthHTTP {
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
	} else if remoteDetails.AuthMode() == model.AuthSSH {
		/*
			TODO - 1. check if gopush key present?
			TODO - 2. generate "gopush" public/private key pair else
		*/
	} else {
		return ErrInvalidAuthMethod
	}
	return cfg.Write(configFile, userDir)
}

func (s *Svc) authInput(provider string) (string, string, error) {
	var username, token string
	username, err := utils.Prompt(false, "%s username", provider)
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	token, err = utils.Prompt(true, "%s token", provider)
	if err != nil {
		return "", "", err
	}
	token = strings.TrimSpace(token)

	return username, token, nil
}
