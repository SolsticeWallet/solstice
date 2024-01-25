package i18n

import (
	"embed"
	"log/slog"

	"github.com/BurntSushi/toml"
	"github.com/jeandeaual/go-locale"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var debugMode bool
var bundle *goi18n.Bundle
var localizer *goi18n.Localizer

//go:embed locale/*.toml
var localeFs embed.FS

func init() {
	debugMode = false
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(localeFs, "locale/en.toml")
	bundle.LoadMessageFileFS(localeFs, "locale/nl.toml")
	// add extra languages here

	lang, err := locale.GetLanguage()
	if err != nil {
		panic(err)
	}
	lang = lang[:2]

	localizer = goi18n.NewLocalizer(bundle, lang, language.English.String())
}

func OverrideLanguage(lang string) {
	localizer = goi18n.NewLocalizer(bundle, lang, language.English.String())
}

func SetDebug(debug bool) {
	debugMode = debug
}

func Localize(key string, args ...string) string {
	if debugMode {
		return key
	}

	if len(args)%2 != 0 {
		slog.Error("i18n: invalid number of arguments")
		return key
	}
	templ := map[string]string{}
	for i := 0; i < len(args); i += 2 {
		templ[args[i]] = args[i+1]
	}

	translation, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templ,
	})
	if err != nil {
		slog.Error("i18n: failed to localize", "key", key, "error", err)
		return key
	}
	return translation
}

func LocalizeN(key string, n int, args ...string) string {
	if debugMode {
		return key
	}

	if len(args)%2 != 0 {
		slog.Error("i18n: invalid number of arguments")
		return key
	}
	templ := map[string]string{}
	for i := 0; i < len(args); i += 2 {
		templ[args[i]] = args[i+1]
	}

	translation, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templ,
		PluralCount:  n,
	})
	if err != nil {
		slog.Error("i18n: failed to localize", "key", key, "n", n, "error", err)
		return key
	}
	return translation
}

func T(key string, args ...string) string {
	return Localize(key, args...)
}

func Tn(key string, n int, args ...string) string {
	return LocalizeN(key, n, args...)
}
