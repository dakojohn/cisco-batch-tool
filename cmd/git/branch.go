package git

import (
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ryclarke/cisco-batch-tool/call"
	"github.com/ryclarke/cisco-batch-tool/config"
	"github.com/ryclarke/cisco-batch-tool/utils"
)

func addBranchCmd() *cobra.Command {
	// branchCmd represents the branch command
	branchCmd := &cobra.Command{
		Use:     "branch <repository> ...",
		Aliases: []string{"checkout"},
		Short:   "Checkout a new branch across repositories",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return utils.ValidateRequiredConfig(config.Branch)
		},
		Run: func(cmd *cobra.Command, args []string) {
			call.Do(args, call.Wrap(gitUpdate, gitCheckout))
		},
	}

	branchCmd.Flags().StringP("branch", "b", "", "branch name (required)")
	viper.BindPFlag(config.Branch, branchCmd.Flags().Lookup("branch"))

	return branchCmd
}

func gitCheckout(name string, ch chan<- string) error {
	branch := viper.GetString(config.Branch)

	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = utils.RepoPath(name)

	output, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("git", "checkout", "-b", branch)
		cmd.Dir = utils.RepoPath(name)

		output, err = cmd.Output()
		if err != nil {
			return err
		}

		ch <- string(output)

		cmd = exec.Command("git", "push", "-u", "origin", branch)
		cmd.Dir = utils.RepoPath(name)

		output, err = cmd.Output()
		if err != nil {
			return err
		}
	}

	ch <- string(output)

	return nil
}
