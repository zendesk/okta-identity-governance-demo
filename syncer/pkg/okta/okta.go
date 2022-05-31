package okta

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zendesk/okta-identity-governance-demo/pkg/yaml"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type OktaClient struct {
	Client   *okta.Client
	Ctx      context.Context
	OrgName  string
	OrgUrl   string
	ApiToken string
}

func NewOktaClient() (OktaClient, error) {
	org := os.Getenv("OKTA_ORG_NAME")
	baseUrl := os.Getenv("OKTA_BASE_URL")
	apiToken := os.Getenv("OKTA_API_TOKEN")
	ctx, cli, err := okta.NewClient(
		context.TODO(),
		okta.WithOrgUrl("https://"+org+"."+baseUrl),
		okta.WithToken(apiToken),
	)
	client := OktaClient{Client: cli, Ctx: ctx}

	return client, err
}

func (c *OktaClient) GetOktaUsers(attributes []yaml.Attribute) ([]yaml.User, error) {

	users, _, err := c.Client.User.ListUsers(c.Ctx, nil)
	if err != nil {
		return nil, err
	}

	transformedUsers := TransformOktaUsers(users, attributes)
	return transformedUsers, err
}

func (c *OktaClient) UpdateUser(user yaml.User, attributes []yaml.Attribute) error {
	userToUpdate, _, err := c.Client.User.GetUser(c.Ctx, user.OktaId)
	if err != nil {
		return err
	}

	newProfile := *userToUpdate.Profile
	for _, attribute := range attributes {

		newProfile[attribute.AttributeName] = strings.Join(user.Attributes[attribute.AttributeName], ",")
	}
	updateUser := &okta.User{
		Profile: &newProfile,
	}

	_, _, err = c.Client.User.UpdateUser(c.Ctx, userToUpdate.Id, *updateUser, nil)
	if err != nil {
		return err
	}

	return nil
}

func TransformOktaUsers(users []*okta.User, attributes []yaml.Attribute) []yaml.User {
	transformedUsers := []yaml.User{}
	for _, user := range users {
		profile := *user.Profile
		updatedAttributes := map[string][]string{}
		for _, attribute := range attributes {
			if value, found := profile[attribute.AttributeName]; found {
				if value.(string) == "" {
					updatedAttributes[attribute.AttributeName] = []string{""}
					continue
				}
				v := fmt.Sprintf("%v", value)
				cleaned := strings.Replace(v, ",", " ", -1)

				updatedAttributes[attribute.AttributeName] = strings.Fields(cleaned)

			} else {
				updatedAttributes[attribute.AttributeName] = []string{""}
			}

		}
		transformedUsers = append(transformedUsers, yaml.User{Email: profile["email"].(string), OktaId: user.Id, Attributes: updatedAttributes})
	}
	return transformedUsers
}
