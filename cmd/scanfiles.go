package cmd

import (
	"github.com/spf13/cobra"
	"text/scanner"
	"log"
	"os"
	"os/exec"
	"io"
	"io/ioutil"
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
/* *************************
	Next Up:
		JSON output
		Reading from a file list of configurations: Rule, Pattern, location
		Dynamic scanners to cull out different types of comments.
***************************/
func scanFiles(cmd *cobra.Command, args []string) error {
	log.Println("Scanfiles Called:")
	log.Printf("\tPattern:%s\n\t\t\tfile-source:%s", pattern, sourceDir)
	or, ow := io.Pipe()

	//This grep command only returns the files names whose contents match the pattern
	grepCmd := exec.Command("grep", "-lr", pattern, sourceDir)
	if output, err := grepCmd.Output(); err == nil {

		fileList := bytes.Split(output, []byte("\n"))
		go processFileList(fileList, pattern, ow)

		if _, err := io.Copy(os.Stdout, or); err != nil {
			log.Fatal(err)
		}
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

func processFileList(files [][]byte, pattern string, output *io.PipeWriter) {
	for i, f := range files {
		if len(f) > 0 {
			log.Printf("Scanning file %d: %s\n", i, string(f))
			r, w := io.Pipe()

			go scanFile(string(f), w)
			cmd := exec.Command("grep", "-n", pattern)
			cmd.Stdin = r
			if res, err := cmd.Output(); err != nil {
				if cmd.ProcessState.ExitCode() != 1 {
					log.Fatal(err)
				}
			} else {
				//here is where we edit res
				splitres := bytes.Split(res,[]byte("\n"))
				for _,ln := range splitres{
					if len(ln) > 0{
						it := []byte{}
						it = append(it, f...)
						it= append(it, []byte(":")...)
						it = append(it, ln...)
						it = append(it, []byte("\n")...)
						output.Write(it)
					}
				}
			}
		} else {
			log.Print("Empty row in files")
		}
	}
	log.Println("Closing...")
	output.Close()
}

func scanFile(file string, w *io.PipeWriter) {
	dat, _ := ioutil.ReadFile(file)
	var s scanner.Scanner
	s.Filename = file
	s.Init(bytes.NewReader(dat))
	s.Mode ^= scanner.ScanIdents
	s.Whitespace = 0
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		_, _ = w.Write([]byte(string(tok)))

	}
	w.Close()
}

func init(){
	scanCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "regex pattern to use for scan")
	scanCmd.Flags().StringVar(&sourceDir, "file-source", ".", "directory to grep")
	scanCmd.MarkFlagRequired("pattern")
}
