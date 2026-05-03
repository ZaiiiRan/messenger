package utils

import (
	"regexp"
	"strings"
)

type ParsedUA struct {
	Browser string
	OS      string
}

var (
	reEdge          = regexp.MustCompile(`Edg(?:e|A|iOS)?/(\d+)`)
	reOpera         = regexp.MustCompile(`OPR/(\d+)`)
	reYaBrowser     = regexp.MustCompile(`YaBrowser/(\d+)`)
	reSamsungBrowser = regexp.MustCompile(`SamsungBrowser/(\d+)`)
	reUCBrowser     = regexp.MustCompile(`UCBrowser/(\d+)`)
	reCriOS         = regexp.MustCompile(`CriOS/(\d+)`)
	reFxiOS         = regexp.MustCompile(`FxiOS/(\d+)`)
	reChrome        = regexp.MustCompile(`Chrome/(\d+)`)
	reFirefox       = regexp.MustCompile(`Firefox/(\d+)`)
	reSafari        = regexp.MustCompile(`Version/(\d+).*Safari`)
	reMSIE          = regexp.MustCompile(`MSIE (\d+)`)

	reAndroid = regexp.MustCompile(`Android (\d+(?:\.\d+)?)`)
	reIPad    = regexp.MustCompile(`iPad.*OS (\d+)[_.](\d+)`)
	reiOS     = regexp.MustCompile(`(?:iPhone OS|CPU OS) (\d+)[_.](\d+)`)
	reMacOS   = regexp.MustCompile(`Mac OS X (\d+)[_.](\d+)`)
	reWindows = regexp.MustCompile(`Windows NT (\d+\.\d+)`)
)

func ParseUserAgent(ua string) ParsedUA {
	return ParsedUA{
		Browser: parseBrowser(ua),
		OS:      parseOS(ua),
	}
}

func parseBrowser(ua string) string {
	switch {
	case reEdge.MatchString(ua):
		return "Edge " + reEdge.FindStringSubmatch(ua)[1]
	case reOpera.MatchString(ua):
		return "Opera " + reOpera.FindStringSubmatch(ua)[1]
	case reYaBrowser.MatchString(ua):
		return "Yandex Browser " + reYaBrowser.FindStringSubmatch(ua)[1]
	case reSamsungBrowser.MatchString(ua):
		return "Samsung Browser " + reSamsungBrowser.FindStringSubmatch(ua)[1]
	case reUCBrowser.MatchString(ua):
		return "UC Browser " + reUCBrowser.FindStringSubmatch(ua)[1]
	case reCriOS.MatchString(ua):
		return "Chrome " + reCriOS.FindStringSubmatch(ua)[1]
	case reFxiOS.MatchString(ua):
		return "Firefox " + reFxiOS.FindStringSubmatch(ua)[1]
	case reChrome.MatchString(ua):
		return "Chrome " + reChrome.FindStringSubmatch(ua)[1]
	case reFirefox.MatchString(ua):
		return "Firefox " + reFirefox.FindStringSubmatch(ua)[1]
	case reSafari.MatchString(ua):
		return "Safari " + reSafari.FindStringSubmatch(ua)[1]
	case reMSIE.MatchString(ua):
		return "Internet Explorer " + reMSIE.FindStringSubmatch(ua)[1]
	case strings.Contains(ua, "Trident/"):
		return "Internet Explorer 11"
	default:
		return "Unknown"
	}
}

func parseOS(ua string) string {
	switch {
	case reAndroid.MatchString(ua):
		return "Android " + reAndroid.FindStringSubmatch(ua)[1]
	case reIPad.MatchString(ua):
		m := reIPad.FindStringSubmatch(ua)
		return "iOS " + m[1] + "." + m[2]
	case reiOS.MatchString(ua):
		m := reiOS.FindStringSubmatch(ua)
		return "iOS " + m[1] + "." + m[2]
	case reMacOS.MatchString(ua):
		m := reMacOS.FindStringSubmatch(ua)
		return "macOS " + m[1] + "." + m[2]
	case reWindows.MatchString(ua):
		return "Windows " + windowsNTVersion(reWindows.FindStringSubmatch(ua)[1])
	case strings.Contains(ua, "CrOS"):
		return "Chrome OS"
	case strings.Contains(ua, "Linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}

func windowsNTVersion(nt string) string {
	switch nt {
	case "10.0":
		return "10/11"
	case "6.3":
		return "8.1"
	case "6.2":
		return "8"
	case "6.1":
		return "7"
	case "6.0":
		return "Vista"
	case "5.1", "5.2":
		return "XP"
	default:
		return nt
	}
}
