package gocd

import (
	"fmt"
	"github.com/hashicorp/go-version"
)

// GetAPIVersion for a given endpoint and method
func (sv *ServerVersion) GetAPIVersion(endpoint string, method string) (apiVersion string, err error) {
	var hasEndpoint, hasMethod bool
	var methods map[string]string
	serverVersionLookup := map[string]interface{}{
		"/api/version": map[string]string{"GET": apiV1},
	}

	if methods, hasEndpoint = serverVersionLookup[endpoint].(map[string]string); hasEndpoint {
		if apiVersion, hasMethod = methods[method]; hasMethod {
			return apiVersion, nil
		}
	}

	return "", fmt.Errorf("could not find API version tag for '%s %s'", method, endpoint)
}

func (sv *ServerVersion) parseVersion() (err error) {
	sv.VersionParts, err = version.NewVersion(sv.Version)
	return
}

// Equal if the two versions are identical
func (sv *ServerVersion) Equal(v *ServerVersion) bool {
	return sv.Version == v.Version
}

// LessThan compares this server version and determines if it is older than the provided server version
func (sv *ServerVersion) LessThan(v *ServerVersion) bool {
	return sv.VersionParts.LessThan(v.VersionParts)
}
