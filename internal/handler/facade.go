package handler

import "github.com/seriouspoop/gopush/model"

type servicer interface {
	LoadProject() error
	InitializeRepo() error
	InitializeRemote() error
	SetUserPreference() error
	SetRemoteAuth() error
	LoadConfig() error
	// FetchAndMerge() error
	Pull(initial bool) error
	StageChanges() error
	SwitchBranchIfExists(branch model.Branch) (bool, error)
	CreateBranchAndSwitch(branch model.Branch) error
	CheckTestsAndRun() (bool, error)
	Push(setUpstreamBranch bool) (output string, err error)
}
