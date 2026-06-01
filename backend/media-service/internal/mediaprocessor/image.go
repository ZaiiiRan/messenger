package mediaprocessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
)

const imageVariantMime = "image/jpeg"

type avatarPreset struct {
	variant filevariant.VariantType
	size    int
}

var avatarPresets = []avatarPreset{
	{filevariant.VariantAvatar, 640},
	{filevariant.VariantAvatarMiddle, 320},
	{filevariant.VariantAvatarSmall, 64},
}

const (
	photoCompressedMaxDimension = 1600
	photoThumbMaxDimension      = 320
)

func (p *ffmpegProcessor) processImage(ctx context.Context, dir, inputPath string, req Request) ([]ProcessedVariant, error) {
	if req.ProcessingType.RequiresCrop() {
		return p.processAvatar(ctx, dir, inputPath, req)
	}
	return p.processPhoto(ctx, dir, inputPath, req)
}

func (p *ffmpegProcessor) processAvatar(ctx context.Context, dir, inputPath string, req Request) ([]ProcessedVariant, error) {
	crop := req.Crop
	if crop == nil {
		return nil, fmt.Errorf("avatar processing requires crop params")
	}

	variants := make([]ProcessedVariant, 0, len(avatarPresets))
	for _, preset := range avatarPresets {
		outPath := filepath.Join(dir, fmt.Sprintf("%s.jpg", preset.variant.String()))
		vf := fmt.Sprintf(
			"crop=%d:%d:%d:%d,scale=%d:%d:flags=lanczos",
			crop.Width, crop.Height, crop.X, crop.Y, preset.size, preset.size,
		)
		if err := p.ff.run(ctx,
			"-i", inputPath,
			"-vf", vf,
			"-frames:v", "1",
			"-q:v", "2",
			"-y", outPath,
		); err != nil {
			return nil, fmt.Errorf("render avatar %s: %w", preset.variant.String(), err)
		}

		data, err := os.ReadFile(outPath)
		if err != nil {
			return nil, fmt.Errorf("read avatar %s: %w", preset.variant.String(), err)
		}
		variants = append(variants, ProcessedVariant{
			Variant:  preset.variant,
			Data:     data,
			MimeType: imageVariantMime,
			Width:    intPtr(preset.size),
			Height:   intPtr(preset.size),
		})
	}
	return variants, nil
}

func (p *ffmpegProcessor) processPhoto(ctx context.Context, dir, inputPath string, req Request) ([]ProcessedVariant, error) {
	srcW, srcH := dimsOrDefault(req.SourceWidth, req.SourceHeight, photoCompressedMaxDimension)

	type photoPreset struct {
		variant filevariant.VariantType
		maxDim  int
		quality string
	}
	presets := []photoPreset{
		{filevariant.VariantCompressed, photoCompressedMaxDimension, "3"},
		{filevariant.VariantThumb, photoThumbMaxDimension, "6"},
	}

	variants := make([]ProcessedVariant, 0, len(presets))
	for _, preset := range presets {
		outW, outH := scaleToFit(srcW, srcH, preset.maxDim)
		outPath := filepath.Join(dir, fmt.Sprintf("%s.jpg", preset.variant.String()))
		vf := fmt.Sprintf(
			"scale='min(%d,iw)':'min(%d,ih)':force_original_aspect_ratio=decrease",
			preset.maxDim, preset.maxDim,
		)
		if err := p.ff.run(ctx,
			"-i", inputPath,
			"-vf", vf,
			"-frames:v", "1",
			"-q:v", preset.quality,
			"-y", outPath,
		); err != nil {
			return nil, fmt.Errorf("render photo %s: %w", preset.variant.String(), err)
		}

		data, err := os.ReadFile(outPath)
		if err != nil {
			return nil, fmt.Errorf("read photo %s: %w", preset.variant.String(), err)
		}
		variants = append(variants, ProcessedVariant{
			Variant:  preset.variant,
			Data:     data,
			MimeType: imageVariantMime,
			Width:    intPtr(outW),
			Height:   intPtr(outH),
		})
	}
	return variants, nil
}

func dimsOrDefault(w, h *int, def int) (int, int) {
	if w == nil || h == nil || *w <= 0 || *h <= 0 {
		return def, def
	}
	return *w, *h
}
