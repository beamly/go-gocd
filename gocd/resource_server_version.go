package gocd

import (
	"fmt"
	"strings"
	"strconv"
)

// GetApiVersion for a given endpoint and method
func (sv *ServerVersion) GetApiVersion(endpoint string, method string) (apiVersion string, err error) {
	var hasEndpoint, hasMethod bool
	var methods map[string]string
	serverVersionLookup := map[string]interface{}{

	}

	if methods, hasEndpoint = serverVersionLookup[endpoint].(map[string]string); hasEndpoint {
		if apiVersion, hasMethod = methods[method]; hasMethod {
			return
		}
	}

	return "", fmt.Errorf("could not find API version tag for '%s %s'", method, endpoint)
}

func (sv *ServerVersion) parseVersion() (err error) {
	var major, minor, patch int
	versionParts := strings.Split(sv.Version, ".")

	if major, err = strconv.Atoi(versionParts[0]); err != nil {
		return
	}

	if minor, err = strconv.Atoi(versionParts[1]); err != nil {
		return
	}

	if patch, err = strconv.Atoi(versionParts[2]); err != nil {
		return
	}

	sv.VersionParts = &ServerVersionParts{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
	return
}

// Equal if the two versions are identical
func (sv *ServerVersion) Equal(v *ServerVersion) bool {
	return sv.Version == v.Version
}

func (sv *ServerVersion) LessThan(v *ServerVersion) bool {
	return sv.VersionParts.Major <= v.VersionParts.Major &&
		sv.VersionParts.Minor <= v.VersionParts.Minor &&
		sv.VersionParts.Patch < v.VersionParts.Patch
}
