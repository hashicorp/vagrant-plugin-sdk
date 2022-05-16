package localizer

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewCoreLocalizer(t *testing.T) {
	oldLC_ALL := os.Getenv("LC_ALL")
	os.Setenv("LC_ALL", "en")
	defer os.Setenv("LC_ALL", oldLC_ALL)

	expectedMsg := "Adding box"
	l, err := NewCoreLocalizer()
	require.NoError(t, err)
	require.NotNil(t, l)
	// Test output prints to terminal
	msg, err := l.LocalizeMsg("box_add", nil)
	require.NoError(t, err)
	require.Contains(t, msg, expectedMsg)
}

func TestNewPluginLocalizer(t *testing.T) {
	oldLC_ALL := os.Getenv("LC_ALL")
	defer os.Setenv("LC_ALL", oldLC_ALL)

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
	l, err := NewPluginLocalizer(localeData...)
	require.NoError(t, err)
	msg, err := l.LocalizeMsg("testdata", nil)
	require.NoError(t, err)
	require.Contains(t, msg, expectedMsgEN)

	os.Setenv("LC_ALL", "es")
	// Test localizing message to spanish
	l, err = NewPluginLocalizer(localeData...)
	require.NoError(t, err)
	msg, err = l.LocalizeMsg("testdata", nil)
	require.NoError(t, err)
	require.Contains(t, msg, expectedMsgES)
}

func TestNewPluginLocalizerBadPath(t *testing.T) {
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
	_, err = NewPluginLocalizer(localeData...)
	require.Error(t, err)
}
