package mediaprocessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
)

type audioEncoding struct {
	args []string
	ext  string
}

func audioEncodingForMime(mime string) audioEncoding {
	switch mime {
	case "audio/mpeg":
		return audioEncoding{args: []string{"-c:a", "libmp3lame", "-b:a", "128k"}, ext: "mp3"}
	case "audio/aac":
		return audioEncoding{args: []string{"-c:a", "aac", "-b:a", "128k"}, ext: "aac"}
	case "audio/ogg":
		return audioEncoding{args: []string{"-c:a", "libvorbis", "-q:a", "4"}, ext: "ogg"}
	case "audio/webm":
		return audioEncoding{args: []string{"-c:a", "libopus", "-b:a", "96k"}, ext: "webm"}
	case "audio/flac":
		return audioEncoding{args: []string{"-c:a", "flac", "-compression_level", "8"}, ext: "flac"}
	case "audio/wav":
		return audioEncoding{args: []string{"-c:a", "libmp3lame", "-b:a", "128k"}, ext: "mp3"}
	default:
		return audioEncoding{args: []string{"-c:a", "libmp3lame", "-b:a", "128k"}, ext: "mp3"}
	}
}

func (p *ffmpegProcessor) processAudio(ctx context.Context, dir, inputPath string, req Request) ([]ProcessedVariant, error) {
	enc := audioEncodingForMime(req.SourceMime)
	outPath := filepath.Join(dir, "compressed."+enc.ext)

	args := []string{"-i", inputPath, "-vn"}
	args = append(args, enc.args...)
	args = append(args, "-y", outPath)

	if err := p.ff.run(ctx, args...); err != nil {
		return nil, fmt.Errorf("compress audio: %w", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read compressed audio: %w", err)
	}

	return []ProcessedVariant{
		{
			Variant:  filevariant.VariantOriginal,
			Data:     data,
			MimeType: req.SourceMime,
			Duration: req.SourceDuration,
		},
	}, nil
}
