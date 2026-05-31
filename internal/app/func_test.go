package app

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fallback string
		bi       *debug.BuildInfo
		want     string
	}{
		{
			name:     "fallback takes precedence",
			fallback: "v1.0.0",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: "v2.0.0"}},
			want:     "v1.0.0",
		},
		{
			name:     "nil build info returns latest",
			fallback: "",
			bi:       nil,
			want:     latest,
		},
		{
			name:     "devel returns latest",
			fallback: "",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}},
			want:     latest,
		},
		{
			name:     "empty version returns latest",
			fallback: "",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: ""}},
			want:     latest,
		},
		{
			name:     "module version used when no fallback",
			fallback: "",
			bi:       &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}},
			want:     "v1.2.3",
		},
		{
			name:     "pseudo version used when no fallback",
			fallback: "",
			bi: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.0.4-0.20260401175942-30fd43debbd8"},
			},
			want: "v0.0.4-0.20260401175942-30fd43debbd8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, resolveVersion(tt.fallback, tt.bi))
		})
	}
}

const (
	fallbackDateTime = "2026-04-01T18:00:00Z"
	wantDateTime     = "2026-04-01T12:00:00Z"
)

func TestResolveBuildTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fallback string
		bi       *debug.BuildInfo
		want     string
	}{
		{
			name:     "nil build info returns fallback",
			fallback: fallbackDateTime,
			bi:       nil,
			want:     fallbackDateTime,
		},
		{
			name:     "vcs.time returned when present",
			fallback: "",
			bi: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.time", Value: wantDateTime},
				},
			},
			want: wantDateTime,
		},
		{
			name:     "fallback returned when vcs.time absent",
			fallback: fallbackDateTime,
			bi: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abc123"},
				},
			},
			want: fallbackDateTime,
		},
		{
			name:     "empty fallback and no vcs.time returns empty",
			fallback: "",
			bi:       &debug.BuildInfo{Settings: []debug.BuildSetting{}},
			want:     "",
		},
		{
			name:     "vcs.time takes precedence over fallback",
			fallback: fallbackDateTime,
			bi: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.time", Value: wantDateTime},
				},
			},
			want: wantDateTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, resolveBuildTime(tt.fallback, tt.bi))
		})
	}
}
