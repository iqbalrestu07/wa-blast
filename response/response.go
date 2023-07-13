// Package response contains response structure that will be used by handler to show response in JSON.
package response

import (
	"strconv"

	"wa-blast/configs"
)

// Success represents base response structure if a request is success
type Success struct {
	Result   interface{}       `json:"result"`
	Metadata *Metadata         `json:"_metadata,omitempty"`
	Header   map[string]string `json:"-"`
}

// Metadata represents base response structure for
type Metadata struct {
	Query           string   `json:"query,omitempty"`
	Count           int      `json:"count"`
	Skip            int64    `json:"skip"`
	Limit           int8     `json:"limit"`
	Order           string   `json:"order,omitempty"`
	Since           int64    `json:"since,omitempty"`
	Timestamp       int64    `json:"last_updated,omitempty"`
	LocationLat     *float32 `json:"location_lat,omitempty"`
	LocationLng     *float32 `json:"location_lng,omitempty"`
	SortField       string   `json:"-"`
	SortArrangement string   `json:"-"`
}

// JWT represents response for jwt authentication
type JWT struct {
	AccessToken string
	ExpiredAt   int64
}

// NewSuccess returns new success response
// If no body argument is set, return Success message string
// If only one body argument is set, set first argument as result body
// Else, returns first argument as result body, cast second argument as metadata and set in metadata response
func NewSuccess(body ...interface{}) *Success {
	switch len(body) {
	case 0:
		return &Success{Result: "Success"}
	case 1:
		return &Success{Result: body[0]}
	default:
		{
			// Cast second parameter as metadata
			metadata, _ := body[1].(*Metadata)
			// Generate response with metadata
			return &Success{Result: body[0], Metadata: metadata}
		}
	}
}

func NewToken(accessToken string, expiredAt int64) map[string]string {
	return map[string]string{
		"x-token": accessToken,
		"expired": strconv.FormatInt(expiredAt, 10),
	}
}

// GetResourceUrl returns generated url to a static file
func GetResourceUrl(resource string, fileName string) string {
	// If file name is empty, return empty url
	if fileName == "" {
		return ""
	}
	// Get base url from config
	baseUrl := configs.MustGetString("base_url." + resource)
	return baseUrl + fileName
}
