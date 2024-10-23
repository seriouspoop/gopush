package svc

import (
	"os"

	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
)

type gitHelper interface {
	GetRepo() error
	CreateRepo() error
	CreateBranch(name model.Branch) error
	GetBranchNames() ([]model.Branch, error)
	CheckoutBranch(name model.Branch) error
	AddRemote(remote *model.Remote) error
	LoadRemote(remoteName string) error
	GetRemoteDetails() (*model.Remote, error)
	ChangeOccured() (bool, error)
	AddThenCommit(commitMsg string) error
	Pull(remote *model.Remote, branch model.Branch, auth *config.Credentials, force bool) error
	Push(remote *model.Remote, branch model.Branch, auth *config.Credentials, force bool) error
}

type scriptHelper interface {
	GetCurrentBranch() (model.Branch, error)
	GenerateMocks() (string, error)
	TestsPresent() (bool, error)
	RunTests() (string, error)
	Exists(path, name string) bool
	CreateFile(path, name string) (*os.File, error)
	CreateDir(path, name string) error
	GenerateSSHKey(path, keyName, mail, passphrase string) error

	PullMerge() (string, error)
}
