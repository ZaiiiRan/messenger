package settings

import "github.com/spf13/viper"

type ProcessorSettings struct {
	FfmpegBin  string `mapstructure:"ffmpeg_bin"`
	FfprobeBin string `mapstructure:"ffprobe_bin"`
	TmpDir     string `mapstructure:"tmp_dir"`
}

func SetProcessorDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".ffmpeg_bin", "ffmpeg")
	v.SetDefault(prefix+".ffprobe_bin", "ffprobe")
	v.SetDefault(prefix+".tmp_dir", "")
}
