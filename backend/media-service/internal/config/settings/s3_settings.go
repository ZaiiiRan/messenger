package settings

import "github.com/spf13/viper"

type S3Settings struct {
	Endpoint          string `mapstructure:"endpoint"`
	AccessKeyID       string `mapstructure:"access_key_id"`
	SecretAccessKey   string `mapstructure:"secret_access_key"`
	SecretAccessToken string `mapstructure:"secret_access_token"`
	Bucket            string `mapstructure:"bucket"`
	Region            string `mapstructure:"region"`
	UseSSL            bool   `mapstructure:"use_ssl"`
}

func SetS3Defaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".endpoint", "localhost:9000")
	v.SetDefault(prefix+".access_key_id", "")
	v.SetDefault(prefix+".secret_access_key", "")
	v.SetDefault(prefix+".secret_access_token", "")
	v.SetDefault(prefix+".bucket", "files")
	v.SetDefault(prefix+".region", "us-east-1")
	v.SetDefault(prefix+".use_ssl", false)
}
