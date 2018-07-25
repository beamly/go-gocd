package gocd

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"sort"
)

type ServerApiVersionMapping struct {
	Api    string
	Server *version.Version
}

func NewServerApiVersionMappingSlice(mappings ...*ServerApiVersionMapping) []*ServerApiVersionMapping {
	return mappings
}

func NewServerApiVersionMapping(serverVersion, apiVersion string) (mapping *ServerApiVersionMapping) {
	mapping = &ServerApiVersionMapping{
		Api: apiVersion,
	}

	var err error
	if mapping.Server, err = version.NewVersion(serverVersion); err != nil {
		panic(err)
	}
	return
}

func NewServerVersionCollection() *ServerVersionCollection {
	return &ServerVersionCollection{
		mapping: make(map[string]*ServerVersionMethodEndpointMapping),
	}
}

func (svc *ServerVersionCollection) WithEndpoint(endpoint string, mapping *ServerVersionMethodEndpointMapping) *ServerVersionCollection {
	svc.mapping[endpoint] = mapping
	return svc
}

type ServerVersionCollection struct {
	mapping map[string]*ServerVersionMethodEndpointMapping
}

func (c *ServerVersionCollection) GetEndpointOk(endpoint string) (endpointMapping *ServerVersionMethodEndpointMapping, hasEndpoint bool) {
	endpointMapping, hasEndpoint = c.mapping[endpoint]
	return
}

func NewServerVersionMethodEndpointMapping() *ServerVersionMethodEndpointMapping {
	return &ServerVersionMethodEndpointMapping{
		methods: make(map[string][]*ServerApiVersionMapping),
	}
}

func (svmep *ServerVersionMethodEndpointMapping) WithMethod(method string, mappings []*ServerApiVersionMapping) *ServerVersionMethodEndpointMapping {
	svmep.methods[method] = mappings
	return svmep
}

type ServerVersionMethodEndpointMapping struct {
	methods map[string][]*ServerApiVersionMapping
}

func (m *ServerVersionMethodEndpointMapping) GetMappingOk(method string) (mappings *ServerApiVersionMappingCollection, hasMethod bool) {
	var mapping []*ServerApiVersionMapping
	mapping, hasMethod = m.methods[method]
	mappings = &ServerApiVersionMappingCollection{mappings: mapping}

	return
}

type ServerApiVersionMappingCollection struct {
	mappings []*ServerApiVersionMapping
}

// GetApiVersion for the highest common version
func (c *ServerApiVersionMappingCollection) GetApiVersion(versionParts *version.Version) (apiVersion string, err error) {
	c.Sort()

	lastMapping := c.mappings[0]
	for _, mapping := range c.mappings {
		if mapping.Server.GreaterThan(versionParts) || mapping.Server.Equal(versionParts) {
			return lastMapping.Api, nil
		}
		lastMapping = mapping
	}
	return "", fmt.Errorf("could not find api version")
}

func (c *ServerApiVersionMappingCollection) Sort() {
	sort.Sort(c)
}

func (c *ServerApiVersionMappingCollection) Len() int {
	return len(c.mappings)
}

func (c *ServerApiVersionMappingCollection) Less(i, j int) bool {
	return c.mappings[i].Server.LessThan(c.mappings[i].Server)
}

func (c *ServerApiVersionMappingCollection) Swap(i, j int) {
	c.mappings[i], c.mappings[j] = c.mappings[j], c.mappings[i]
}

// GetAPIVersion for a given endpoint and method
func (sv *ServerVersion) GetAPIVersion(endpoint string, method string) (apiVersion string, err error) {

	serverVersionLookup := NewServerVersionCollection().WithEndpoint(
		"/api/version", NewServerVersionMethodEndpointMapping().WithMethod(
			"GET", NewServerApiVersionMappingSlice(
				NewServerApiVersionMapping("1.0.0", apiV1),
			)))

	if methods, hasEndpoint := serverVersionLookup.GetEndpointOk(endpoint); hasEndpoint {
		if versionMapping, hasVersionMapping := methods.GetMappingOk(method); hasVersionMapping {
			return versionMapping.GetApiVersion(sv.VersionParts)
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
