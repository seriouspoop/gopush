package svc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
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

func (s *Svc) Pull(initial bool) error {
	pullBranch, err := s.bash.GetCurrentBranch()
	if err != nil {
		return err
	}
	remote, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}
	if initial {
		output, err := s.bash.PullBranch(remote.Name, pullBranch, true)
		fmt.Println(output)
		return err
	}

	var providerAuth *config.Credentials
	if remote.AuthMode() == model.AuthHTTP {
		if s.cfg == nil {
			return ErrConfigNotLoaded
		}
		providerAuth = s.cfg.ProviderAuth(remote.Provider())
		if providerAuth == nil {
			return ErrAuthNotFound
		}
	} else if remote.AuthMode() == model.AuthSSH {
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
	pullErr := s.git.Pull(remote, pullBranch, providerAuth)
	for errors.Is(pullErr, ErrInvalidPassphrase) {
		passphrase, err := utils.Prompt(true, false, "invalid passphrase")
		if err != nil {
			return err
		}
		s.passphrase = model.Password(passphrase)
		providerAuth = &config.Credentials{
			Token: s.passphrase.String(),
		}
		pullErr = s.git.Pull(remote, pullBranch, providerAuth)
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
			fmt.Printf("Branch %s already exists. Switching branch...\n", branch.String())
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
	selecttTemplate := &promptui.SelectTemplates{
		Active: "\U0001F892 {{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     "Select commit type",
		Items:     []string{"fix", "feature", "chore", "refactor", "ci"},
		Templates: selecttTemplate,
	}
	_, commitType, err := prompt.Run()
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
	}
	return nil
}

func (s *Svc) Push(setUpstreamBranch bool) (output string, err error) {
	currBranch, err := s.bash.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	if setUpstreamBranch {
		output, err = s.bash.Push(currBranch, true)
		if err != nil {
			fmt.Println(output)
			return "", err
		}
	} else {
		remoteDetails, err := s.git.GetRemoteDetails()
		if err != nil {
			return "", err
		}

		var providerAuth *config.Credentials
		if remoteDetails.AuthMode() == model.AuthHTTP {
			if s.cfg == nil {
				return "", ErrConfigNotLoaded
			}
			providerAuth := s.cfg.ProviderAuth(remoteDetails.Provider())
			if providerAuth == nil {
				return "", ErrAuthNotFound
			}
		} else if remoteDetails.AuthMode() == model.AuthSSH {
			if !s.passphrase.Valid() {
				passphrase, err := utils.Prompt(true, false, "passphrase")
				if err != nil {
					return "", err
				}
				s.passphrase = model.Password(passphrase)
			}
			providerAuth = &config.Credentials{
				Token: s.passphrase.String(),
			}
		} else {
			return "", ErrInvalidAuthMethod
		}
		pushErr := s.git.Push(remoteDetails, currBranch, providerAuth)
		for errors.Is(pushErr, ErrInvalidPassphrase) {
			passphrase, err := utils.Prompt(true, false, "invalid passphrase")
			if err != nil {
				return "", err
			}
			s.passphrase = model.Password(passphrase)
			providerAuth = &config.Credentials{
				Token: s.passphrase.String(),
			}
			pushErr = s.git.Push(remoteDetails, currBranch, providerAuth)
		}
		if pushErr != nil {
			return "", pushErr
		}
		utils.Logger(utils.LOG_SUCCESS, "push successful")
	}
	return "", err
}
