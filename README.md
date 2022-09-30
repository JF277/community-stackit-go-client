# Community Go Client for STACKIT

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/community-stackit-go-client)](https://goreportcard.com/report/github.com/SchwarzIT/community-stackit-go-client) [![UnitTests](https://github.com/SchwarzIT/community-stackit-go-client/actions/workflows/go.yml/badge.svg)](https://github.com/SchwarzIT/community-stackit-go-client/actions/workflows/go.yml) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/SchwarzIT/community-stackit-go-client) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0)

<br />

This repo's goal is to create a go-based http client for consuming STACKIT APIs

The client is community-supported and not an official STACKIT release, it is maintained by internal Schwarz IT teams integrating with STACKIT


## Usage example



- If you're not sure how to get this information, please contact [STACKIT support](https://support.stackit.cloud)
- To use the Service Account it must be assigned relevant roles using the [Membership API](https://api.stackit.schwarz/membership-service/openapi.v1.html#operation/post-organizations-organizationId-projects-projectId-roles-roleName-service-accounts)
- If your Service Account needs to operate outside the scope of your project, you may need to contact STACKIT to assign further permissions

```
package main

import (
	"context"
	"fmt"
	"os"

	client "github.com/SchwarzIT/community-stackit-go-client"
)

func main() {
	c, err := client.New(context.Background(), &client.Config{
		ServiceAccountID: os.Getenv("STACKIT_SERVICE_ACCOUNT_ID"), [^1]
		Token:            os.Getenv("STACKIT_SERVICE_ACCOUNT_TOKEN"),
		OrganizationID:   os.Getenv("STACKIT_CUSTOMER_ACCOUNT_ID"), [^2]
	})
	if err != nil {
		panic(err)
	}

	projectID := "1234-56789-101112"
	bucketName := "example"

	err = c.ObjectStorage.Buckets.Create(context.TODO(), projectID, bucketName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("bucket '%s' created successfully", bucketName)
}

```

[^1]: In order to use the client, a Service Account and Token needs to be created [using the Service Account API](https://api.stackit.schwarz/service-account/openapi.v1.html#operation/post-projects-projectId-service-accounts-v2)

[^2]: The Customer account ID (or Organization ID) must also be known in advanced. Every project belongs to a Customer Account.