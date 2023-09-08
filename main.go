package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/anfernee/gomod/pkg/module"
	"golang.org/x/mod/modfile"
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

		m := module.New(require.Mod.Path)
		latest, err := m.Latest()
		if err != nil {
			log.Printf("failed to load %s: %v", require.Mod.Path, err)
			continue
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
