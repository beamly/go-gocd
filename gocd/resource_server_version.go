package gocd

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"sort"
)

var serverVersionLookup *serverVersionCollection

func init() {
	serverVersionLookup = newServerVersionLookup().WithEndpoint(
		"/api/version", newMethodToVersionsMapping().WithMethod(
			"GET", newServerAPISlice(
				newServerAPI("1.0.0", apiV1),
			)))
}

// GetAPIVersion for a given endpoint and method
func (sv *ServerVersion) GetAPIVersion(endpoint string, method string) (apiVersion string, err error) {

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

//
// Structures for storing, creating, and parsing the method/endpoint/server-version/api-version mapping
//

// following type definitions makes the map[...]... below a bit easier to understand.
type endpointS string
type methodS string

// serverVersionToAcceptMapping links an Accept header value and a Server version
type serverVersionToAcceptMapping struct {
	API    string
	Server *version.Version
}

type serverVersionMethodEndpointMapping struct {
	methods map[methodS][]*serverVersionToAcceptMapping
}

type serverVersionCollection struct {
	mapping map[endpointS]*serverVersionMethodEndpointMapping
}

type serverAPIVersionMappingCollection struct {
	mappings []*serverVersionToAcceptMapping
}

// newServerAPISlice provides some syntactic sugar to make the chaining resources a bit easier
// to read.
func newServerAPISlice(mappings ...*serverVersionToAcceptMapping) []*serverVersionToAcceptMapping {
	return mappings
}

// newServerAPI creates a new server/api version mapping and panics on any errors. These
// values will be hardcoded, so it should fail when loaded.
func newServerAPI(serverVersion, apiVersion string) (mapping *serverVersionToAcceptMapping) {
	mapping = &serverVersionToAcceptMapping{
		API: apiVersion,
	}

	var err error
	if mapping.Server, err = version.NewVersion(serverVersion); err != nil {
		panic(err)
	}
	return
}

func newServerVersionLookup() *serverVersionCollection {
	return &serverVersionCollection{
		mapping: make(map[endpointS]*serverVersionMethodEndpointMapping),
	}
}

func (svc *serverVersionCollection) WithEndpoint(endpoint string, mapping *serverVersionMethodEndpointMapping) *serverVersionCollection {
	svc.mapping[endpointS(endpoint)] = mapping
	return svc
}

func (svc *serverVersionCollection) GetEndpointOk(endpoint string) (endpointMapping *serverVersionMethodEndpointMapping, hasEndpoint bool) {
	endpointMapping, hasEndpoint = svc.mapping[endpointS(endpoint)]
	return
}

func newMethodToVersionsMapping() *serverVersionMethodEndpointMapping {
	return &serverVersionMethodEndpointMapping{
		methods: make(map[methodS][]*serverVersionToAcceptMapping),
	}
}

func (svmep *serverVersionMethodEndpointMapping) WithMethod(method string, mappings []*serverVersionToAcceptMapping) *serverVersionMethodEndpointMapping {
	svmep.methods[methodS(method)] = mappings
	return svmep
}

func (svmep *serverVersionMethodEndpointMapping) GetMappingOk(method string) (mappings *serverAPIVersionMappingCollection, hasMethod bool) {
	var mapping []*serverVersionToAcceptMapping
	mapping, hasMethod = svmep.methods[methodS(method)]
	mappings = &serverAPIVersionMappingCollection{mappings: mapping}

	return
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

// Sort the version collections
func (c *serverAPIVersionMappingCollection) Sort() {
	sort.Sort(c)
}

// Len of the versions in this collection.
func (c *serverAPIVersionMappingCollection) Len() int {
	return len(c.mappings)
}

// Less compares two server versions to see which is lower.
func (c *serverAPIVersionMappingCollection) Less(i, j int) bool {
	return c.mappings[i].Server.LessThan(c.mappings[j].Server)
}

// Swap the position of two server versions.
func (c *serverAPIVersionMappingCollection) Swap(i, j int) {
	c.mappings[i], c.mappings[j] = c.mappings[j], c.mappings[i]
}
