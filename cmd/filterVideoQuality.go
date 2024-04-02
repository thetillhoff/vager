/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package cmd

import (
	"github.com/spf13/cobra"
	vager "github.com/thetillhoff/vager/pkg/vager"
)

// filterVideoQuality represents the filterVideoQuality command
var filterVideoQuality = &cobra.Command{
	Use:   "filterVideoQuality",
	Short: "FilterVideoQuality checks for each subfolder whether there are multiple videos with different resolutions and removes all but the highest (max 1080p)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vager.FilterVideoQuality(args[0], dryRun, verbose)
	},
}

func init() {
	rootCmd.AddCommand(filterVideoQuality)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filterVideoQuality.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filterVideoQuality.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
