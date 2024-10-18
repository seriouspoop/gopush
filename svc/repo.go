package svc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
	"github.com/seriouspoop/gopush/utils"
)

func (s *Svc) LoadProject() error {
	err := s.git.GetRepo()
	if err != nil {
		return err
	}

	defaultRemote := DefaultRemote
	if s.cfg != nil {
		defaultRemote = s.cfg.DefaultRemote
	}
	err = s.git.LoadRemote(defaultRemote)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) InitializeRepo() error {
	utils.Logger(utils.LOG_INFO, "Initializing repository...")
	err := s.git.CreateRepo()
	if err != nil {
		return err
	}

	utils.Logger(utils.LOG_SUCCESS, "repository initialized")
	return nil
}

func (s *Svc) InitializeRemote() error {
	if s.cfg == nil {
		return ErrConfigNotLoaded
	}
	var remoteURL string
	utils.Logger(utils.LOG_INFO, "Adding remote...")

	remoteURL, err := utils.Prompt(false, false, "remote url")
	if err != nil {
		return err
	}

	remoteURL = strings.TrimSpace(remoteURL)
	remote := &model.Remote{
		Name: s.cfg.DefaultRemote,
		Url:  remoteURL,
	}
	err = s.git.AddRemote(remote)
	if err != nil {
		return err
	}

	utils.Logger(utils.LOG_SUCCESS, "remote added")
	return nil
}

func (s *Svc) Pull(force bool) error {
	pullBranch, err := s.bash.GetCurrentBranch()
	if err != nil {
		return err
	}
	remoteDetails, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}

	var providerAuth *config.Credentials
	if remoteDetails.AuthMode() == model.AuthHTTP {
		if s.cfg == nil {
			return ErrConfigNotLoaded
		}
		providerAuth = s.cfg.ProviderAuth(remoteDetails.Provider())
		if providerAuth == nil {
			return ErrAuthNotFound
		}
	} else if remoteDetails.AuthMode() == model.AuthSSH {
		if !s.passphrase.Valid() {
			passphrase, err := utils.Prompt(true, false, "passphrase")
			if err != nil {
				return err
			}
			s.passphrase = model.Password(passphrase)
		}
		providerAuth = &config.Credentials{
			Token: s.passphrase.String(),
		}
	} else {
		return ErrInvalidAuthMethod
	}
	pullErr := s.git.Pull(remoteDetails, pullBranch, providerAuth, force)
	if pullErr == nil {
		utils.Logger(utils.LOG_SUCCESS, "changes pulled")
	}
	for errors.Is(pullErr, ErrInvalidPassphrase) {
		passphrase, err := utils.Prompt(true, false, "invalid passphrase")
		if err != nil {
			return err
		}
		s.passphrase = model.Password(passphrase)
		providerAuth = &config.Credentials{
			Token: s.passphrase.String(),
		}
		pullErr = s.git.Pull(remoteDetails, pullBranch, providerAuth, force)
	}
	if errors.Is(pullErr, ErrKeyNotSupported) {
		message := fmt.Sprintf("copy contents of %s.pub and upload the keys on %s", filepath.Join(os.Getenv("HOME"), gopushDir, keyName), remoteDetails.Provider().String())
		utils.Logger(utils.LOG_STRICT_INFO, message)
	} else if errors.Is(pullErr, ErrAlreadyUpToDate) {
		utils.Logger(utils.LOG_SUCCESS, "already up-to-date")
		return nil
	} else if errors.Is(pullErr, ErrMergeFailed) {
		output, err := s.bash.PullMerge()
		if err != nil {
			utils.Logger(utils.LOG_FAILURE, output)
			return err
		}
		utils.Logger(utils.LOG_SUCCESS, "changes merged")
		return nil
	}
	return pullErr
}

func (s *Svc) SwitchBranchIfExists(branch model.Branch) (bool, error) {
	branches, err := s.git.GetBranchNames()
	if err != nil {
		return false, err
	}
	for _, br := range branches {
		if br.String() == branch.String() {
			err = s.git.CheckoutBranch(branch)
			return true, err
		}
	}
	return false, nil
}

func (s *Svc) CreateBranchAndSwitch(branch model.Branch) error {
	err := s.git.CreateBranch(branch)
	if err != nil {
		return err
	}
	return s.git.CheckoutBranch(branch)
}

func generateCommitMsg() (string, error) {
	items := []string{"fix", "feature", "chore", "refactor", "ci"}
	commitType, err := utils.Select(items)
	if err != nil {
		return "", err
	}

	shortner := map[string]string{
		"refactor": "ref",
		"feature":  "feat",
	}

	if _, ok := shortner[commitType]; ok {
		commitType = shortner[commitType]
	}

	msg, err := utils.Prompt(false, false, "commit message")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s: %s", commitType, msg), nil
}

func (s *Svc) StageChanges() error {
	change, err := s.git.ChangeOccured()
	if err != nil {
		return err
	}
	if change {
		commitMsg, err := generateCommitMsg()
		if err != nil {
			return err
		}
		err = s.git.AddThenCommit(commitMsg)
		if err != nil {
			return err
		}
		utils.Logger(utils.LOG_SUCCESS, "files added")
	} else {
		utils.Logger(utils.LOG_SUCCESS, "no files changed")
	}
	return nil
}

func (s *Svc) Push(force bool) error {
	currBranch, err := s.bash.GetCurrentBranch()
	if err != nil {
		return err
	}
	remoteDetails, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}

	var providerAuth *config.Credentials
	if remoteDetails.AuthMode() == model.AuthHTTP {
		if s.cfg == nil {
			return ErrConfigNotLoaded
		}
		providerAuth = s.cfg.ProviderAuth(remoteDetails.Provider())
		if providerAuth == nil {
			return ErrAuthNotFound
		}
	} else if remoteDetails.AuthMode() == model.AuthSSH {
		if !s.passphrase.Valid() {
			passphrase, err := utils.Prompt(true, false, "passphrase")
			if err != nil {
				return err
			}
			s.passphrase = model.Password(passphrase)
		}
		providerAuth = &config.Credentials{
			Token: s.passphrase.String(),
		}
	} else {
		return ErrInvalidAuthMethod
	}

	if providerAuth == nil {
		return ErrAuthLoadFailed
	}
	pushErr := s.git.Push(remoteDetails, currBranch, providerAuth, force)
	for errors.Is(pushErr, ErrInvalidPassphrase) {
		passphrase, err := utils.Prompt(true, false, "invalid passphrase")
		if err != nil {
			return err
		}
		s.passphrase = model.Password(passphrase)
		providerAuth = &config.Credentials{
			Token: s.passphrase.String(),
		}
		pushErr = s.git.Push(remoteDetails, currBranch, providerAuth, force)
	}
	if pushErr != nil {
		if errors.Is(pushErr, ErrKeyNotSupported) {
			message := fmt.Sprintf("copy contents of %s.pub and upload the keys on %s", filepath.Join(os.Getenv("HOME"), gopushDir, keyName), remoteDetails.Provider().String())
			utils.Logger(utils.LOG_STRICT_INFO, message)
		}
		if errors.Is(pushErr, ErrAlreadyUpToDate) {
			utils.Logger(utils.LOG_SUCCESS, "already up-to-date")
			return nil
		}
		return pushErr
	}
	utils.Logger(utils.LOG_SUCCESS, "push successful")
	return err
}
