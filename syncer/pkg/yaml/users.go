package yaml

import "log"

// User represents the profile attributes and identifiers for an Okta user
type User struct {
	Email      string
	OktaID     string
	Attributes map[string][]string
}

// CompareUsers compares YAML and Okta Users and returns a list of Users requiring an update
func CompareUsers(yamlUsers, oktaUsers []User, attributes []Attribute) ([]User, error) {
	changeSet := []User{}

	// Find existing YAML users needing an update
	for _, yamlUser := range yamlUsers {
		found := false
		for _, oktaUser := range oktaUsers {
			if yamlUser.Email == oktaUser.Email {
				found = true
				if !mapsEqual(yamlUser.Attributes, oktaUser.Attributes) {
					log.Printf("Updating user in Okta: %v \t\n Current: %v\t\n Desired: %v\n\n", yamlUser.Email, oktaUser.Attributes, yamlUser.Attributes)
					yamlUser.OktaID = oktaUser.OktaID
					changeSet = append(changeSet, yamlUser)
				}
			}
		}
		if !found {
			log.Printf("User not found in Okta: %v\n", yamlUser.Email)
		}
	}

	// Find users missing from YAML
	for _, oktaUser := range oktaUsers {
		if len(oktaUser.Attributes) == 0 {
			continue
		}
		foundInYaml := false
		for _, yamlUser := range yamlUsers {
			if yamlUser.Email == oktaUser.Email {
				foundInYaml = true
			}
		}

		// Add missing Okta User to changeset with no attributes
		if !foundInYaml {
			profileEmpty := true
			for _, profileAttribute := range oktaUser.Attributes {
				if (len(profileAttribute) > 1) || (profileAttribute[0] != "") {
					profileEmpty = false
				}
			}
			if !profileEmpty {
				for _, attribute := range attributes {
					oktaUser.Attributes[attribute.AttributeName] = []string{""}
				}
				log.Printf("Clearing out permissions from Okta user missing from YAML: %v\n", oktaUser.Email)
				changeSet = append(changeSet, oktaUser)
			}
		}
	}

	return changeSet, nil
}
