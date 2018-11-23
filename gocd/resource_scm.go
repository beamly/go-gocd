package gocd

// SetVersion sets a version string for this SCM object
func (scm *SCM) SetVersion(version string) {
	scm.Version = version
}

// GetVersion retrieves a version string for this CM object
func (scm SCM) GetVersion() (version string) {
	return scm.Version
}

// RemoveLinks from the pipeline object for json marshalling.
func (scm *SCM) RemoveLinks() {
	scm.Links = nil
}

// GetLinks from pipeline
func (scm *SCM) GetLinks() *HALLinks {
	return scm.Links
}
