package svc

import (
	"os"

	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
)

type gitHelper interface {
	GetRepo() error
	CreateRepo() error
	// Fetch() error
	// Merge(remoteName string, branchName model.Branch) error
	CreateBranch(name model.Branch) error
	GetBranchNames() ([]model.Branch, error)
	CheckoutBranch(name model.Branch) error
	AddRemote(remote *model.Remote) error
	LoadRemote(remoteName string) error
	GetRemoteDetails() (*model.Remote, error)
	ChangeOccured() (bool, error)
	AddThenCommit() error
	Pull(remote *model.Remote, branch model.Branch, authType model.AuthMode, auth *config.Credentials) error
	Push(remote *model.Remote, branch model.Branch, authType model.AuthMode, auth *config.Credentials) error
}

type scriptHelper interface {
	GetCurrentBranch() (model.Branch, error)
	PullBranch(remoteName string, branch model.Branch, force bool) (string, error)
	GenerateMocks() (string, error)
	TestsPresent() (bool, error)
	RunTests() (string, error)
	Push(branch model.Branch, withUpStream bool) (string, error)
	FileExists(filename, path string) bool
	CreateFile(filename, path string) (*os.File, error)
	SetUpstream(remoteName string, branch model.Branch) error
}
