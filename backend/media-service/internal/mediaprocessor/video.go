package mediaprocessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
)

const (
	videoVariantMime = "video/mp4"
	videoThumbMime   = "image/jpeg"
	videoThumbMaxDim = 320
)

type videoRendition struct {
	variant filevariant.VariantType
	height  int
}

var videoRenditions = []videoRendition{
	{filevariant.Variant360p, 360},
	{filevariant.Variant720p, 720},
	{filevariant.Variant1080p, 1080},
}

func (p *ffmpegProcessor) processVideo(ctx context.Context, dir, inputPath string, req Request) ([]ProcessedVariant, error) {
	srcW, srcH := dimsOrDefault(req.SourceWidth, req.SourceHeight, 720)

	variants := make([]ProcessedVariant, 0, len(videoRenditions)+2)

	thumb, err := p.renderVideoThumb(ctx, dir, inputPath, srcW, srcH, req.SourceDuration)
	if err != nil {
		return nil, err
	}
	variants = append(variants, thumb)

	compressedHeight := even(min(srcH, 720))
	compressed, err := p.renderVideoRendition(ctx, dir, inputPath, filevariant.VariantCompressed, compressedHeight, srcW, srcH, "30", req.SourceDuration)
	if err != nil {
		return nil, err
	}
	variants = append(variants, compressed)

	for _, r := range videoRenditions {
		if r.height > srcH {
			continue
		}
		rendition, err := p.renderVideoRendition(ctx, dir, inputPath, r.variant, even(r.height), srcW, srcH, "23", req.SourceDuration)
		if err != nil {
			return nil, err
		}
		variants = append(variants, rendition)
	}

	return variants, nil
}

func (p *ffmpegProcessor) renderVideoThumb(ctx context.Context, dir, inputPath string, srcW, srcH int, duration *int) (ProcessedVariant, error) {
	outPath := filepath.Join(dir, "thumb.jpg")
	seek := "1"
	if duration != nil && *duration <= 1 {
		seek = "0"
	}
	vf := fmt.Sprintf(
		"scale='min(%d,iw)':'min(%d,ih)':force_original_aspect_ratio=decrease",
		videoThumbMaxDim, videoThumbMaxDim,
	)
	if err := p.ff.run(ctx,
		"-ss", seek,
		"-i", inputPath,
		"-frames:v", "1",
		"-vf", vf,
		"-q:v", "4",
		"-y", outPath,
	); err != nil {
		return ProcessedVariant{}, fmt.Errorf("render video thumb: %w", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return ProcessedVariant{}, fmt.Errorf("read video thumb: %w", err)
	}

	outW, outH := scaleToFit(srcW, srcH, videoThumbMaxDim)
	return ProcessedVariant{
		Variant:  filevariant.VariantThumb,
		Data:     data,
		MimeType: videoThumbMime,
		Width:    intPtr(outW),
		Height:   intPtr(outH),
		Duration: duration,
	}, nil
}

func (p *ffmpegProcessor) renderVideoRendition(
	ctx context.Context,
	dir, inputPath string,
	variant filevariant.VariantType,
	height, srcW, srcH int,
	crf string,
	duration *int,
) (ProcessedVariant, error) {
	outPath := filepath.Join(dir, fmt.Sprintf("%s.mp4", variant.String()))
	vf := fmt.Sprintf("scale=-2:%d", height)
	if err := p.ff.run(ctx,
		"-i", inputPath,
		"-vf", vf,
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-crf", crf,
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-y", outPath,
	); err != nil {
		return ProcessedVariant{}, fmt.Errorf("render video rendition %s: %w", variant.String(), err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return ProcessedVariant{}, fmt.Errorf("read video rendition %s: %w", variant.String(), err)
	}

	outW := height
	if srcH > 0 {
		outW = even(srcW * height / srcH)
	}
	return ProcessedVariant{
		Variant:  variant,
		Data:     data,
		MimeType: videoVariantMime,
		Width:    intPtr(outW),
		Height:   intPtr(height),
		Duration: duration,
	}, nil
}
