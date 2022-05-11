package localizer

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

type testingUI struct {
	terminal.UI
	buf bytes.Buffer
}

func NewTestingUI(buf bytes.Buffer) *testingUI {
	return &testingUI{UI: terminal.NonInteractiveUI(context.Background()), buf: buf}
}

func (ui *testingUI) Output(msg string, raw ...interface{}) {
	ui.UI.Output(msg, terminal.WithWriter(&ui.buf), raw)
}

func TestNewCoreLocalizer(t *testing.T) {
	oldLC_ALL := os.Getenv("LC_ALL")
	os.Setenv("LC_ALL", "en")
	defer os.Setenv("LC_ALL", oldLC_ALL)

	var buf bytes.Buffer
	ui := NewTestingUI(buf)

	expectedMsg := "Adding box"
	l, err := NewCoreLocalizer(ui)
	require.NoError(t, err)
	require.NotNil(t, l)
	// Test output prints to terminal
	err = l.Output("box_add", nil)
	require.NoError(t, err)
	require.Contains(t, ui.buf.String(), expectedMsg)
	// Test localizing message
	msg, err := l.LocalizeMsg("box_add", nil)
	require.NoError(t, err)
	require.Contains(t, msg, expectedMsg)
}

func TestNewPluginLocalizer(t *testing.T) {
	oldLC_ALL := os.Getenv("LC_ALL")
	defer os.Setenv("LC_ALL", oldLC_ALL)

	var buf bytes.Buffer
	ui := NewTestingUI(buf)

	expectedMsgEN := "yes i am test data"
	jsonStrEN, err := json.Marshal(map[string]string{
		"testdata": expectedMsgEN,
	})
	require.NoError(t, err)

	expectedMsgES := "sí, soy datos de prueba"
	jsonStrES, err := json.Marshal(map[string]string{
		"testdata": expectedMsgES,
	})
	require.NoError(t, err)

	localeData := []LocaleData{
		{
			LocaleData: jsonStrEN,
			LocalePath: "stub.json",
			Languages:  []language.Tag{language.English, language.AmericanEnglish},
		},
		{
			LocaleData: jsonStrES,
			LocalePath: "stub.json",
			Languages:  []language.Tag{language.Spanish},
		},
	}

	os.Setenv("LC_ALL", "en")
	// Test localizing message to english
	l, err := NewPluginLocalizer(ui, localeData...)
	require.NoError(t, err)
	msg, err := l.LocalizeMsg("testdata", nil)
	require.NoError(t, err)
	require.Contains(t, msg, expectedMsgEN)

	os.Setenv("LC_ALL", "es")
	// Test localizing message to spanish
	l, err = NewPluginLocalizer(ui, localeData...)
	require.NoError(t, err)
	msg, err = l.LocalizeMsg("testdata", nil)
	require.NoError(t, err)
	require.Contains(t, msg, expectedMsgES)
}

func TestNewPluginLocalizerBadPath(t *testing.T) {
	var buf bytes.Buffer
	ui := NewTestingUI(buf)

	expectedMsgEN := "yes i am test data"
	jsonStrEN, err := json.Marshal(map[string]string{
		"testdata": expectedMsgEN,
	})
	require.NoError(t, err)

	localeData := []LocaleData{
		{
			LocaleData: jsonStrEN,
			LocalePath: "",
			Languages:  []language.Tag{language.English, language.AmericanEnglish},
		},
	}
	_, err = NewPluginLocalizer(ui, localeData...)
	require.Error(t, err)
}
