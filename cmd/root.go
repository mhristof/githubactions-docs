package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mhristof/githubactions-docs/action"
	"github.com/mhristof/githubactions-docs/log"
	"github.com/spf13/cobra"
)

var version = "devel"

var rootCmd = &cobra.Command{
	Use:     "githubactions-docs",
	Short:   "Generate documentation for Github Actions similar to terraform-docs",
	Args:    cobra.ArbitraryArgs,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)
		var files []string

		files = append(files, args...)

		if len(args) == 0 {
			files = append(files, "action.yml")
		}

		for _, file := range files {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				log.WithFields(log.Fields{
					"file": file,
				}).Error("File not found")
				continue
			}

			data, err := ioutil.ReadFile(file)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"file": file,
				}).Error("Could not load file")
				continue
			}

			cfg, err := action.Load(data)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Could not decode action file")
				continue
			}

			fmt.Println(cfg.Markdown())
		}
	},
}

// Verbose Increase verbosity
func Verbose(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Panic(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().BoolP("dryrun", "n", false, "Dry run")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
