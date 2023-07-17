/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package cmd

import (
	"github.com/spf13/cobra"
	videomanager "github.com/thetillhoff/video-manager/pkg/video-manager"
)

// flattenCmd represents the flatten command
var flattenCmd = &cobra.Command{
	Use:   "flatten",
	Short: "Flattens file structures with similar videos so that only highest resolution remains (max 1080p)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		videomanager.Flatten(args[0], dryRun, verbose)
	},
}

func init() {
	rootCmd.AddCommand(flattenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flattenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flattenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
