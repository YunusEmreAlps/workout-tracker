package app

import (
	"fmt"
	"strings"

	"github.com/invopop/ctxi18n"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/jovandeginste/workout-tracker/v2/pkg/database"
	"github.com/jovandeginste/workout-tracker/v2/views/helpers"
	"github.com/labstack/echo/v5"
)

const (
	BrowserLanguage   = "browser"
	BrowserTheme      = "browser"
	DefaultTotalsShow = database.WorkoutTypeRunning
)

func (a *App) ConfigureLocalizer() error {
	if err := ctxi18n.LoadWithDefault(a.Translations, "en"); err != nil {
		return err
	}

	a.translator = ctxi18n.Match(string(ctxi18n.DefaultLocale))

	return nil
}

func getLanguagePrefix(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if idx := strings.IndexAny(lang, "-_"); idx != -1 {
		return lang[:idx]
	}

	// Backwards-compat: Norwegian was previously stored/negotiated as "no".
	if lang == "no" {
		lang = "nb"
	}

	return lang
}

func (a *App) langFromContextString(ctx *echo.Context) string {
	langs := langFromContext(ctx)
	resolved := parseAndResolveLanguages(langs)
	unique := deduplicateLanguages(resolved)

	return strings.Join(unique, ",")
}

func resolveLanguageCode(codeStr string) string {
	if codeStr == "" {
		return ""
	}

	// 1. Direct match check
	code := i18n.Code(codeStr)
	if ctxi18n.Get(code) != nil {
		return codeStr
	}

	// 2. Prefix match check
	reqPrefix := getLanguagePrefix(codeStr)

	for _, supTag := range helpers.SupportedLanguages() {
		sup := supTag.String()
		if getLanguagePrefix(sup) == reqPrefix {
			return sup
		}
	}
	return ""
}

func parseAndResolveLanguages(langs []any) []string {
	res := []string{}

	for _, lang := range langs {
		l, ok := lang.(string)
		if !ok || l == "" {
			continue
		}

		parsed := i18n.ParseAcceptLanguage(l)
		for _, code := range parsed {
			if match := resolveLanguageCode(string(code)); match != "" {
				res = append(res, match)
			}
		}
	}

	return res
}

func deduplicateLanguages(langs []string) []string {
	seen := make(map[string]bool)
	unique := []string{}

	for _, code := range langs {
		if !seen[code] {
			seen[code] = true
			unique = append(unique, code)
		}
	}

	return unique
}

func langFromContext(ctx *echo.Context) []any {
	return []any{
		ctx.QueryParam("lang"),
		ctx.Get("user_language"),
		ctx.Request().Header.Get("Accept-Language"),
	}
}

func (a *App) i18nN(ctx *echo.Context, message string, count int, vars ...any) string {
	t := a.translatorFromContext(ctx)
	if t.Has(message) {
		return t.N(message, count, vars...)
	}

	return fmt.Sprintf("%s[%d]: %v", message, count, vars)
}

func (a *App) i18nT(ctx *echo.Context, message string, vars ...any) string {
	t := a.translatorFromContext(ctx)
	if t.Has(message) {
		return t.T(message, vars...)
	}

	return fmt.Sprintf("%s: %v", message, vars)
}

func (a *App) translatorFromContext(ctx *echo.Context) *i18n.Locale {
	if l := ctxi18n.Locale(ctx.Request().Context()); l != nil {
		return l
	}

	return a.translator
}
