package svc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
	"github.com/seriouspoop/gopush/utils"
)

const (
	configFile    = "gopush_config.toml"
	DefaultRemote = "origin"
	gopushDir     = ".gopush"
	keyName       = "gopush_key"
)

func (s *Svc) createConfigPath() (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if !s.bash.Exists(gopushDir, userDir) {
		err := s.bash.CreateDir(gopushDir, userDir)
		if err != nil {
			return "", err
		}
	}

	gopushDirPath := filepath.Join(userDir, gopushDir)

	if !s.bash.Exists(configFile, gopushDirPath) {
		_, err := s.bash.CreateFile(configFile, gopushDirPath)
		if err != nil {
			return "", err
		}
	}
	return gopushDirPath, nil
}

func (s *Svc) LoadConfig() error {
	gopushDirPath, err := s.createConfigPath()
	if err != nil {
		return err
	}
	if !s.bash.Exists(configFile, gopushDirPath) {
		return ErrFileNotFound
	}

	s.cfg, err = config.Read(configFile, gopushDirPath)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) SetUserPreference() error {
	gopushDirPath, err := s.createConfigPath()
	if err != nil {
		return err
	}
	// cfg == nil if file is empty
	cfg, err := config.Read(configFile, gopushDirPath)
	if err != nil {
		return err
	}
	utils.Logger(utils.LOG_INFO, "Gathering default settings...")
	if cfg.DefaultRemote == "" {
		remoteName, err := utils.Prompt(false, true, "remote (default=origin)")
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
		branchPrefix, err := utils.Prompt(false, true, "branch prefix (default=empty)")
		if err != nil {
			return err
		}
		branchPrefix = strings.TrimSpace(branchPrefix)
		cfg.BranchPrefix = branchPrefix
	}
	return cfg.Write(configFile, gopushDirPath)
}

func (s *Svc) SetRemoteHTTPAuth() error {
	remoteDetails, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}

	if remoteDetails.AuthMode() != model.AuthHTTP {
		return ErrInvalidAuthMethod
	}

	gopushDirPath, err := s.createConfigPath()
	if err != nil {
		return err
	}
	// cfg == nil if file is empty
	cfg, err := config.Read(configFile, gopushDirPath)
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
	return cfg.Write(configFile, gopushDirPath)
}

func (s *Svc) authInput(provider string) (string, string, error) {
	var username, token string
	username, err := utils.Prompt(false, false, "%s username", provider)
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	token, err = utils.Prompt(true, false, "%s token", provider)
	if err != nil {
		return "", "", err
	}
	token = strings.TrimSpace(token)

	return username, token, nil
}

func (s *Svc) SetRemoteSSHAuth() error {
	remoteDetails, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}

	if remoteDetails.AuthMode() != model.AuthSSH {
		return ErrInvalidAuthMethod
	}

	gopushDirPath, err := s.createConfigPath()
	if err != nil {
		return err
	}

	utils.Logger(utils.LOG_INFO, "Gathering ssh keys...")
	if !s.bash.Exists(keyName, gopushDirPath) {
		// generate ssh key pair
		mail, err := utils.Prompt(false, false, "mail")
		if err != nil {
			return err
		}
		passphrase, err := utils.Prompt(true, false, "passphrase")
		if err != nil {
			return err
		}
		err = s.bash.GenerateSSHKey(keyName, gopushDirPath, mail, passphrase)
		if err != nil {
			return err
		}
		utils.Logger(utils.LOG_SUCCESS, "keys generated")
		message := fmt.Sprintf("copy contents of %s.pub and upload the keys on %s", filepath.Join(gopushDirPath, keyName), remoteDetails.Provider().String())
		utils.Logger(utils.LOG_STRICT_INFO, message)
		return ErrWaitExit
	} else {
		utils.Logger(utils.LOG_SUCCESS, "key found")
	}
	return err
}
