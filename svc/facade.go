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
	Push(remote *model.Remote, branch model.Branch, auth *config.Credentials) error
}

type scriptHelper interface {
	GetCurrentBranch() (model.Branch, error)
	PullBranch(remoteName string, branch model.Branch, force bool) (string, error)
	GenerateMocks() (string, error)
	TestsPresent() (bool, error)
	RunTests() (string, error)
	Push(branch model.Branch, withUpStream bool) (string, error)
	Exists(path, name string) bool
	CreateFile(path, name string) (*os.File, error)
	CreateDir(path, name string) error
	SetUpstream(remoteName string, branch model.Branch) error
	GenerateSSHKey(path, keyName, mail, passphrase string) error

	// TODO -> replace this with go-git merge
	PullMerge() (string, error)
}
