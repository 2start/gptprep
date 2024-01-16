package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gprep",
		Short: "gprep prints filenames to be loaded into the clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle flags and execute the main logic
			extensions := viper.GetStringSlice("extension")
			excludes := viper.GetStringSlice("exclude")

			files, err := findFiles(extensions, excludes)
			if err != nil {
				fmt.Println("Error finding files:", err)
				return
			}
			
			// Print the filenames
			for _, file := range files {
				fmt.Println(file)
			}
			
			concatenatedContent, err := concatenateFiles(files)
			if err != nil {
				fmt.Println("Error concatenating file content:", err)
				return
			}
			err = LoadToClipboard(concatenatedContent)
			if err != nil {
				fmt.Println("Error loading content to clipboard:", err)
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
