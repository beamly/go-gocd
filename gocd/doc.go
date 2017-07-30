// Copyright 2017 Drew J. Sonne. All rights reserved.
//
// Use of this source code is governed by a LGPLv3-style
// license that can be found in the LICENSE file.

/*
Package gocd provides a client for using the GoCD Server API.

Usage:

	import "github.com/drewsonne/go-gocd/gocd"

Construct a new GoCD client, then use the various services on the client to
access different parts of the GoCD Server API. For example:

	client := gocd.NewClient("https://goserver:8154/go, &gocd.Auth{
		Username: os.GetEnv("GOCD_USERNAME"),
		Password: os.GetEnv("GOCD_PASSWORD"),
	}, nil, false)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Agents.List(ctx)

Some API methods have optional parameters that can be passed. For example:

	client := github.NewClient(nil)

The services of a client divide the API into logical chunks and correspond to
the structure of the GoCD API documentation at
https://api.gocd.org/17.7.0/.

*/

package gocd
