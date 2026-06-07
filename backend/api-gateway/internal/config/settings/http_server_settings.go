package settings

import "github.com/spf13/viper"

type HTTPServerSettings struct {
	Port string `mapstructure:"port"`

	ReadTimeout  uint `mapstructure:"read_timeout"`
	WriteTimeout uint `mapstructure:"write_timeout"`
	IdleTimeout  uint `mapstructure:"idle_timeout"`
	BodyLimit    int  `mapstructure:"body_limit"`

	TrustedProxyHeader string `mapstructure:"trusted_proxy_header"`
}

func SetHTTPServerDefaults(v *viper.Viper, prefix string, defaultPort string) {
	v.SetDefault(prefix+".port", defaultPort)
	v.SetDefault(prefix+".read_timeout", 15)
	v.SetDefault(prefix+".write_timeout", 15)
	v.SetDefault(prefix+".idle_timeout", 60)
	v.SetDefault(prefix+".body_limit", 4*1024*1024)
	// Empty by default: c.IP() uses the socket remote address. Set this to a
	// header name (e.g. "X-Real-IP") only when the gateway sits behind a
	// trusted reverse proxy that populates it.
	v.SetDefault(prefix+".trusted_proxy_header", "")
}
