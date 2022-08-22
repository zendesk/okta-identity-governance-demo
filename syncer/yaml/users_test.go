package yaml_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/zendesk/okta-identity-governance-demo/yaml"
)

func TestCompareUsers(t *testing.T) {
	testUser1 := []yaml.User{
		{
			Email: "test@example.com",
			Attributes: map[string][]string{
				"test-app": {"admin"},
			},
		},
	}

	testUser2 := []yaml.User{
		{
			Email: "test2@example.com",
			Attributes: map[string][]string{
				"test-app": {"admin", "read-only"},
			},
		},
	}

	testAttribute := []yaml.Attribute{
		{
			AttributeName: "test-app",
		},
	}

	t.Parallel()
	inputs := []struct {
		name        string
		oktaUsers   []yaml.User
		yamlUsers   []yaml.User
		returnUsers []yaml.User
		hasError    bool
	}{
		{
			"no changes",
			testUser1,
			testUser1,
			[]yaml.User{},
			false,
		},
		{
			"adding attribute",
			[]yaml.User{
				{
					Email:      "test@example.com",
					Attributes: map[string][]string{},
				},
			},
			testUser1,
			testUser1,
			false,
		},
		{
			"removing attribute",
			testUser1,
			[]yaml.User{
				{
					Email:      "test@example.com",
					Attributes: map[string][]string{},
				},
			},
			[]yaml.User{
				{
					Email:      "test@example.com",
					Attributes: map[string][]string{},
				},
			},
			false,
		},
		{
			"missing Okta user and update",
			append(testUser1, testUser2...),
			[]yaml.User{
				{
					Email: "test2@example.com",
					Attributes: map[string][]string{
						"test-app": {"read-only"},
					},
				},
			},
			[]yaml.User{
				{
					Email: "test2@example.com",
					Attributes: map[string][]string{
						"test-app": {"read-only"},
					},
				},
				{
					Email: "test@example.com",
					Attributes: map[string][]string{
						"test-app": {""},
					},
				},
			},
			false,
		},
	}

	for _, test := range inputs {
		test := test
		t.Run(fmt.Sprintf("when %s", test.name), func(tt *testing.T) {
			tt.Parallel()
			g := NewGomegaWithT(tt)

			actual, err := yaml.CompareUsers(test.yamlUsers, test.oktaUsers, testAttribute)
			if test.hasError {
				g.Expect(err).To(HaveOccurred(), fmt.Sprintf("when %s", test.name))
			} else {
				g.Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("when %s", test.name))
			}

			g.Expect(actual).To(Equal(test.returnUsers), fmt.Sprintf("when %s", test.name))
		})
	}
}
