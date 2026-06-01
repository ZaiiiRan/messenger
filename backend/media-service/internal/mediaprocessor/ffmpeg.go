package mediaprocessor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ffmpeg struct {
	ffmpegBin  string
	ffprobeBin string
}

func newFFmpeg(ffmpegBin, ffprobeBin string) *ffmpeg {
	if ffmpegBin == "" {
		ffmpegBin = "ffmpeg"
	}
	if ffprobeBin == "" {
		ffprobeBin = "ffprobe"
	}
	return &ffmpeg{ffmpegBin: ffmpegBin, ffprobeBin: ffprobeBin}
}

func (f *ffmpeg) run(ctx context.Context, args ...string) error {
	baseArgs := []string{"-hide_banner", "-loglevel", "error", "-nostdin"}
	cmd := exec.CommandContext(ctx, f.ffmpegBin, append(baseArgs, args...)...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg: %w: %s", err, stderr.String())
	}
	return nil
}

func even(n int) int {
	if n < 2 {
		return 2
	}
	if n%2 != 0 {
		n--
	}
	return n
}

func scaleToFit(w, h, max int) (int, int) {
	if w <= 0 || h <= 0 {
		return max, max
	}
	if w <= max && h <= max {
		return w, h
	}
	if w >= h {
		return max, max * h / w
	}
	return max * w / h, max
}

func writeTemp(dir, name string, data []byte) (string, error) {
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}
	return path, nil
}
