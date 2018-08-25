package gocd

// MaterialAttribute describes the behaviour of the GoCD material structures for a pipeline
type MaterialAttribute interface {
	equal(attributes MaterialAttribute) (bool, error)
	GenerateGeneric() map[string]interface{}
	HasFilter() bool
	GetFilter() *MaterialFilter
}

// MaterialAttributesGit describes a git material
type MaterialAttributesGit struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	Branch string `json:"branch,omitempty"`

	SubmoduleFolder string `json:"submodule_folder,omitempty"`
	ShallowClone    bool   `json:"shallow_clone,omitempty"`

	Destination  string          `json:"destination,omitempty"`
	Filter       *MaterialFilter `json:"filter,omitempty"`
	InvertFilter bool            `json:"invert_filter"`
	AutoUpdate   bool            `json:"auto_update,omitempty"`
}

// MaterialAttributesSvn describes a material type
type MaterialAttributesSvn struct {
	Name              string `json:"name,omitempty"`
	URL               string `json:"url,omitempty"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"encrypted_password"`

	CheckExternals bool `json:"check_externals"`

	Destination  string          `json:"destination,omitempty"`
	Filter       *MaterialFilter `json:"filter,omitempty"`
	InvertFilter bool            `json:"invert_filter"`
	AutoUpdate   bool            `json:"auto_update,omitempty"`
}

// MaterialAttributesHg describes a Mercurial material type
type MaterialAttributesHg struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url"`

	Destination  string          `json:"destination"`
	Filter       *MaterialFilter `json:"filter,omitempty"`
	InvertFilter bool            `json:"invert_filter"`
	AutoUpdate   bool            `json:"auto_update,omitempty"`
}

// MaterialAttributesP4 describes a Perforce material type
type MaterialAttributesP4 struct {
	Name       string `json:"name,omitempty"`
	Port       string `json:"port"`
	UseTickets bool   `json:"use_tickets"`
	View       string `json:"view"`

	Username          string `json:"username"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"encrypted_password"`

	Destination  string          `json:"destination"`
	Filter       *MaterialFilter `json:"filter,omitempty"`
	InvertFilter bool            `json:"invert_filter"`
	AutoUpdate   bool            `json:"auto_update,omitempty"`
}

// MaterialAttributesTfs describes a Team Foundation Server material
// codebeat:disable[TOO_MANY_IVARS]
type MaterialAttributesTfs struct {
	Name string `json:"name,omitempty"`

	URL         string `json:"url"`
	ProjectPath string `json:"project_path"`
	Domain      string `json:"domain"`

	Username          string `json:"username"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"encrypted_password"`

	Destination  string          `json:"destination"`
	Filter       *MaterialFilter `json:"filter,omitempty"`
	InvertFilter bool            `json:"invert_filter"`
	AutoUpdate   bool            `json:"auto_update,omitempty"`
}
// codebeat:enable[TOO_MANY_IVARS]

// MaterialAttributesDependency describes a Pipeline dependency material
type MaterialAttributesDependency struct {
	Name       string `json:"name,omitempty"`
	Pipeline   string `json:"pipeline"`
	Stage      string `json:"stage"`
	AutoUpdate bool   `json:"auto_update,omitempty"`
}

// MaterialAttributesPackage describes a package reference
type MaterialAttributesPackage struct {
	Ref string `json:"ref"`
}

// MaterialAttributesPlugin describes a plugin material
type MaterialAttributesPlugin struct {
	Ref string `json:"ref"`

	Destination  string          `json:"destination"`
	Filter       *MaterialFilter `json:"filter,omitempty"`
	InvertFilter bool            `json:"invert_filter"`
}
