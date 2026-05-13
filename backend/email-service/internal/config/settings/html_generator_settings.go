package settings

import "github.com/spf13/viper"

type HTMLGeneratorSettings struct {
	BaseUrlForActivation    string `mapstructure:"base_url_for_activation"`
	BaseUrlForPasswordReset string `mapstructure:"base_url_for_password_reset"`
	BaseUrlForEmailChange   string `mapstructure:"base_url_for_email_change"`
}

func SetHTMLGeneratorDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".base_url_for_activation", "http://localhost:8080/activate")
	v.SetDefault(prefix+".base_url_for_password_reset", "http://localhost:8080/password-reset")
	v.SetDefault(prefix+".base_url_for_email_change", "http://localhost:8080/email-change")
}
