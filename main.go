package main

import (
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("mirc-portable", "mIRC"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	if err := os.MkdirAll(app.DataPath, 0o755); err != nil {
		log.Fatal().Err(err).Msg("Cannot create data path")
	}
	app.Process = filepath.Join(app.AppPath, "mirc.exe")
	app.Args = []string{
		"-r" + app.DataPath,
		"-noreg",
	}

	// Update settings
	settingsPath := filepath.Join(app.DataPath, "mirc.ini")
	if _, err := os.Stat(settingsPath); err == nil {
		ini.PrettyFormat = false
		log.Info().Msg("Update settings...")
		cfg, err := ini.LoadSources(ini.LoadOptions{
			IgnoreInlineComment:         true,
			SkipUnrecognizableLines:     false,
			UnescapeValueDoubleQuotes:   true,
			UnescapeValueCommentSymbols: true,
			PreserveSurroundedQuote:     true,
			SpaceBeforeInlineComment:    true,
			UnparseableSections: []string{
				"chanfolder",
				"dirs",
				"colors",
				"extensions",
				"ssl",
				"about",
				"afiles",
				"options",
				"pfiles",
				"rfiles",
				"user",
				"agent",
				"windows",
				"text",
				"files",
				"mirc",
				"dde",
				"ports",
			},
		}, settingsPath)
		if err == nil {
			cfg.Section("update").Key("check").SetValue("0")
			if err := cfg.SaveTo(settingsPath); err != nil {
				log.Error().Err(err).Msg("Write settings")
			}
		} else {
			log.Error().Err(err).Msg("Load mirc.ini file")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
