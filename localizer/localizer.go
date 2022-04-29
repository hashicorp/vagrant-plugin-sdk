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

func NewLocalizer(lang string, terminal terminal.UI) (localizer *Localizer, err error) {
	if lang == "" {
		lang = language.English.String()
	}
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localDataPath := "localizer/locales/" + lang + ".json"
	enData, err := Asset(localDataPath)
	if err != nil {
		return nil, err
	}
	_, err = bundle.ParseMessageFileBytes(enData, localDataPath)
	if err != nil {
		return nil, err
	}
	l := i18n.NewLocalizer(bundle, lang)
	return &Localizer{
		localizer: l, bundle: bundle, terminal: terminal,
	}, nil
}

func (l *Localizer) LocalizeMsg(msg string) error {
	config := i18n.LocalizeConfig{
		MessageID: msg,
	}
	localizedMsg, err := l.localizer.Localize(&config)
	if err != nil {
		return err
	}
	l.terminal.Output(localizedMsg)
	return nil
}
