package gocd

import (
	"errors"
	"fmt"
)

// ErrBadServerVersionRange describes a case where the minimum server version is greater than the maximum server version
var ErrBadServerVersionRange = errors.New("minimum server version must not be larger than maximum server version")

// ServerVersionRange describe the domain of Server versions and allows them to be compared
type ServerVersionRange struct {
	min *ServerVersion
	max *ServerVersion
}

// MustServerVersionRange creates a Server Version range and ignores errors
func MustServerVersionRange(min *ServerVersion, max *ServerVersion) (svr *ServerVersionRange) {
	svr, _ = NewServerVersionRange(min, max)
	return
}

// NewServerVersionRange creates a Server Version range
func NewServerVersionRange(min *ServerVersion, max *ServerVersion) (*ServerVersionRange, error) {
	min.parseVersion()
	max.parseVersion()
	if max.LessThan(min) {
		return nil, ErrBadServerVersionRange
	}
	return &ServerVersionRange{
		min: min,
		max: max,
	}, nil
}

// Contains checks if a server version falls within the bounds of the range
func (svr *ServerVersionRange) Contains(version *ServerVersion) bool {
	lowerBound := svr.min.LessThan(version)
	upperBound := version.LessThan(svr.max)
	equality := version.Equal(svr.max)
	return lowerBound && upperBound || equality
}

// String representation of the Server version
func (svr *ServerVersionRange) String() string {
	return fmt.Sprintf("(%s, %s]", svr.min.String(), svr.max.String())
}
