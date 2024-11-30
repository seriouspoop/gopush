package internal

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/seriouspoop/gopush/gogresSvc"
	"github.com/seriouspoop/gopush/internal/handler"
	"github.com/seriouspoop/gopush/repo/git"
	"github.com/seriouspoop/gopush/repo/script"
	"github.com/spf13/cobra"
)

type Root struct {
	git  *git.Git
	bash *script.Bash
	s    *gogresSvc.Svc
}

func NewRoot() (*Root, error) {
	gitHelper, err := git.New(&git.Errors{
		RemoteNotFound:       gogresSvc.ErrRemoteNotFound,
		RemoteNotLoaded:      gogresSvc.ErrRemoteNotLoaded,
		RemoteAlreadyExists:  gogresSvc.ErrRemoteAlreadyExists,
		RepoAlreadyExists:    gogresSvc.ErrRepoAlreadyExists,
		RepoNotFound:         gogresSvc.ErrRepoNotFound,
		PullFailed:           gogresSvc.ErrPullFailed,
		AuthNotFound:         gogresSvc.ErrAuthNotFound,
		InvalidAuthMethod:    gogresSvc.ErrInvalidAuthMethod,
		InvalidPassphrase:    gogresSvc.ErrInvalidPassphrase,
		KeyNotSupported:      gogresSvc.ErrKeyNotSupported,
		AlreadyUpToDate:      gogresSvc.ErrAlreadyUpToDate,
		RemoteBranchNotFound: gogresSvc.ErrRemoteBranchNotFound,
	})
	if err != nil {
		return nil, err
	}

	bashHelper := script.New(&script.Error{
		FileNotExists: gogresSvc.ErrFileNotFound,
	})

	s := gogresSvc.New(gitHelper, bashHelper)

	return &Root{
		git:  gitHelper,
		bash: bashHelper,
		s:    s,
	}, nil
}

func (r *Root) RootCMD() *cobra.Command {
	rootCMD := &cobra.Command{
		Use:     "gopush",
		Version: "1.1.2",
		Short: heredoc.Doc(`

			 ██████╗   ██████╗  ██████╗  ██╗   ██╗ ███████╗ ██╗  ██╗
			██╔════╝  ██╔═══██╗ ██╔══██╗ ██║   ██║ ██╔════╝ ██║  ██║
			██║  ███╗ ██║   ██║ ██████╔╝ ██║   ██║ ███████╗ ███████║
			██║   ██║ ██║   ██║ ██╔═══╝  ██║   ██║ ╚════██║ ██╔══██║
			╚██████╔╝ ╚██████╔╝ ██║      ╚██████╔╝ ███████║ ██║  ██║
			 ╚═════╝   ╚═════╝  ╚═╝       ╚═════╝  ╚══════╝ ╚═╝  ╚═╝
			`),
		Long: "",
	}

	//TODO - remove git and bash dependency from handler
	rootCMD.AddCommand(handler.Run(r.s))
	rootCMD.AddCommand(handler.Init(r.s))

	return rootCMD
}
