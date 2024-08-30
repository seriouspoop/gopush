package handler

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/seriouspoop/gopush/svc"
	"github.com/seriouspoop/gopush/utils"
	"github.com/spf13/cobra"
)

func Init(s servicer) *cobra.Command {
	// TODO -> verbose implementation
	// var verbose bool
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

			err = s.LoadProject()
			if err != nil {
				if errors.Is(err, svc.ErrRepoNotFound) {
					err := s.InitializeRepo()
					if err != nil {
						return err
					}
					utils.Logger(utils.LOG_SUCCESS, "repository initialized")
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
				utils.Logger(utils.LOG_SUCCESS, "remote initialized")
			}

			utils.Logger(utils.LOG_INFO, "Generating config file...")
			err = s.SetRemoteAuth()
			if err != nil {
				return err
			}
			utils.Logger(utils.LOG_SUCCESS, "authorization set")

			err = s.LoadConfig()
			if err != nil {
				return err
			}
			utils.Logger(utils.LOG_SUCCESS, "config file generated")

			// staging current files
			utils.Logger(utils.LOG_INFO, "Staging changes...")
			err = s.StageChanges()
			if err != nil {
				return err
			}

			utils.Logger(utils.LOG_INFO, "Pulling commits from main...")
			err = s.Pull(true)
			if err != nil {
				if errors.Is(err, svc.ErrPullFailed) {
					utils.Logger(utils.LOG_INFO, "Remote pull failed, try pulling manually.")
				}
				return err
			}
			utils.Logger(utils.LOG_SUCCESS, "pulled changes")
			fmt.Println(heredoc.Doc(`

				Now you will be able to use "gopush run" command for you workflow.
				See "gopush run --help" for more details.`))
			return nil
		},
	}

	// initCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output for each step")

	return initCmd
}
