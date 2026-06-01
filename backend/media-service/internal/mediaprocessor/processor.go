package mediaprocessor

import (
	"context"
	"fmt"
	"os"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/processing"
)

type Processor interface {
	Process(ctx context.Context, req Request) ([]ProcessedVariant, error)
}

type ffmpegProcessor struct {
	ff      *ffmpeg
	workDir string
}

func New(ffmpegBin, ffprobeBin, tmpDir string) Processor {
	return &ffmpegProcessor{
		ff:      newFFmpeg(ffmpegBin, ffprobeBin),
		workDir: tmpDir,
	}
}

func (p *ffmpegProcessor) Process(ctx context.Context, req Request) ([]ProcessedVariant, error) {
	dir, err := os.MkdirTemp(p.workDir, "media-*")
	if err != nil {
		return nil, fmt.Errorf("create work dir: %w", err)
	}
	defer os.RemoveAll(dir)

	inputPath, err := writeTemp(dir, "input"+extForMime(req.SourceMime), req.SourceData)
	if err != nil {
		return nil, err
	}

	switch req.ProcessingType.Category() {
	case processing.CategoryImage:
		return p.processImage(ctx, dir, inputPath, req)
	case processing.CategoryVideo:
		return p.processVideo(ctx, dir, inputPath, req)
	case processing.CategoryAudio:
		return p.processAudio(ctx, dir, inputPath, req)
	default:
		return nil, fmt.Errorf("unsupported processing category for type %q", req.ProcessingType.String())
	}
}

func extForMime(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "video/avi":
		return ".avi"
	case "video/quicktime":
		return ".mov"
	case "video/x-matroska":
		return ".mkv"
	case "audio/mpeg":
		return ".mp3"
	case "audio/ogg":
		return ".ogg"
	case "audio/wav":
		return ".wav"
	case "audio/webm":
		return ".weba"
	case "audio/aac":
		return ".aac"
	case "audio/flac":
		return ".flac"
	default:
		return ".bin"
	}
}
