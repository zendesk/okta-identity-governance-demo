package yaml

func GenerateUsers(attributes []Attribute, teams []Team) ([]User, error) {
	// Create map of all YAML emails users
	userMap := map[string]User{}
	for _, team := range teams {
		for _, member := range team.Members {
			userMap[member] = User{Email: member, Attributes: make(map[string][]string)}
		}
	}

	// Create empty Attribute values for each user
	for email, user := range userMap {
		for _, attribute := range attributes {
			userAttributes := user.Attributes
			userAttributes[attribute.AttributeName] = []string{""}
			user.Attributes = userAttributes
			userMap[email] = user
		}
	}

	// Loop through attributes to assign values
	for _, attribute := range attributes {
		for attributeTeam, attributeValues := range attribute.AttributeMap {
			for _, team := range teams {
				if attributeTeam == team.TeamName {
					userMap = addValuesToUsers(userMap, team.Members, attribute.AttributeName, attributeValues)
				}
			}
		}
	}

	// Return slice off map
	users := []User{}

	for _, user := range userMap {
		users = append(users, user)
	}

	return users, nil
}

func addValuesToUsers(userMap map[string]User, members []string, attribute string, values []string) map[string]User {
	for _, member := range members {
		if user, found := userMap[member]; found {
			if len(user.Attributes[attribute]) == 1 && user.Attributes[attribute][0] == "" {
				user.Attributes[attribute] = values
			} else {
				for _, value := range values {
					if !contains(user.Attributes[attribute], value) {
						user.Attributes[attribute] = append(user.Attributes[attribute], value)
					}
				}
			}
			userMap[member] = user
		}
	}
	return userMap
}
