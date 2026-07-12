package templatehelpers

import (
	"math"
	"strings"
	"time"

	"github.com/biter777/countries"
	"github.com/spf13/cast"
)

const InvalidValue = "N/A"

func HumanCadence(cad float64) string {
	return RoundFloat64(cad)
}

func HumanPower(pow float64) string {
	return RoundFloat64(pow)
}

func HumanCalories(cal float64) string {
	return RoundFloat64(cal)
}

func NumericDuration(d time.Duration) float64 {
	return d.Seconds()
}

var languageToCountryCodes = map[string][]string{
	"nl":      {"NL", "BE"},
	"en":      {"US", "GB"},
	"fr":      {"FR", "BE"},
	"de":      {"DE", "AT", "CH"},
	"es":      {"ES"},
	"it":      {"IT"},
	"pl":      {"PL"},
	"pt":      {"PT"},
	"pt-br":   {"BR"},
	"ru":      {"RU"},
	"sv":      {"SE"},
	"tr":      {"TR"},
	"ota":     {"TR"},
	"fi":      {"FI"},
	"fa":      {"IR"},
	"id":      {"ID"},
	"nb":      {"NO"},
	"nb-no":   {"NO"},
	"no":      {"NO"},
	"zh":      {"CN"},
	"zh-hans": {"CN"},
}

func LanguageToFlag(code string) string {
	code = strings.ToLower(strings.ReplaceAll(code, "_", "-"))

	ccs, found := languageToCountryCodes[code]
	if !found {
		return "👽"
	}

	flags := make([]string, 0, len(ccs))
	for _, cc := range ccs {
		flags = append(flags, CountryToFlag(cc))
	}

	result := strings.Join(flags, "")

	if code == "en" {
		result = "🌐" + result
	}

	return result
}

func CountryToFlag(cc string) string {
	ccc := countries.ByName(cc)
	return ccc.Emoji()
}

func HumanElevationFor(unit string) func(float64) string {
	switch unit {
	case "ft":
		return HumanElevationFt
	default:
		return HumanElevationM
	}
}

func HumanHeightSingleFor(unit string) func(float64) string {
	switch unit {
	case "in":
		return HumanHeightInch
	default:
		return HumanHeightCMNoSuffix
	}
}

func HumanHeightFor(unit string) func(float64) string {
	switch unit {
	case "in":
		return HumanHeightFeetInch
	default:
		return HumanHeightCM
	}
}

func HumanWeightFor(unit string) func(float64) string {
	switch unit {
	case "lbs":
		return HumanWeightPounds
	default:
		return HumanWeightKG
	}
}

func HumanDistanceFor(unit string) func(float64) string {
	switch unit {
	case "mi":
		return HumanDistanceMile
	case "nm":
		return HumanDistanceNM
	default:
		return HumanDistanceKM
	}
}

func HumanSpeedFor(unit string) func(float64) string {
	switch unit {
	case "mph":
		return HumanSpeedMilePH
	case "kn":
		return HumanSpeedKnots
	default:
		return HumanSpeedKPH
	}
}

func HumanTempoFor(unit string) func(float64) string {
	switch unit {
	case "min/mi", "mi":
		return HumanTempoMile
	case "min/nm", "nm":
		return HumanTempoNM
	default:
		return HumanTempoKM
	}
}

func HeightToDatabase(v float64, unit string) float64 {
	switch unit {
	case "in":
		return v * CmPerInch
	default:
		return v
	}
}

func WeightToDatabase(v float64, unit string) float64 {
	switch unit {
	case "lbs":
		return v / PoundsPerKG
	default:
		return v
	}
}

func DistanceToDatabase(v float64, unit string) float64 {
	switch unit {
	case "mi":
		return v * MeterPerMile
	case "m":
		return v
	case "nm":
		return v * MeterPerNM
	default:
		return v * MeterPerKM
	}
}

// RoundFloat64 rounds a float64 to 2 decimal places and returns it as a string.
func RoundFloat64(f float64) string {
	return cast.ToString(math.Round(f*100) / 100)
}
