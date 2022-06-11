package main

import (
	"log"
	"os"

	"github.com/zendesk/okta-identity-governance-demo/pkg/okta"
	"github.com/zendesk/okta-identity-governance-demo/pkg/yaml"
)

func main() {
	// Collect the users and attributes from local YAML files
	attributes, yamlUsers, err := yaml.GetYaml()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	client, err := okta.NewOktaClient()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Collect the profiles from Okta and transform into objects
	oktaUsers, err := client.GetOktaUsers(attributes)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Find the difference between attributes in YAML vs Okta
	changeset, err := yaml.CompareUsers(yamlUsers, oktaUsers, attributes)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Only update users that need it
	for _, user := range changeset {
		err := client.UpdateUser(user, attributes)
		if err != nil {
			log.Printf("Error updating users: %v\n", err)
		}
	}

	log.Printf("\nOkta Sync successful for %v users\n\n", len(changeset))
}
