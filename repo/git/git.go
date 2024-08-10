package git

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-git/go-git/v5"
	gitCfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
)

type Errors struct {
	RemoteNotFound      error
	RemoteNotLoaded     error
	RemoteAlreadyExists error
	RepoAlreadyExists   error
	RepoNotFound        error
	PullFailed          error
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

// func (g *Git) Fetch() error {
// 	if g.remote == nil {
// 		return g.err.RemoteNotLoaded
// 	}
// 	fetchOpts := &git.FetchOptions{
// 		RemoteName: g.remote.Config().Name,
// 		RemoteURL:  g.remote.Config().URLs[0],
// 		Force:      false,
// 		Prune:      false,
// 		Progress:   os.Stdout,
// 		RefSpecs: []gitCfg.RefSpec{
// 			gitCfg.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
// 		},
// 		InsecureSkipTLS: true,
// 	}
// 	err := g.remote.Fetch(fetchOpts)
// 	if err == git.NoErrAlreadyUpToDate {
// 		return nil
// 	}

// 	return err
// }

func (g *Git) Pull(remote *model.Remote, branch model.Branch, auth *config.Credentials) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{
		RemoteName:    remote.Name,
		RemoteURL:     remote.Url,
		ReferenceName: plumbing.NewBranchReferenceName(branch.String()),
		SingleBranch:  true,
		Auth: &http.BasicAuth{
			Username: auth.Username,
			Password: auth.Token,
		},
		Progress: os.Stdout,
		Force:    false,
	})
	if errors.Is(err, git.ErrNonFastForwardUpdate) {
		return g.err.PullFailed
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

func (g *Git) conventionalCommit() string {
	var selection int
	var commitPrefix, commitMsg string
	isValidSelection := false
	for !isValidSelection {
		fmt.Print(heredoc.Doc(`
			Select the type of commit:
			  1. fix
			  2. feature
			  3. chore
			  4. refactor
			  5. ci
			Select an option: `))
		fmt.Scanf("%d", &selection)
		if selection < 1 || selection > 5 {
			isValidSelection = false
			fmt.Println("invalid selection.")
			continue
		}
		isValidSelection = true
	}
	switch selection {
	case 1:
		commitPrefix = "fix"
	case 2:
		commitPrefix = "feat"
	case 3:
		commitPrefix = "chore"
	case 4:
		commitPrefix = "ref"
	case 5:
		commitPrefix = "ci"
	default:
		commitPrefix = ""
	}
	fmt.Print("Enter commit message: ")
	reader := bufio.NewReader(os.Stdin)
	commitMsg, _ = reader.ReadString('\n')
	commitMsg = strings.TrimSpace(commitMsg)
	return fmt.Sprintf("%s: %s", commitPrefix, commitMsg)
}

func (g *Git) AddThenCommit() error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = w.Commit(g.conventionalCommit(), &git.CommitOptions{
		All:               true,
		AllowEmptyCommits: false,
		Amend:             false,
	})
	return err
}

func (g *Git) Push(remote *model.Remote, branch model.Branch, auth *config.Credentials) error {
	err := g.remote.Push(&git.PushOptions{
		RemoteName: remote.Name,
		RemoteURL:  remote.Url,
		Prune:      false,
		RefSpecs: []gitCfg.RefSpec{
			gitCfg.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/remotes/origin/%s", branch.String(), branch.String())),
		},
		Force:    false,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: auth.Username,
			Password: auth.Token,
		},
	})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}
