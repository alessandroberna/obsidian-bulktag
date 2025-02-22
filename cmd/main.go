package cmd

import (
	"fmt"
	"obsidian-bulktag/internal"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "obsidian-bulktag [folder]",
	Short: "Obsidian bulktag helps you add tags to your Obsidian vault",
	Long: `Obsidian bulktag is a CLI tool to semi-automatically add tags to your Obsidian notes based on folder structure.
It provides a navigable interface in the terminal to select folders and apply tags.`,
	Args: cobra.ExactArgs(1), // Expect exactly one argument: the folder path
	Run: func(cmd *cobra.Command, args []string) {
		internal.Config.Root = args[0]
		if _, err := os.Stat(internal.Config.Root); os.IsNotExist(err) {
			fmt.Println("Error: Folder does not exist:", internal.Config.Root)
			os.Exit(1)
		}
		if err := internal.Main(); err != nil {
			fmt.Println("UI Error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&internal.Config.ShowFiles, "no-files", true, "Do not show files in the file picker")
	rootCmd.Flags().Lookup("no-files").NoOptDefVal = "false"
}
