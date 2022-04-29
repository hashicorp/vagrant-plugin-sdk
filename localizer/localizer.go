package localizer

import (
	"encoding/json"

	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer struct {
	localizer *i18n.Localizer
	bundle    *i18n.Bundle
	terminal  terminal.UI
}

func NewPluginLocalizer(lang language.Tag, localeData []byte, localePath string, terminal terminal.UI) (localizer *Localizer, err error) {
	bundle := i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = bundle.ParseMessageFileBytes(localeData, localePath)
	if err != nil {
		return nil, err
	}
	l := i18n.NewLocalizer(bundle, lang.String())
	return &Localizer{
		localizer: l, bundle: bundle, terminal: terminal,
	}, nil
}

func NewCoreLocalizer(lang language.Tag, terminal terminal.UI) (localizer *Localizer, err error) {
	bundle := i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localDataPath := "localizer/locales/" + lang.String() + ".json"
	localeData, err := Asset(localDataPath)
	if err != nil {
		return nil, err
	}
	_, err = bundle.ParseMessageFileBytes(localeData, localDataPath)
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
