package cmd

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

// Declaring flags
var (
	Branch       string
	BuildVersion bool
	Debug        bool
	Domain       string
	Environment  string
	FileName     string
	Remove       bool
	Version      bool
)

func init() {
	flag.StringVarP(&Branch, "branch", "b", "", "project branch")
	flag.BoolVarP(&Debug, "debug", "d", false, "replace to use with local debug")
	flag.StringVarP(&Domain, "domain", "D", "", "project domain")
	flag.StringVarP(&Environment, "env", "e", "", "environment deploy")
	flag.StringVarP(&FileName, "file", "f", "", "go.mod file path")
	flag.BoolVarP(&Remove, "remove", "r", false, "remove all replaces")
	flag.BoolVarP(&Version, "version", "v", false, "show version")

	// hidden params
	flag.BoolVarP(&BuildVersion, "VERSION", "V", false, "show build version")
	flag.CommandLine.MarkHidden("VERSION")
}

// Arg returns the i'th command-line argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func Arg(i int) string {
	return flag.Arg(i)
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	return flag.Args()
}

// StartFlags initialize flags arguments to the app.
func StartFlags() {
	flag.Usage = showUsageFlags
	flag.Parse()
}

func showUsageFlags() {
	fmt.Fprintf(os.Stdout, "go-mod-replace\n\n")
	fmt.Fprintf(os.Stdout, "Usage: %s [optional flags]\n\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "Optional Flags:\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stdout, "\n")
}
