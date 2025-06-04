package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zaynkorai/resolve/ai"
)

var errorOut string

var rootCmd = &cobra.Command{
	Use:   "resolve",
	Short: "resolve terminal command error using AI",
	Run: func(cmd *cobra.Command, args []string) {
		command, _ := cmd.Flags().GetString("cmd")
		details, _ := cmd.Flags().GetString("details")
		printLastCommandOutput(context.Background(), os.Getenv("GOOGLE_API_KEY"), command, details)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func printLastCommandOutput(ctx context.Context, googleAPIKey, command, details string) {
	if command == "" {
		fmt.Println(color.RedString("No command provided. Use --cmd or -c to specify a command."))
		return
	}

	args := []string{"-c", command}

	cmd := exec.Command("bash", args...)
	combinedOutput, err := cmd.CombinedOutput()

	if err != nil {
		agent, err := ai.NewAgents(ctx, googleAPIKey)
		if err != nil {
			fmt.Println(color.RedString("Failed to initialize AI agent: %v", err))
			return
		}
		fmt.Println(color.YellowString("Asking about issue to LLM..."))
		answer, err := agent.ResolveGivenIssue(ctx, details, string(combinedOutput))
		if err != nil {
			fmt.Println(color.RedString("Error from LLM: %v", err))
			return
		}
		fmt.Println(color.GreenString("Here is the solution!!!"))

		render(answer)
	} else {
		fmt.Println(color.GreenString("Command executed successfully. No action needed"))
	}
}

func init() {
	rootCmd.Flags().StringP("cmd", "c", "", "Specify a command to run which failed with error")
	rootCmd.Flags().StringP("details", "d", "", "Any context you want to add the issue.")
	rootCmd.Flags().StringVarP(&errorOut, "error", "e", "", "Error you want to search about.")
}

func render(content string) {
	renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("auto"))
	if err != nil {
		fmt.Printf("Error rendering answer: %v\n", err)
		return
	}

	out, err := renderer.Render(content)
	if err != nil {
		fmt.Printf("Error rendering markdown: %v\n", err)
		return
	}
	fmt.Print(out)
}
