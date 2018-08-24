package gocd

import (
	"regexp"
	"strconv"
)

const (
	apiV0 = ""
	// Version 1 of the GoCD API.
	apiV1 = "application/vnd.go.cd.v1+json"
	// Version 2 of the GoCD API.
	apiV2 = "application/vnd.go.cd.v2+json"
	// Version 3 of the GoCD API.
	apiV3 = "application/vnd.go.cd.v3+json"
	// Version 4 of the GoCD API.
	apiV4 = "application/vnd.go.cd.v4+json"
	// Version 5 of the GoCD API.
	apiV5 = "application/vnd.go.cd.v5+json"
	// Version 6 of the GoCD API.
	apiV6 = "application/vnd.go.cd.v6+json"
)

// APIVersion describe the Accept header type of an API version
type APIVersion string

// APIVersionRange describe the domain of API versions and allows them to be compared
type APIVersionRange struct {
	min APIVersion
	max APIVersion
}

var (
	// APIVersion0 represents no Api Version
	APIVersion0 = APIVersion(apiV0)
	// APIVersion1 represents Api version 1
	APIVersion1 = APIVersion(apiV1)
	// APIVersion2 represents Api version 2
	APIVersion2 = APIVersion(apiV2)
	// APIVersion3 represents Api version 3
	APIVersion3 = APIVersion(apiV3)
	// APIVersion4 represents Api version 4
	APIVersion4 = APIVersion(apiV4)
	// APIVersion5 represents Api version 5
	APIVersion5 = APIVersion(apiV5)
	// APIVersion6 represents Api version 6
	APIVersion6     = APIVersion(apiV6)
	apiVersionRegex = regexp.MustCompile("application/vnd.go.cd.v(\\d+)\\+json")
)

// Int value for the api version
func (av APIVersion) Int() (i int) {
	if string(av) == "" {
		return 0
	}
	matches := apiVersionRegex.FindStringSubmatch(av.String())
	// We match for integers, so this shouldn't fail
	i, _ = strconv.Atoi(matches[1])
	return
}

// LessThan compares two api versions
func (av APIVersion) LessThan(version APIVersion) bool {
	return av.Int() < version.Int()
}

// Equal compares the quality of two api versions
func (av APIVersion) Equal(version APIVersion) bool {
	return string(av) == string(version)
}

// String representation of the Api version
func (av APIVersion) String() string {
	return string(av)
}

// NewAPIVersionRange creates an API range
func NewAPIVersionRange(min APIVersion, max APIVersion) *APIVersionRange {
	return &APIVersionRange{
		min: min,
		max: max,
	}
}

// Contains checks if an api version falls within the bounds of the range
func (avr *APIVersionRange) Contains(version APIVersion) bool {
	upperBound := version.LessThan(avr.max)
	lowerBound := avr.min.LessThan(version)
	equal := avr.max.Equal(version)
	return upperBound && lowerBound || equal
}
