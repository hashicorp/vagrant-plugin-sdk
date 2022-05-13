package localizer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	DEFAULT_LANGUAGE = "en"
)

type LocaleData struct {
	LocaleData []byte
	LocalePath string
	Languages  []language.Tag
}

type Localizer struct {
	localizer *i18n.Localizer
	bundle    *i18n.Bundle
	terminal  terminal.UI
}

func NewPluginLocalizer(terminal terminal.UI, data ...LocaleData) (localizer *Localizer, err error) {
	lang, err := getLang()
	if err != nil {
		return nil, err
	}
	bundle, err := loadLocaleData(lang, data...)
	if err != nil {
		return nil, err
	}
	l := i18n.NewLocalizer(bundle, lang.String())
	return &Localizer{
		localizer: l, bundle: bundle, terminal: terminal,
	}, nil
}

func NewCoreLocalizer(terminal terminal.UI) (localizer *Localizer, err error) {
	lang, err := getLang()
	if err != nil {
		return nil, err
	}
	localDataPathEN := "localizer/locales/en.json"
	localeData, err := Asset(localDataPathEN)
	if err != nil {
		return nil, err
	}
	bundle, err := loadLocaleData(lang, LocaleData{
		LocaleData: localeData,
		LocalePath: localDataPathEN,
		Languages:  []language.Tag{language.English, language.AmericanEnglish, language.BritishEnglish},
	})
	if err != nil {
		return nil, err
	}
	l := i18n.NewLocalizer(bundle, lang.String())
	return &Localizer{
		localizer: l, bundle: bundle, terminal: terminal,
	}, nil
}

func (l *Localizer) LocalizeMsg(msg string, templateData interface{}) (string, error) {
	config := i18n.LocalizeConfig{
		MessageID:    msg,
		TemplateData: templateData,
	}
	return l.localizer.Localize(&config)
}

func (l *Localizer) Output(msg string, templateData interface{}) error {
	config := i18n.LocalizeConfig{
		MessageID:    msg,
		TemplateData: templateData,
	}
	localizedMsg, err := l.localizer.Localize(&config)
	if err != nil {
		return err
	}
	l.terminal.Output(localizedMsg)
	return nil
}

func getLang() (locale language.Tag, err error) {
	// TODO: make this work for Windows too
	lang := os.Getenv("LC_ALL")
	if lang == "" {
		langVar := os.Getenv("LANG")
		// LANG returns the locale and encoding. Extract just the locale
		lang = strings.Split(langVar, ".")[0]
	}
	if lang == "" {
		lang = DEFAULT_LANGUAGE
	}
	return language.Parse(lang)
}

func loadLocaleData(lang language.Tag, data ...LocaleData) (bundle *i18n.Bundle, err error) {
	bundle = i18n.NewBundle(lang)
	for _, d := range data {
		msg, err := i18n.ParseMessageFileBytes(
			d.LocaleData,
			d.LocalePath,
			map[string]i18n.UnmarshalFunc{"json": json.Unmarshal},
		)
		if err != nil {
			return nil, err
		}
		for _, lang := range d.Languages {
			if err := bundle.AddMessages(lang, msg.Messages...); err != nil {
				return nil, err
			}
		}
	}
	return
}
