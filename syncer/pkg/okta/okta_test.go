package okta_test

import (
	"fmt"
	"testing"

	ookta "github.com/okta/okta-sdk-golang/v2/okta"
	. "github.com/onsi/gomega"
	"github.com/zendesk/okta-identity-governance-demo/pkg/okta"
	"github.com/zendesk/okta-identity-governance-demo/pkg/yaml"
)

func TestTransformOktaUsers(t *testing.T) {
	oktaUsers := []*ookta.User{
		{
			Id: "00u179bsnjqfXS6dH697",
			Profile: &ookta.UserProfile{
				"email":    "test@example.com",
				"test-app": "admin",
			},
		},
		{
			Id: "00u190bsnjqfXS6dH697",
			Profile: &ookta.UserProfile{
				"email":      "test2@example.com",
				"second-app": "role1,role2",
			},
		},
	}
	testUsers := []yaml.User{
		{
			Email:  "test@example.com",
			OktaID: "00u179bsnjqfXS6dH697",
			Attributes: map[string][]string{
				"test-app":   {"admin"},
				"second-app": {""},
			},
		},
		{
			Email:  "test2@example.com",
			OktaID: "00u190bsnjqfXS6dH697",
			Attributes: map[string][]string{
				"test-app":   {""},
				"second-app": {"role1", "role2"},
			},
		},
	}

	testAttribute := []yaml.Attribute{
		{
			AttributeName: "test-app",
		},
		{
			AttributeName: "second-app",
		},
	}

	t.Parallel()
	inputs := []struct {
		name        string
		oktaUsers   []*ookta.User
		returnUsers []yaml.User
		hasError    bool
	}{
		{
			"no changes",
			oktaUsers,
			testUsers,
			false,
		},
	}

	for _, test := range inputs {
		test := test
		t.Run(fmt.Sprintf("when %s", test.name), func(tt *testing.T) {
			tt.Parallel()
			g := NewGomegaWithT(tt)

			actual := okta.TransformOktaUsers(oktaUsers, testAttribute)

			g.Expect(actual).To(Equal(test.returnUsers), fmt.Sprintf("when %s", test.name))
		})
	}
}
