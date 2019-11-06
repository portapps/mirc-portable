//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"

	"github.com/go-ini/ini"
	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/utl"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("mirc-portable", "mIRC"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "mirc.exe")
	app.Args = []string{
		"-r" + app.DataPath,
		"-noreg",
	}

	// Update settings
	settingsPath := utl.PathJoin(app.DataPath, "mirc.ini")
	if _, err := os.Stat(settingsPath); err == nil {
		ini.PrettyFormat = false
		Log.Info().Msg("Update settings...")
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
				Log.Error().Err(err).Msg("Write settings")
			}
		} else {
			Log.Error().Err(err).Msg("Load mirc.ini file")
		}
	}

	app.Launch(os.Args[1:])
}
