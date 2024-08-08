package handler

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/seriouspoop/gopush/svc"
	"github.com/spf13/cobra"
)

func Init(s servicer) *cobra.Command {
	var verbose bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "initializes git repo with all the config setting",
		Long: heredoc.Doc(`
			Initialize a git repo and creates a toml config for configuration settings,
			remote name, repo, branch prefix etc. are stored in the gopush_config.toml
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("no args required for run")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SetErrPrefix("❌ Error:")

			// Generate gopush_config.toml
			err := s.SetUserPreference()
			if err != nil {
				return err
			}
			//TODO - migrate to svc methods
			err = s.LoadProject()
			if err != nil {
				if errors.Is(err, svc.ErrRepoNotFound) {
					err := s.InitializeRepo()
					if err != nil {
						return err
					}
				}

				// load config for add remote
				err = s.LoadConfig()
				if err != nil {
					return err
				}
				// Add remote
				err = s.InitializeRemote()
				if err != nil {
					return err
				}
			}
			fmt.Println("✅ Repository and Remote initialized.")
			err = s.SetRemoteAuth()
			if err != nil {
				return err
			}
			fmt.Println("✅ config file generated.")

			err = s.LoadConfig()
			if err != nil {
				return err
			}

			// staging current files
			err = s.StageChanges()
			if err != nil {
				return err
			}
			fmt.Println("✅ changes staged")
			// //TODO - implement pull origin main & merge remote/origin/main
			fmt.Println("Pulling commits from main...")
			err = s.Pull(true)
			if err != nil {
				if errors.Is(err, svc.ErrPullFailed) {
					fmt.Println("Remote pull failed, try pulling manually.")
				}
				return err
			}
			// output, err := bash.PullBranch(currBranch)
			// if err != nil {
			// 	fmt.Println(output)
			// 	return ErrMergeOverwritingConflict
			// }
			// if verbose {
			// 	fmt.Println(output)
			// }
			fmt.Println(heredoc.Doc(`
				✅ Pulled remote changes

				Now you will be able to use "gopush run" command for you workflow.
				See "gopush run --help" for more details.
			`))
			return nil
		},
	}

	initCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output for each step")

	return initCmd
}
