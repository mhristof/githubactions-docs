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
	Version: version,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			panic(fmt.Sprintf("Error, file %s does not exist", args[0]))
		}

		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"args[0]": args[0],
			}).Error("Could not load file")
		}

		cfg, err := action.Load(data)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Could not decode action file")
		}

		fmt.Println(cfg.Markdown())
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
