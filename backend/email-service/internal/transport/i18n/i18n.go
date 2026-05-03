package i18n

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed translations/*.json
var fs embed.FS

var bundle *i18n.Bundle

var SupportedLanguages []string
var initOnce sync.Once

func Init() {
	initOnce.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		SupportedLanguages = SupportedLanguages[:0]

		files, _ := fs.ReadDir("translations")
		for _, file := range files {
			name := file.Name()
			content, _ := fs.ReadFile("translations/" + name)
			bundle.ParseMessageFileBytes(content, name)

			lang := strings.TrimSuffix(name, filepath.Ext(name))
			SupportedLanguages = append(SupportedLanguages, lang)
		}
	})
}

func NewLocalizer(lang string) *i18n.Localizer {
	Init()

	if lang == "" || !isSupported(lang) {
		lang = "en"
	}
	return i18n.NewLocalizer(bundle, lang)
}

func isSupported(lang string) bool {
	for _, l := range SupportedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}
