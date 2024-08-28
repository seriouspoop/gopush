package svc

import (
	"fmt"
	"os"
	"strings"

	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
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
	fmt.Println("Gathering default setting...")
	if cfg.DefaultRemote == "" {
		fmt.Print("-  Remote (default=origin): ")
		remoteName, err := s.r.ReadString('\n')
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
		fmt.Print(`-  Branch Prefix (default=empty): `)
		branchPrefix, err := s.r.ReadString('\n')
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

	provider := remoteDetails.Provider()
	if cfg.ProviderAuth(provider) == nil {
		fmt.Println("\nAuth credentials not found.")
		fmt.Println("Gathering auth details...")
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
	}
	return cfg.Write(configFile, userDir)
}

func (s *Svc) authInput(provider string) (string, string, error) {
	var username, token string
	fmt.Printf("-  %s Username: ", provider)
	username, err := s.r.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)
	fmt.Printf("-  %s Token: ", provider)
	token, err = s.r.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	token = strings.TrimSpace(token)
	return username, token, nil
}
