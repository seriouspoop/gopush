package svc

import (
	"fmt"
	"strings"

	"github.com/seriouspoop/gopush/model"
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
	fmt.Println("✅ Repository initialized.")
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
	fmt.Println("✅ Remote added.")
	return nil
}

func (s *Svc) Pull(force bool) error {
	currBranch, err := s.bash.GetCurrentBranch()
	if err != nil {
		return err
	}
	remote, err := s.git.GetRemoteDetails()
	if err != nil {
		return err
	}
	if force {
		_, err = s.bash.PullBranch(remote.Name, currBranch, true)
		return err
	}
	if s.cfg == nil {
		return ErrConfigNotLoaded
	}
	return s.git.Pull(remote, currBranch, s.cfg.ProviderAuth(remote.Provider()))
}

func (s *Svc) StageChanges() error {
	return s.git.AddThenCommit()
}
