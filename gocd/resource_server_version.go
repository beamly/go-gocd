package gocd

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"sort"
)

var serverVersionLookup *serverVersionCollection

func init() {
	serverVersionLookup = &serverVersionCollection{
		mapping: map[endpointS]*serverAPIVersionMappingCollection{
			"/api/version": newVersionCollection(
				newServerAPI("16.6.0", apiV1)),
			"/api/admin/pipelines/:pipeline_name": newVersionCollection(
				newServerAPI("18.7.0", apiV6),
				newServerAPI("17.12.0", apiV5),
				newServerAPI("17.4.0", apiV4)),
		},
	}
}

// GetAPIVersion for a given endpoint and method
func (sv *ServerVersion) GetAPIVersion(endpoint string) (apiVersion string, err error) {

	if versions, hasEndpoint := serverVersionLookup.GetEndpointOk(endpoint); hasEndpoint {
		return versions.GetAPIVersion(sv.VersionParts)
	}

	return "", fmt.Errorf("could not find API version tag for '%s'", endpoint)
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
// Structures for storing, creating, and parsing the endpoint/server-version/api-version mapping
//

// following type definitions makes the map[...]... below a bit easier to understand.
type endpointS string

// serverVersionToAcceptMapping links an Accept header value and a Server version
type serverVersionToAcceptMapping struct {
	API    string
	Server *version.Version
}

type serverVersionCollection struct {
	mapping map[endpointS]*serverAPIVersionMappingCollection
}

type serverAPIVersionMappingCollection struct {
	mappings []*serverVersionToAcceptMapping
}

// newServerAPISlice provides some syntactic sugar to make the chaining resources a bit easier
// to read.
func newVersionCollection(mappings ...*serverVersionToAcceptMapping) *serverAPIVersionMappingCollection {
	return &serverAPIVersionMappingCollection{
		mappings: mappings,
	}
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

func (svc *serverVersionCollection) GetEndpointOk(endpoint string) (endpointMapping *serverAPIVersionMappingCollection, hasEndpoint bool) {
	endpointMapping, hasEndpoint = svc.mapping[endpointS(endpoint)]
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
	return "", fmt.Errorf("could not find api version for server version '%s'", versionParts.String())
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
