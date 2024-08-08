package handler

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
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

			//TODO - use svc package
			// fmt.Println("Fetching remote changes...")
			// err = s.FetchAndMerge()
			// if err != nil {
			// 	return err
			// }
			// fmt.Println("✅ changes fetched.")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// branchExist := false
			// branch := model.Branch(newBranch)
			// if branch.Valid() {
			// 	branches, err := git.GetBranchNames()
			// 	if err != nil {
			// 		return err
			// 	}
			// 	for _, br := range branches {
			// 		if br.String() == branch.String() {
			// 			fmt.Printf("Branch %s already exists. Switching branch...\n", branch.String())
			// 			branchExist = true
			// 			break
			// 		}
			// 	}
			// 	if !branchExist {
			// 		err := git.CreateBranch(branch)
			// 		if err != nil {
			// 			return err
			// 		}
			// 	}

			// 	err = git.CheckoutBranch(branch)
			// 	if err != nil {
			// 		return err
			// 	}

			// } else if branch.String() != "" {
			// 	fmt.Println("Branch name: ", branch.String())
			// 	return ErrBranchInvalid
			// }

			// present, err := bash.TestsPresent()
			// if err != nil {
			// 	return err
			// }
			// if present {
			// 	fmt.Println("\nGenerating and Running tests...")
			// 	output, _ := bash.GenerateMocks()
			// 	fmt.Print(output)
			// 	output, err := bash.RunTests()
			// 	fmt.Print(output)
			// 	if err != nil {
			// 		return ErrTestsFailed
			// 	}
			// 	fmt.Println("✅ All Tests passed. Staging changes.")
			// } else {
			// 	fmt.Println("No test files found, continuing...")
			// }
			// change, err := git.ChangeOccured()
			// if err != nil {
			// 	return err
			// }
			// if change {
			// 	err := git.AddThenCommit()
			// 	if err != nil {
			// 		return err
			// 	}
			// 	fmt.Println("✅ Files added.")
			// }
			// currBranch, err := bash.GetCurrentBranch()
			// if err != nil {
			// 	return err
			// }
			// if setUpstreamBranch {
			// 	output, err := bash.Push(currBranch, true)
			// 	if err != nil {
			// 		fmt.Println(output)
			// 		return err
			// 	}
			// } else {
			// 	output, err := bash.Push(currBranch, false)
			// 	if err != nil {
			// 		fmt.Println(output)
			// 		return err
			// 	}
			// }
			// fmt.Println("✅ Push Successful.")
			return nil
		},
	}
	runCmd.PersistentFlags().StringVarP(&newBranch, newBranchFlag, "b", "", "create new branch and set-upstream")
	runCmd.PersistentFlags().BoolVarP(&setUpstreamBranch, setUpstreamFlag, "u", false, "upstreams the given branch to remote")
	runCmd.MarkFlagsMutuallyExclusive(newBranchFlag, setUpstreamFlag)
	return runCmd
}
