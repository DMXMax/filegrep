package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Two stage scan to find patterns outside of comments",
	Long: `Scan scans files in two stages, first a general grep to find candidate files.
The results are passed to a filter to strip comments, based on file extension`,
	RunE: scanFiles,
}

func scanFiles(cmd *cobra.Command, args []string) error {
	log.Println("Scanfiles Called")
	return nil
}
