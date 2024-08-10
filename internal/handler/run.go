package handler

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/seriouspoop/gopush/model"
	"github.com/seriouspoop/gopush/svc"
	"github.com/spf13/cobra"
)

const (
	newBranchFlag   = "new-branch"
	setUpstreamFlag = "set-upstream"
)

func Run(s servicer) *cobra.Command {
	var newBranch string
	setUpstreamBranch := false

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "runs tests and push on remote.",
		Long: heredoc.Doc(`

			run command generates the mocks in current directory with go generate ./...
			then runs tests with go test ./...
			If all tests are passed, then modified files are staged following 
			push on the current repo's remote counterpart.

			[NOTE] Before pushing changes, changes from the remote main are pulled and are 
			attempted to merge to current branch. 
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("no args required for run")
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SetErrPrefix("❌ Error:")
			// Load repo and remote
			err := s.LoadProject()
			if err != nil {
				return err
			}

			// Load config
			err = s.LoadConfig()
			if err != nil {
				return err
			}

			fmt.Println("Pulling remote changes...")
			err = s.Pull(false)
			if err != nil {
				return err
			}
			fmt.Println("✅ changes fetched.")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := model.Branch(newBranch)
			if branch.Valid() {
				branchExist := false
				branchExist, err := s.SwitchBranchIfExists(branch)
				if err != nil {
					return err
				}
				if !branchExist {
					err = s.CreateBranchAndSwitch(branch)
					if err != nil {
						return err
					}
				}
			}

			// Generate Tests and Run
			testPresent, err := s.CheckTestsAndRun()
			if err != nil {
				return err
			}
			if testPresent {
				fmt.Println("✅ All tests passed. Staging changes...")
			} else {
				fmt.Println("No tests found. Staging changes...")
			}

			// stage changes
			err = s.StageChanges()
			if err != nil {
				return err
			}

			//TODO -> pull and merge from main

			// Push changes
			fmt.Println("\nPushing changes...")
			output, err := s.Push(setUpstreamBranch)
			if err != nil {
				if output != "" {
					fmt.Println(output)
				}
				if errors.Is(err, svc.ErrAuthNotFound) {
					fmt.Println(heredoc.Doc(`
					Auth credentials for current remote are missing.
					Run "gopush init" first to setup auth credentials.
					`))
				}
				return err
			}
			if output != "" {
				fmt.Println(output)
			}
			return nil
		},
	}
	runCmd.PersistentFlags().StringVarP(&newBranch, newBranchFlag, "b", "", "create new branch and set-upstream")
	runCmd.PersistentFlags().BoolVarP(&setUpstreamBranch, setUpstreamFlag, "u", false, "upstreams the given branch to remote")
	runCmd.MarkFlagsMutuallyExclusive(newBranchFlag, setUpstreamFlag)
	return runCmd
}
