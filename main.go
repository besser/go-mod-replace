package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/besser/go-mod-replace/cmd"
	"golang.org/x/mod/modfile"
)

var Version string = "v1.0.0"

func main() {
	cmd.StartFlags()

	const (
		cmpEnv = "compose"
		devEnv = "development"
		stgEnv = "staging"
		prdEnv = "production"
	)

	if cmd.BuildVersion {
		fmt.Print(Version)
		os.Exit(0)
	}

	if cmd.Version {
		fmt.Printf("Version: %s (%s)\n", Version, runtime.Version())
		os.Exit(0)
	}

	env := strings.TrimSpace(strings.ToLower(cmd.Environment))

	if !cmd.Remove && len(cmd.Domain) == 0 {
		log.Fatalln("go.mod path not found")
	}

	if len(cmd.FilePath) == 0 {
		log.Fatalln("Domain value not found")
	}

	if (env == cmpEnv || env == stgEnv) && len(cmd.Branch) == 0 {
		log.Fatalln("Branch value is required")
	}

	fileName := fmt.Sprintf("%s/go.mod", cmd.FilePath)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read file %s: %w", fileName, err))
	}

	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse file %s: %w", fileName, err))
	}

	// cleanup replaces
	for _, replace := range file.Replace {
		file.DropReplace(replace.Old.Path, replace.Old.Version)
	}

	if env != prdEnv && !cmd.Remove {
		for _, req := range file.Require {
			if strings.Contains(req.Mod.Path, cmd.Domain) {
				if env == devEnv || cmd.Debug {
					debugPath := strings.ReplaceAll(req.Mod.Path, cmd.Domain, "./..")
					file.AddReplace(req.Mod.Path, "", debugPath, "")
				} else if env == cmpEnv || env == stgEnv || len(cmd.Branch) > 0 {
					file.AddReplace(req.Mod.Path, "", req.Mod.Path, cmd.Branch)
				}
			}
		}
	}

	file.Cleanup()

	newData, err := file.Format()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to format file: %w", err))
	}

	if e := os.WriteFile(fileName, newData, 0600); e != nil {
		log.Fatal(fmt.Errorf("failed to write file %s: %w", fileName, e))
	}
}
