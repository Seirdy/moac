package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"moac-pwgen": main1,
	}))
}

func TestScripts(t *testing.T) {
	update := flag.Bool("u", false, "update testscript output files")

	testscript.Run(t, testscript.Params{
		Dir:           filepath.Join("testdata", "scripts"),
		UpdateScripts: *update,
		TestWork:      false,
	})
}
