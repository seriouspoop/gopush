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
	if !s.bash.Exists(userDir, gopushDir) {
		err := s.bash.CreateDir(userDir, gopushDir)
		if err != nil {
			return "", err
		}
	}

	gopushDirPath := filepath.Join(userDir, gopushDir)

	if !s.bash.Exists(gopushDirPath, configFile) {
		_, err := s.bash.CreateFile(gopushDirPath, configFile)
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
	if !s.bash.Exists(gopushDirPath, configFile) {
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
	} else {
		utils.Logger(utils.LOG_SUCCESS, "remote name found")
	}
	if cfg.BranchPrefix == "" {
		branchPrefix, err := utils.Prompt(false, true, "branch prefix (default=empty)")
		if err != nil {
			return err
		}
		branchPrefix = strings.TrimSpace(branchPrefix)
		cfg.BranchPrefix = branchPrefix
	} else {
		utils.Logger(utils.LOG_SUCCESS, "branch prefix found")
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
	if !s.bash.Exists(gopushDirPath, keyName) {
		// generate ssh key pair
		mail, err := utils.Prompt(false, false, "mail")
		if err != nil {
			return err
		}
		passphrase, err := utils.Prompt(true, false, "passphrase")
		if err != nil {
			return err
		}
		err = s.bash.GenerateSSHKey(gopushDirPath, keyName, mail, passphrase)
		if err != nil {
			return err
		}
		utils.Logger(utils.LOG_SUCCESS, "keys generated")
		// add keys to known hosts
		hostCode := fmt.Sprintf("Host %s\n  AddKeysToAgent yes\n  IdentityFile \"%s\"", remoteDetails.Provider().HostURL(), filepath.Join(gopushDirPath, keyName))
		fileContent, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
		if err != nil {
			return err
		}
		fileContent = []byte(strings.Join([]string{string(fileContent), hostCode}, "\n"))
		err = os.Remove(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(os.Getenv("HOME"), ".ssh", "config"), fileContent, 0777)
		if err != nil {
			return err
		}
		utils.Logger(utils.LOG_SUCCESS, "key added to host")
		message := fmt.Sprintf("copy contents of %s.pub and upload the keys on %s", filepath.Join(gopushDirPath, keyName), remoteDetails.Provider().String())
		utils.Logger(utils.LOG_STRICT_INFO, message)
		return ErrWaitExit
	} else {
		utils.Logger(utils.LOG_SUCCESS, "key found")
	}
	return err
}
