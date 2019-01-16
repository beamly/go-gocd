package gocd_test

import (
	"context"
	"fmt"
	"github.com/beamly/go-gocd/gocd"
)

func ExampleAgentsService_List() {
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	// list all agents in use by the GoCD Server
	var a []*gocd.Agent
	var err error
	var r *gocd.APIResponse
	if a, r, err = c.Agents.List(context.Background()); err != nil {
		if r.HTTP.StatusCode == 404 {
			fmt.Println("Couldn't find agent")
		} else {
			panic(err)
		}
	}

	fmt.Println(a)
}

func ExampleConfigRepoService_List() {
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	l, _, err := c.ConfigRepos.List(context.Background())
	if err != nil {
		panic(err)
	}
	// Loops through the list of repositories to display some basic informations
	for _, r := range l {
		fmt.Printf("Pipeline: %s\n\tMaterial type: %s\n", r.ID, r.Material.Type)
		if r.Material.Type == "git" {
			fmt.Printf("\tMaterial url: %s\n", r.Material.Attributes.(*gocd.MaterialAttributesGit).URL)
		}
		fmt.Printf("\tNumber of configuration parameters: %d\n\n", len(r.Configuration))
	}
}

func ExampleConfigRepoService_Get() {
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	r, _, err := c.ConfigRepos.Get(context.Background(), "my_repo_config_id")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Pipeline: %s\n\tMaterial type: %s\n", r.ID, r.Material.Type)
	if r.Material.Type == "git" {
		fmt.Printf("\tMaterial url: %s\n", r.Material.Attributes.(*gocd.MaterialAttributesGit).URL)
	}
	fmt.Printf("\tNumber of configuration parameters: %d\n\n", len(r.Configuration))
}

func ExampleConfigRepoService_Create() {
	// This example creates a config repo that uses a git material and a json config plugin
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	name := "my_pipeline_name"
	repoURL := "git@github.com:example/myrepo.git"
	repo := gocd.ConfigRepo{ID: name, PluginID: "json.config.plugin", Material: gocd.Material{Type: "git", Attributes: &gocd.MaterialAttributesGit{URL: repoURL, Branch: "master", AutoUpdate: true}}}

	r, _, err := c.ConfigRepos.Create(context.Background(), &repo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pipeline created:\n\tName: %s\n\tMaterial type: %s\n", r.ID, r.Material.Type)
	if r.Material.Type == "git" {
		fmt.Printf("\tMaterial url: %s\n", r.Material.Attributes.(*gocd.MaterialAttributesGit).URL)
	}
	fmt.Printf("\tNumber of configuration parameters: %d\n\n", len(r.Configuration))
}

func ExamplePipelinesService_Pause() {
	// This example pauses the pipeline "my_pipeline_name"
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	_, _, err := c.Pipelines.Pause(context.Background(), "my_pipeline_name")
	if err != nil {
		panic(err)
	}
}

func ExamplePipelinesService_Unpause() {
	// This example unpauses the pipeline "my_pipeline_name"
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	_, _, err := c.Pipelines.Unpause(context.Background(), "my_pipeline_name")
	if err != nil {
		panic(err)
	}
}

func ExamplePipelineGroupsService_List() {
	// This example list the pipeline names that belong to the group "foo"
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	groupName := "foo" // If you set your group name to an empty string you will get all the groups
	g, _, err := c.PipelineGroups.List(context.Background(), groupName)
	if err != nil {
		panic(err)
	}
	for _, grp := range *g {
		fmt.Printf("Pipelines in the %s group:\n", grp.Name)
		for _, elt := range grp.Pipelines {
			fmt.Printf("  - %s\n", elt.Name)
		}
	}
}

func ExamplePipelineConfigsService_Delete() {
	// This example deletes the pipeline "my_pipeline_name"
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	_, _, err := c.PipelineConfigs.Delete(context.Background(), "my_pipeline_name")
	if err != nil {
		panic(err)
	}
}

func ExamplePipelineConfigsService_Get() {
	// This example prints out the entire configuration of a pipeline
	cfg := gocd.Configuration{
		Server:   "https://my_gocd/go/", // don't forget the "/go/" at the end of the url to avoid issues!
		Username: "ApiUser",
		Password: "MySecretPassword",
	}

	c := cfg.Client()

	p, _, err := c.PipelineConfigs.Get(context.Background(), "my_pipeline_name")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pipeline configuration:\n")
	fmt.Printf("  - Name: %s\n", p.Name)
	fmt.Printf("  - Group: %s\n", p.Group)
	fmt.Printf("  - Label: %s\n", p.Label)
	fmt.Printf("  - Label template: %s\n", p.LabelTemplate)
	pLocking := "disabled"
	if p.EnablePipelineLocking {
		pLocking = "enabled"
	}
	fmt.Printf("  - Pipeline locking: %s\n", pLocking)
	fmt.Printf("  - Template: %s\n", p.Template)
	if p.Origin != nil {
		fmt.Printf("  - Origin (%s): %s", p.Origin.Type, p.Origin.File)
	}
	fmt.Printf("  - Parameters:\n")
	for _, item := range p.Parameters {
		fmt.Printf("    - %s: %s\n", item.Name, item.Value)
	}
	fmt.Printf("  - Environment variables:\n")
	for _, item := range p.EnvironmentVariables {
		fmt.Printf("    - %s: %s %s (%t)\n", item.Name, item.Value, item.EncryptedValue, item.Secure)
	}
	fmt.Printf("  - Materials:\n")
	for _, item := range p.Materials {
		fmt.Printf("    - Type: %s\n", item.Type)
		fmt.Printf("      Fingerprint: %s\n", item.Fingerprint)
		fmt.Printf("      Description: %s\n", item.Description)
		fmt.Printf("      Attributes:\n")
		m := item.Attributes.GenerateGeneric()
		for k, v := range m {
			fmt.Printf("      - %s: %#v\n", k, v)
		}
	}
	fmt.Printf("  - Stages:\n")
	for _, item := range p.Stages {
		fmt.Printf("    - Name: %s\n", item.Name)
		fmt.Printf("      FetchMaterials: %t\n", item.FetchMaterials)
		fmt.Printf("      CleanWorkingDirectory: %t\n", item.CleanWorkingDirectory)
		fmt.Printf("      NeverCleanupArtifacts: %t\n", item.NeverCleanupArtifacts)
		if item.Approval != nil {
			fmt.Printf("      Approval:\n        Type: %s\n", item.Approval.Type)
			if item.Approval.Authorization != nil {
				fmt.Printf("        Users: %q\n", item.Approval.Authorization.Users)
				fmt.Printf("        Roles: %q\n", item.Approval.Authorization.Roles)
			}
		}
		fmt.Printf("      EnvironmentVariables:\n")
		for _, i := range item.EnvironmentVariables {
			fmt.Printf("        - %s: %s %s (%t)\n", i.Name, i.Value, i.EncryptedValue, i.Secure)
		}
		fmt.Printf("      Resources: %#v\n", item.Resources)
		fmt.Printf("      Jobs:\n")

		for _, i := range item.Jobs {
			fmt.Printf("        - %#v\n", i)
		}
	}
	fmt.Printf("  - Version: %s\n", p.Version)
}
