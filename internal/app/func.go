package app

import (
	"runtime/debug"
	"time"

	"github.com/urfave/cli/v2"
)

// WithValues set values.
func WithValues(
	name, description, version, copyright, authorName, authorEmail, buildTime string,
) Option {
	return func(app *cli.App) {
		app.Name = name
		app.Description = description
		app.Version = resolveVersion(version)
		app.Copyright = copyright

		app.Authors = []*cli.Author{
			{
				Name:  authorName,
				Email: authorEmail,
			},
		}

		app.HideVersion = false

		resolvedBuildTime := resolveBuildTime(buildTime)
		if resolvedBuildTime == "" {
			return
		}

		parsedBuildTime, err := time.Parse(time.RFC3339, resolvedBuildTime)
		if err != nil {
			panic(err)
		}

		app.Compiled = parsedBuildTime
	}
}

func resolveVersion(fallback string) string {
	bi, ok := debug.ReadBuildInfo()
	if !ok || fallback == "latest" {
		return fallback
	}

	if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}

	return fallback
}

func resolveBuildTime(fallback string) string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fallback
	}

	for _, s := range bi.Settings {
		if s.Key == "vcs.time" {
			return s.Value
		}
	}

	return fallback
}
