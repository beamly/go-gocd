package gocd

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"sort"
)

type serverAPIVersionMapping struct {
	API    string
	Server *version.Version
}

func newServerAPIVersionMappingSlice(mappings ...*serverAPIVersionMapping) []*serverAPIVersionMapping {
	return mappings
}

func newServerAPIVersionMapping(serverVersion, apiVersion string) (mapping *serverAPIVersionMapping) {
	mapping = &serverAPIVersionMapping{
		API: apiVersion,
	}

	var err error
	if mapping.Server, err = version.NewVersion(serverVersion); err != nil {
		panic(err)
	}
	return
}

func newServerVersionCollection() *serverVersionCollection {
	return &serverVersionCollection{
		mapping: make(map[string]*serverVersionMethodEndpointMapping),
	}
}

type serverVersionCollection struct {
	mapping map[string]*serverVersionMethodEndpointMapping
}

func (svc *serverVersionCollection) WithEndpoint(endpoint string, mapping *serverVersionMethodEndpointMapping) *serverVersionCollection {
	svc.mapping[endpoint] = mapping
	return svc
}
func (svc *serverVersionCollection) GetEndpointOk(endpoint string) (endpointMapping *serverVersionMethodEndpointMapping, hasEndpoint bool) {
	endpointMapping, hasEndpoint = svc.mapping[endpoint]
	return
}

func newServerVersionMethodEndpointMapping() *serverVersionMethodEndpointMapping {
	return &serverVersionMethodEndpointMapping{
		methods: make(map[string][]*serverAPIVersionMapping),
	}
}

type serverVersionMethodEndpointMapping struct {
	methods map[string][]*serverAPIVersionMapping
}

func (svmep *serverVersionMethodEndpointMapping) WithMethod(method string, mappings []*serverAPIVersionMapping) *serverVersionMethodEndpointMapping {
	svmep.methods[method] = mappings
	return svmep
}

func (svmep *serverVersionMethodEndpointMapping) GetMappingOk(method string) (mappings *serverAPIVersionMappingCollection, hasMethod bool) {
	var mapping []*serverAPIVersionMapping
	mapping, hasMethod = svmep.methods[method]
	mappings = &serverAPIVersionMappingCollection{mappings: mapping}

	return
}

type serverAPIVersionMappingCollection struct {
	mappings []*serverAPIVersionMapping
}

// GetAPIVersion for the highest common version
func (c *serverAPIVersionMappingCollection) GetAPIVersion(versionParts *version.Version) (apiVersion string, err error) {
	c.Sort()

	lastMapping := c.mappings[0]
	for _, mapping := range c.mappings {
		if mapping.Server.GreaterThan(versionParts) || mapping.Server.Equal(versionParts) {
			return lastMapping.API, nil
		}
		lastMapping = mapping
	}
	return "", fmt.Errorf("could not find api version")
}

func (c *serverAPIVersionMappingCollection) Sort() {
	sort.Sort(c)
}

func (c *serverAPIVersionMappingCollection) Len() int {
	return len(c.mappings)
}

func (c *serverAPIVersionMappingCollection) Less(i, j int) bool {
	return c.mappings[i].Server.LessThan(c.mappings[j].Server)
}

func (c *serverAPIVersionMappingCollection) Swap(i, j int) {
	c.mappings[i], c.mappings[j] = c.mappings[j], c.mappings[i]
}

// GetAPIVersion for a given endpoint and method
func (sv *ServerVersion) GetAPIVersion(endpoint string, method string) (apiVersion string, err error) {

	serverVersionLookup := newServerVersionCollection().WithEndpoint(
		"/api/version", newServerVersionMethodEndpointMapping().WithMethod(
			"GET", newServerAPIVersionMappingSlice(
				newServerAPIVersionMapping("1.0.0", apiV1),
			)))

	if methods, hasEndpoint := serverVersionLookup.GetEndpointOk(endpoint); hasEndpoint {
		if versionMapping, hasVersionMapping := methods.GetMappingOk(method); hasVersionMapping {
			return versionMapping.GetAPIVersion(sv.VersionParts)
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
