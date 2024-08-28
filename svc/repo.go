package svc

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
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
	fmt.Println("Initializing repository...")
	err := s.git.CreateRepo()
	if err != nil {
		return err
	}

	utils.Logger(utils.STATUS_SUCCESS, "Repository initialized")
	return nil
}

func (s *Svc) InitializeRemote() error {
	if s.cfg == nil {
		return ErrConfigNotLoaded
	}
	var remoteURL string
	fmt.Println("Adding remote...")
	fmt.Print("-  Enter Remote URL: ")
	remoteURL, err := s.r.ReadString('\n')
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
	fmt.Println("\U00002714 Remote added.")
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
	if s.cfg == nil {
		return ErrConfigNotLoaded
	}
	providerAuth := s.cfg.ProviderAuth(remote.Provider())
	if providerAuth == nil {
		return ErrAuthNotFound
	}
	return s.git.Pull(remote, pullBranch, remote.AuthMode(), providerAuth)
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

	message := promptui.Prompt{
		Label: "Commit Message",
	}
	msg, err := message.Run()
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
		utils.Logger(utils.STATUS_SUCCESS, "Files added")
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
		if s.cfg == nil {
			return "", ErrConfigNotLoaded
		}
		providerAuth := s.cfg.ProviderAuth(remoteDetails.Provider())
		if providerAuth == nil {
			return "", ErrAuthNotFound
		}
		err = s.git.Push(remoteDetails, currBranch, remoteDetails.AuthMode(), providerAuth)
		if err != nil {
			return "", err
		}
	}
	utils.Logger(utils.STATUS_SUCCESS, "Push Successful")
	return
}
