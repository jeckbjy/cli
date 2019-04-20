package cli

import (
	"regexp"
	"strconv"
	"strings"
)

func strToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return val
}

func strToUint(s string) uint {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}

	return uint(val)
}

func strToInt64(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return val
}

func strToUint64(s string) uint64 {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}

	return val
}

func strToBool(s string) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}

	return val
}

func strToF32(s string) float32 {
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}

	return float32(val)
}

func strToF64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return float64(val)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// toSnakeCase return xx_xx_xx
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// toKebabCase return xx-xx-xx
func toKebabCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}
