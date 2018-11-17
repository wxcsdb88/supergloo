package util

// PACKAGE WILL BE MOVED TO GO-CHECKPOINT

import (
	"context"
	"path/filepath"
	"time"

	checkpoint "github.com/solo-io/go-checkpoint"
)

// Telemetry sends telemetry information about solo.io products to Checkpoint server
func Telemetry(version string, t time.Time) {
	sigfile := filepath.Join(HomeDir(), ".soloio.sig")
	configDir, err := ConfigDir()
	if err == nil {
		sigfile = filepath.Join(configDir, "soloio.sig")
	}
	ctx := context.Background()
	report := &checkpoint.ReportParams{
		Product:       "test",
		Version:       version,
		StartTime:     t,
		EndTime:       time.Now(),
		SignatureFile: sigfile,
	}
	checkpoint.Report(ctx, report)
}
