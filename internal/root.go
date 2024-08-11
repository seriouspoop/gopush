package internal

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/seriouspoop/gopush/internal/handler"
	"github.com/seriouspoop/gopush/repo/git"
	"github.com/seriouspoop/gopush/repo/script"
	"github.com/seriouspoop/gopush/svc"
	"github.com/spf13/cobra"
)

type Root struct {
	git  *git.Git
	bash *script.Bash
	s    *svc.Svc
}

func NewRoot() (*Root, error) {
	gitHelper, err := git.New(&git.Errors{
		RemoteNotFound:      svc.ErrRemoteNotFound,
		RemoteNotLoaded:     svc.ErrRemoteNotLoaded,
		RemoteAlreadyExists: svc.ErrRemoteAlreadyExists,
		RepoAlreadyExists:   svc.ErrRepoAlreadyExists,
		RepoNotFound:        svc.ErrRepoNotFound,
		PullFailed:          svc.ErrPullFailed,
		AuthNotFound:        svc.ErrAuthNotFound,
		InvalidAuthMethod:   svc.ErrInvalidAuthMethod,
	})
	if err != nil {
		return nil, err
	}

	bashHelper := script.New(&script.Error{
		FileNotExists: svc.ErrFileNotFound,
	})

	s := svc.New(gitHelper, bashHelper)

	return &Root{
		git:  gitHelper,
		bash: bashHelper,
		s:    s,
	}, nil
}

func (r *Root) RootCMD() *cobra.Command {
	rootCMD := &cobra.Command{
		Use:     "gopush",
		Version: "1.0.2",
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
