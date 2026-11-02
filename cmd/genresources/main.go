// Command genresources writes the provider's resource schemas from the
// config kind JSON Schemas.
//
// By default it reads the schemas embedded in the config-schemas module this
// repository requires, so bumping that module and running go generate is the
// whole of adding a kind:
//
//	go run ./cmd/genresources -out internal/provider
//
// Each kind becomes zz_generated_<kind>.go holding a kindSpec: the framework
// schema attributes and the field map between them and the document spec.
// A valueFrom secret becomes a write-only <name>_wo argument and a
// <name>_wo_version trigger.
//
// -schemas reads a directory instead, which the tests use to generate a test
// kind with every supported shape.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	configschemas "github.com/vulnetix/config-schemas"
)

const schemasModule = "github.com/vulnetix/config-schemas"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "genresources:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fl := flag.NewFlagSet("genresources", flag.ContinueOnError)
	out := fl.String("out", "internal/provider", "directory to write generated files to")
	dir := fl.String("schemas", "", "read schemas from this directory instead of the embedded config-schemas")
	suffix := fl.String("suffix", "", "file name suffix before .go, such as _test")
	varName := fl.String("var", "generatedKinds", "name of the generated kind list")
	version := fl.String("version", "", "config-schemas version for the header (default: from the build)")
	if err := fl.Parse(args); err != nil {
		return err
	}

	var (
		fsys    fs.FS
		root    string
		ver     = *version
		sources = "embedded"
	)
	if *dir != "" {
		fsys, root, sources = os.DirFS(*dir), ".", *dir
	} else {
		fsys, root = configschemas.FS(), "schemas"
	}
	if ver == "" {
		ver = moduleVersion()
		if *dir != "" {
			ver = "test schemas"
		}
	}

	set, err := loadSchemas(fsys, root)
	if err != nil {
		return err
	}
	g := &generator{set: set, version: ver}
	kinds := set.kinds()
	if len(kinds) == 0 {
		return fmt.Errorf("no kinds in %s schemas", sources)
	}

	if err := removeStale(*out, *suffix); err != nil {
		return err
	}
	var outs []*kindOut
	for _, k := range kinds {
		o, err := g.generate(k)
		if err != nil {
			return err
		}
		outs = append(outs, o)
		if err := os.WriteFile(filepath.Join(*out, "zz_generated_"+o.TypeName+*suffix+".go"), o.Source, 0o644); err != nil {
			return err
		}
	}
	reg, err := registry(ver, *varName, outs)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(*out, "zz_generated_kinds"+*suffix+".go"), reg, 0o644)
}

// removeStale deletes generated files a removed kind would otherwise leave
// behind. Files with another suffix belong to another run.
func removeStale(dir, suffix string) error {
	matches, err := filepath.Glob(filepath.Join(dir, "zz_generated_*.go"))
	if err != nil {
		return err
	}
	for _, m := range matches {
		base := strings.TrimSuffix(filepath.Base(m), ".go")
		if strings.HasSuffix(base, "_test") != (suffix == "_test") {
			continue
		}
		if err := os.Remove(m); err != nil {
			return err
		}
	}
	return nil
}

// moduleVersion is the config-schemas version compiled in.
func moduleVersion() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, d := range bi.Deps {
			if d.Path == schemasModule {
				if d.Replace != nil {
					return d.Replace.Version
				}
				return d.Version
			}
		}
	}
	return "(unknown version)"
}
