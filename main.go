package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anfernee/gomod/pkg/module"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

const (
	// Ignore this module
	commentIgnore = "gomod:ignore"

	// Only allow patch version updates
	commentPatchOnly = "gomod:patch-only"
)

var (
	dir = flag.String("dir", ".", "directory to golang project")
)

func main() {
	flag.Parse()

	gomod := filepath.Join(*dir, "go.mod")

	data, err := os.ReadFile(gomod)
	if err != nil {
		log.Fatal(err)
	}

	file, err := modfile.Parse(gomod, data, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, require := range file.Require {
		log.Printf("require: %s %s", require.Mod.Path, require.Mod.Version)

		// Skip indirect dependencies
		if require.Indirect {
			continue
		}

		// Skip stdlib dependencies
		if require.Mod.Path == "std" {
			continue
		}

		if ignoreModule(require) {
			continue
		}

		m := module.New(require.Mod.Path)

		var (
			latest string
			err    error
		)

		if patchOnlyModule(require) {
			versions, err := m.Versions()
			if err != nil {
				log.Printf("failed to get versions for %s: %v", require.Mod.Path, err)
				continue
			}

			latest = latestPatch(require.Mod.Version, versions)

		} else {
			latest, err = m.Latest()
			if err != nil {
				log.Printf("failed to get latest version for %s: %v", require.Mod.Path, err)
				continue
			}
		}

		if require.Mod.Version != latest {
			require.Mod.Version = latest
		}
	}
	file.SetRequire(file.Require)

	if data, err = file.Format(); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(gomod, data, 0644); err != nil {
		log.Fatal(err)
	}
}

func ignoreModule(require *modfile.Require) bool {
	if require.Syntax.Before != nil {
		for _, comment := range require.Syntax.Before {
			if strings.Trim(comment.Token, "/ ") == commentIgnore {
				return true
			}
		}
	}
	return false
}

func patchOnlyModule(require *modfile.Require) bool {
	if require.Syntax.Before != nil {
		for _, comment := range require.Syntax.Before {
			if strings.Trim(comment.Token, "/ ") == commentPatchOnly {
				return true
			}
		}
	}
	return false
}

func latestPatch(current string, versions []string) string {
	cur := semver.MajorMinor(current)
	latest := current
	for _, version := range versions {
		if semver.MajorMinor(version) == cur && semver.Compare(version, latest) > 0 {
			latest = version
		}
	}
	return latest
}
