package gocd

import (
	"fmt"
	"strings"
	"strconv"
)

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

func (v *ServerVersion) parseVersion() (err error) {
	var major, minor, patch int
	versionParts := strings.Split(v.Version, ".")

	if major, err = strconv.Atoi(versionParts[0]); err != nil {
		return
	}

	if minor, err = strconv.Atoi(versionParts[1]); err != nil {
		return
	}

	if patch, err = strconv.Atoi(versionParts[2]); err != nil {
		return
	}

	v.VersionParts = &ServerVersionParts{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
	return
}

// Equal if the two versions are identical
func (sc *ServerVersion) Equal(v *ServerVersion) bool {
	return sc.Version == v.Version
}

func (sc *ServerVersion) LessThan(v *ServerVersion) bool {
	return sc.VersionParts.Major <= v.VersionParts.Major &&
		sc.VersionParts.Minor <= v.VersionParts.Minor &&
		sc.VersionParts.Patch < v.VersionParts.Patch
}
