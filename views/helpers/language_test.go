package helpers

import (
	"testing"

	"github.com/jovandeginste/workout-tracker/v2/pkg/templatehelpers"
	"github.com/stretchr/testify/assert"
)

func TestSupportedLanguagesHaveFlags(t *testing.T) {
	for _, lang := range SupportedLanguages() {
		flag := templatehelpers.LanguageToFlag(lang.String())
		assert.NotEqual(t, "👽", flag, "Language %s should have a flag mapped", lang.String())
	}
}
