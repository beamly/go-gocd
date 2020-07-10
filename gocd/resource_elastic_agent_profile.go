package gocd

// SetVersion sets a version string for this role
func (eap *ElasticAgentProfile) SetVersion(version string) {
	eap.Version = version
}

// GetVersion retrieves a version string for this role
func (eap ElasticAgentProfile) GetVersion() (version string) {
	return eap.Version
}

// RemoveLinks from the pipeline object for json marshalling.
func (eap *ElasticAgentProfile) RemoveLinks() {
	eap.Links = nil
}

// GetLinks from pipeline
func (eap *ElasticAgentProfile) GetLinks() *HALLinks {
	return eap.Links
}
