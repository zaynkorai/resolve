package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var command string
var errorOut string
var prompt string
var exitError string

var rootCmd = &cobra.Command{
	Use:   "resolve",
	Short: "resolve terminal command error using AI",
	Run: func(cmd *cobra.Command, args []string) {
		printLastCommandOutput()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func printLastCommandOutput() {

	cmdStr := "bash"
	args := []string{"-c", command}

	cmd := exec.Command(cmdStr, args...)
	combinedOutput, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Command Output: %s\n", combinedOutput)
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitError = fmt.Sprintf("Command exited with exitError: %s\n", exitErr)
		}
	} else {
		fmt.Printf("Command executed successfully. No action needed: %s\n\n", combinedOutput)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&command, "cmd", "c", "", "Specify a command to run")
	rootCmd.Flags().StringVarP(&errorOut, "error", "e", "", "error you want to search about")
	rootCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "prompt for captured error")
}
