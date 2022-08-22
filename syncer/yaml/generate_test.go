package yaml_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/zendesk/okta-identity-governance-demo/yaml"
)

func TestGenerateUsers(t *testing.T) {
	testUsers := []yaml.User{
		{
			Email: "test@example.com",
			Attributes: map[string][]string{
				"test-app": {"admin"},
			},
		},
	}

	testTeams := []yaml.Team{
		{
			TeamName: "test-team",
			Members: []string{
				"test@example.com",
			},
		},
	}

	testAttribute := []yaml.Attribute{
		{
			AttributeName: "test-app",
			AttributeMap:  map[string][]string{"test-team": {"admin"}},
		},
	}

	t.Parallel()
	inputs := []struct {
		name        string
		attributes  []yaml.Attribute
		teams       []yaml.Team
		returnUsers []yaml.User
		hasError    bool
	}{
		{
			"no changes",
			testAttribute,
			testTeams,
			testUsers,
			false,
		},
	}

	for _, test := range inputs {
		test := test
		t.Run(fmt.Sprintf("when %s", test.name), func(tt *testing.T) {
			tt.Parallel()
			g := NewGomegaWithT(tt)

			actual, err := yaml.GenerateUsers(test.attributes, test.teams)
			if test.hasError {
				g.Expect(err).To(HaveOccurred(), fmt.Sprintf("when %s", test.name))
			} else {
				g.Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("when %s", test.name))
			}

			g.Expect(actual).To(Equal(test.returnUsers), fmt.Sprintf("when %s", test.name))
		})
	}
}
