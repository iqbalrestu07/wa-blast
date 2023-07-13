package util

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// Init initialize utility functions. App will exit if an error is occurred while initializing utilities
func Init() {
	// Init error codes
	initErrorCodes()
	initSnowflake()
	initRandomTime()
}

// An ApiError represent standard API error that contains HTTP Status (Status) and API-scoped Error Code (Code).
type ApiError struct {
	Status  int    `yaml:"status" json:"-"`
	Code    string `yaml:"code" json:"code"`
	Message string `yaml:"message" json:"message"`
}

// Error is an implementation of built-in error type interface
func (e ApiError) Error() string {
	return e.Message
}

// Error Codes from error_codes.yml
var errorCodes = make(map[string]ApiError)

// Errors
var ErrBadRequest, ErrInternalServer, ErrUserNotFound ApiError

// initErrorCodes load error codes from file. App will exit when an error occurred.
func initErrorCodes() {
	// Read file
	bytes, err := ioutil.ReadFile("errors_codes.yml")
	if err != nil {
		fmt.Printf("Unable to read error_codes.yml file. Error: %s\n", err.Error())
		os.Exit(5)
	}
	// Parse error codes file
	err = yaml.Unmarshal(bytes, &errorCodes)
	if err != nil {
		fmt.Printf("Unable to parse error_codes.yml file. Error: %s\n", err.Error())
		os.Exit(6)
	}
	// Init bad request
	ErrBadRequest = NewError("400")
	ErrInternalServer = NewError("500")
	ErrUserNotFound = NewError("-1002")
}

// NewError returns ApiError that is defined in error_codes.yml
func NewError(code string) ApiError {
	return errorCodes[code]
}

// CastError cast error interface as an ApiError
func CastError(err error) ApiError {
	apiErr, ok := err.(ApiError)
	if !ok {
		// If assert type fail, create new internal error
		apiErr = NewError("500")
	}
	return apiErr
}

func ParseFloat32(input string, defaultValue float32) float32 {
	o, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return defaultValue
	}
	return float32(o)
}

func ParseFloat64(input string, defaultValue float64) float64 {
	o, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return defaultValue
	}
	return o
}

func ParseInt(input string, defaultValue int) int {
	o, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return o
}

func ParseInt8(input string, defaultValue int8) int8 {
	o, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		return defaultValue
	}
	return int8(o)
}

func ParseInt16(input string, defaultValue int16) int16 {
	o, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		return defaultValue
	}
	return int16(o)
}

func ParseInt64(input string, defaultValue int64) int64 {
	o, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return defaultValue
	}
	return o
}

// GetNullStringPtr Converts sql.NullString to *string so string will be serialized to null in json instead sql.NullString
func GetNullStringPtr(str *sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &str.String
}

func UniqueString(Slice []string) []string {

	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range Slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func initRandomTime() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func LeftZeroPad(number, padWidth int64) string {
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", padWidth), number)
}

// match html tag and replace it with ""
func RemoveHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func ClearNumber(number string) string {
	if string(number[0]) == "+" {
		return string(number[1:len(number)])
	} else if string(number[0]) == "0" {
		return "62" + string(number[1:len(number)])
	}

	return number
}

func TypeTarget(number string) string {
	numberTarget := ClearNumber(number)
	if len([]rune(numberTarget)) > 15 {
		return "g.us"
	}
	return "s.whatsapp.net"
}
