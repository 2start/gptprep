package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/2start/gptprep/internal/clipboard"
	"github.com/2start/gptprep/internal/filesearch"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gprep",
		Short: "gprep prints filenames to be loaded into the clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle flags and execute the main logic
			extensions := viper.GetStringSlice("extension")
			excludes := viper.GetStringSlice("exclude")

			files, concatenatedContent, err := filesearch.SearchAndConcatenateFiles(extensions, excludes)
			if err != nil {
				fmt.Println("Error finding and combining files:", err)
				return
			}
			
			if err != nil {
				fmt.Println("Error concatenating file content:", err)
				return
			}

			err = clipboard.LoadToClipboard(concatenatedContent)
			if err != nil {
				fmt.Println("Error loading content to clipboard:", err)
			}

			fmt.Println("Files loaded into the clipboard:")
			for _, file := range files {
				fmt.Println(file)
			}
		},
	}

	rootCmd.PersistentFlags().StringSlice("extension", []string{}, "File extensions to include")
	rootCmd.PersistentFlags().StringSlice("exclude", []string{".gitignore"}, "Patterns to exclude")
	viper.BindPFlag("extension", rootCmd.PersistentFlags().Lookup("extension"))
	viper.BindPFlag("exclude", rootCmd.PersistentFlags().Lookup("exclude"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
