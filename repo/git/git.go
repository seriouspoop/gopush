package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	gitCfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
)

const (
	gopushDir = ".gopush"
	keyName   = "gopush_key"
)

type Errors struct {
	RemoteNotFound      error
	RemoteNotLoaded     error
	RemoteAlreadyExists error
	RepoAlreadyExists   error
	RepoNotFound        error
	PullFailed          error
	AuthNotFound        error
	InvalidAuthMethod   error
	InvalidPassphrase   error
	KeyNotSupported     error
}

type Git struct {
	rootDir string
	repo    *git.Repository
	remote  *git.Remote
	err     *Errors
}

func New(gitErrors *Errors) (*Git, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Git{
		rootDir: rootDir,
		err:     gitErrors,
	}, nil
}

func (g *Git) GetRepo() error {
	repo, err := git.PlainOpen(g.rootDir)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return g.err.RepoNotFound
	}
	g.repo = repo
	return err
}

func (g *Git) CreateBranch(name model.Branch) error {
	return g.repo.CreateBranch(&gitCfg.Branch{
		Name:        name.String(),
		Remote:      "origin",
		Merge:       plumbing.NewBranchReferenceName(name.String()),
		Description: fmt.Sprintf("Desc: %s", name),
	})
}

func (g *Git) CreateRepo() error {
	repo, err := git.PlainInitWithOptions(g.rootDir, &git.PlainInitOptions{
		Bare: false,
		InitOptions: git.InitOptions{
			DefaultBranch: plumbing.Main,
		},
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return g.err.RepoAlreadyExists
	}
	g.repo = repo
	return err
}

func (g *Git) LoadRemote(remoteName string) error {
	remotes, _ := g.repo.Remotes()
	if len(remotes) == 0 {
		return g.err.RemoteNotFound
	}
	for _, remote := range remotes {
		if remote.Config().Name == remoteName {
			g.remote = remote
			break
		}
	}
	return nil
}

func (g *Git) GetRemoteDetails() (*model.Remote, error) {
	if g.remote == nil {
		return nil, g.err.RemoteNotLoaded
	}
	remoteDetails := &model.Remote{
		Name: g.remote.Config().Name,
		Url:  g.remote.Config().URLs[0],
	}
	return remoteDetails, nil
}

func (g *Git) AddRemote(remote *model.Remote) error {
	r, err := g.repo.CreateRemote(&gitCfg.RemoteConfig{
		Name:   remote.Name,
		URLs:   []string{remote.Url},
		Mirror: false,
	})
	if errors.Is(err, git.ErrRemoteExists) {
		return g.err.RemoteAlreadyExists
	}
	g.remote = r
	return err
}

func (g *Git) Pull(remote *model.Remote, branch model.Branch, auth *config.Credentials) error {
	if auth == nil {
		return g.err.AuthNotFound
	}
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	var Auth transport.AuthMethod
	if remote.AuthMode() == model.AuthHTTP {
		Auth = &http.BasicAuth{
			Username: auth.Username,
			Password: auth.Token,
		}
	} else if remote.AuthMode() == model.AuthSSH {
		sshPath := filepath.Join(os.Getenv("HOME"), gopushDir, keyName)
		sshKey, _ := os.ReadFile(sshPath)
		publicKey, err := ssh.NewPublicKeys("git", sshKey, auth.Token)
		if err != nil {
			if strings.Contains(err.Error(), "decryption password incorrect") {
				return g.err.InvalidPassphrase
			}
			return err
		}
		Auth = publicKey
	} else {
		return g.err.InvalidAuthMethod
	}
	err = w.Pull(&git.PullOptions{
		RemoteName:    remote.Name,
		RemoteURL:     remote.Url,
		ReferenceName: plumbing.NewBranchReferenceName(branch.String()),
		SingleBranch:  true,
		Auth:          Auth,
		Force:         false,
	})

	if strings.Contains(err.Error(), "unable to authenticate") {
		return g.err.KeyNotSupported
	} else if errors.Is(err, git.ErrNonFastForwardUpdate) {
		return g.err.PullFailed
	} else if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

func (g *Git) GetBranchNames() ([]model.Branch, error) {
	iter, err := g.repo.Branches()
	if err != nil {
		return nil, err
	}
	branches := []model.Branch{}
	iter.ForEach(func(r *plumbing.Reference) error {
		branches = append(branches, model.Branch(r.Name()[11:]))
		return nil
	})
	return branches, nil
}

func (g *Git) CheckoutBranch(name model.Branch) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name.String()),
		Keep:   true,
	})
	return err
}

func (g *Git) ChangeOccured() (bool, error) {
	w, err := g.repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := w.Status()
	if err != nil {
		return false, err
	}
	return !status.IsClean() || status.IsUntracked(g.rootDir), nil
}

func (g *Git) AddThenCommit(commitMsg string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		All:               true,
		AllowEmptyCommits: false,
		Amend:             false,
	})
	return err
}

func (g *Git) Push(remote *model.Remote, branch model.Branch, auth *config.Credentials) error {
	if auth == nil {
		return g.err.AuthNotFound
	}
	if remote == nil {
		return g.err.RemoteNotLoaded
	}
	var Auth transport.AuthMethod
	if remote.AuthMode() == model.AuthHTTP {
		Auth = &http.BasicAuth{
			Username: auth.Username,
			Password: auth.Token,
		}
	} else if remote.AuthMode() == model.AuthSSH {
		sshPath := filepath.Join(os.Getenv("HOME"), gopushDir, keyName)
		sshKey, _ := os.ReadFile(sshPath)
		publicKey, err := ssh.NewPublicKeys("git", sshKey, auth.Token)
		if err != nil {
			return err
		}
		Auth = publicKey
	} else {
		return g.err.InvalidAuthMethod
	}
	err := g.remote.Push(&git.PushOptions{
		RemoteName: remote.Name,
		RemoteURL:  remote.Url,
		Prune:      false,
		RefSpecs: []gitCfg.RefSpec{
			// final refspecs
			gitCfg.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", branch.String(), branch.String())),
		},
		Force: true,
		Auth:  Auth,
	})
	fmt.Println("hello")
	if strings.Contains(err.Error(), "unable to authenticate") {
		return g.err.KeyNotSupported
	} else if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}
