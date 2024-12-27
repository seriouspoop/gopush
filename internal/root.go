package internal

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/seriouspoop/gopush/gopushSvc"
	"github.com/seriouspoop/gopush/internal/handler"
	"github.com/seriouspoop/gopush/repo/git"
	"github.com/seriouspoop/gopush/repo/script"
	"github.com/spf13/cobra"
)

type Root struct {
	git  *git.Git
	bash *script.Bash
	s    *gopushSvc.Svc
}

func NewRoot() (*Root, error) {
	gitHelper, err := git.New(&git.Errors{
		RemoteNotFound:       gopushSvc.ErrRemoteNotFound,
		RemoteNotLoaded:      gopushSvc.ErrRemoteNotLoaded,
		RemoteAlreadyExists:  gopushSvc.ErrRemoteAlreadyExists,
		RepoAlreadyExists:    gopushSvc.ErrRepoAlreadyExists,
		RepoNotFound:         gopushSvc.ErrRepoNotFound,
		PullFailed:           gopushSvc.ErrPullFailed,
		AuthNotFound:         gopushSvc.ErrAuthNotFound,
		InvalidAuthMethod:    gopushSvc.ErrInvalidAuthMethod,
		InvalidPassphrase:    gopushSvc.ErrInvalidPassphrase,
		KeyNotSupported:      gopushSvc.ErrKeyNotSupported,
		AlreadyUpToDate:      gopushSvc.ErrAlreadyUpToDate,
		RemoteBranchNotFound: gopushSvc.ErrRemoteBranchNotFound,
	})
	if err != nil {
		return nil, err
	}

	bashHelper := script.New(&script.Error{
		FileNotExists: gopushSvc.ErrFileNotFound,
	})

	s := gopushSvc.New(gitHelper, bashHelper)

	return &Root{
		git:  gitHelper,
		bash: bashHelper,
		s:    s,
	}, nil
}

func (r *Root) RootCMD() *cobra.Command {
	rootCMD := &cobra.Command{
		Use:     "gopush",
		Version: "1.1.3",
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
