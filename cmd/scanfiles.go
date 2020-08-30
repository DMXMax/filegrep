package cmd

import (
	"github.com/spf13/cobra"
	"log"
	_"os"
	"os/exec"
	_"io"
	"bytes"
)

var (
	scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Two stage scan to find patterns outside of comments",
	Long: `Scan scans files in two stages, first a general grep to find candidate files.
The results are passed to a filter to strip comments, based on file extension`,
	RunE: scanFiles,
}

	pattern string
	sourceDir string
)

func scanFiles(cmd *cobra.Command, args []string) error {
	log.Println("Scanfiles Called:")
	log.Printf("\tPattern:%s\n\t\t\tfile-source:%s", pattern, sourceDir)
	//or, ow := io.Pipe()

	//This grep command only returns the files names whose contents match the pattern
	grepCmd := exec.Command("grep", "-lr", pattern, sourceDir)
	if output, err := grepCmd.Output(); err == nil {

		fileList := bytes.Split(output, []byte("\n"))
		/*go processFileList(fileList, pattern, ow)

		if _, err := io.Copy(os.Stdout, or); err != nil {
			log.Fatal(err)
		}*/
		log.Printf("%#v\n",fileList)
		log.Println("Done")
	} else {
		if grepCmd.ProcessState.ExitCode() == 1 {
			log.Println("Scan: No files Found")
			return nil
		} else {
			return err
		}
	}

	return nil
}

func init(){
	scanCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "regex pattern to use for scan")
	scanCmd.Flags().StringVar(&sourceDir, "file-source", ".", "directory to grep")
	scanCmd.MarkFlagRequired("pattern")
}
