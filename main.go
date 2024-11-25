package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/besser/go-mod-replace/cmd"
	"github.com/besser/go-mod-replace/git"

	"golang.org/x/mod/modfile"
)

var Version string = "v1.2.2"

func main() {
	cmd.StartFlags()

	const (
		DEV_ENV = "dev"
		STG_ENV = "stg"
		HML_ENV = "hml"
		PRD_ENV = "prd"

		GOMODULEFILE = "go.mod"
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

	if !cmd.Remove && len(env) == 0 {
		log.Fatalln("Environment value is required")
	}

	if env == STG_ENV && len(cmd.Branch) == 0 {
		log.Fatalln("Branch value is required")
	}

	goModFileName := GOMODULEFILE
	if len(cmd.FileName) > 0 {
		goModFileName = cmd.FileName
	}

	data, err := os.ReadFile(goModFileName)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read file %s: %w", goModFileName, err))
	}

	modFile, err := modfile.Parse(goModFileName, data, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse file %s: %w", goModFileName, err))
	}

	// cleanup replaces
	for _, replace := range modFile.Replace {
		modFile.DropReplace(replace.Old.Path, replace.Old.Version)
	}

	if !cmd.Remove {
		for _, req := range modFile.Require {
			if strings.Contains(req.Mod.Path, cmd.Domain) {
				switch env {
				// DEVELOPMENT
				case DEV_ENV:
					debugPath := strings.ReplaceAll(req.Mod.Path, cmd.Domain, "./..")
					modFile.AddReplace(req.Mod.Path, "", debugPath, "")

				// STAGING or HOMOLOG
				case STG_ENV, HML_ENV:
					modFile.AddReplace(req.Mod.Path, "", req.Mod.Path, cmd.Branch)

				// PRODUCTION
				case PRD_ENV:
					repoCommon := fmt.Sprintf("https://x-token-auth:%s@%s", strings.ReplaceAll(cmd.Token, "\"", ""), req.Mod.Path)

					r, err := git.GetRepo(repoCommon)
					if err != nil {
						log.Fatal(fmt.Errorf("failed to open '%s' repo: %w", repoCommon, err))
					}

					latestTagName, err := git.GetLatestTagFromRepository(r)
					if err != nil {
						log.Fatal(fmt.Errorf("failed to get latest tag from 'commons' repo: %w", err))
					}

					if err := modFile.AddRequire(req.Mod.Path, latestTagName); err != nil {
						log.Fatal(fmt.Errorf("failed to add '%s' to go.mod: %w", req.Mod.Path, err))
					}
				}
			}
		}
	}

	modFile.Cleanup()

	newData, err := modFile.Format()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to format file: %w", err))
	}

	if e := os.WriteFile(goModFileName, newData, 0600); e != nil {
		log.Fatal(fmt.Errorf("failed to write file %s: %w", goModFileName, e))
	}

	if cmd.Show {
		log.Println(string(newData))
	}
}
