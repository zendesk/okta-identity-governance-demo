package yaml

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Attribute struct {
	Type          string              `yaml:"type"`
	AttributeName string              `yaml:"attribute_name"`
	AttributeMap  map[string][]string `yaml:"attribute_map"`
}

type Team struct {
	Type     string   `yaml:"type"`
	TeamName string   `yaml:"team_name"`
	Members  []string `yaml:"members"`
}

func getYamlFiles() ([]string, error) {
	files := []string{}
	yamlPaths := []string{"../attributes/", "../teams/"}
	for _, y := range yamlPaths {
		err := filepath.Walk(y,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.Contains(path, ".yaml") {
					files = append(files, path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return files, nil
}

func importYamlFile(path string) (map[interface{}]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(b, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}
	return m, err
}

func GetYaml() ([]Attribute, []User, error) {
	attributes := []Attribute{}
	teams := []Team{}

	files, err := getYamlFiles()
	if err != nil {
		return nil, nil, err
	}

	for _, file := range files {
		m, err := importYamlFile(file)
		if err != nil {
			return nil, nil, err
		}
		if m["type"] == "attribute" {
			attr := Attribute{}
			b, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, nil, err
			}
			err = yaml.Unmarshal(b, &attr)
			if err != nil {
				return nil, nil, err
			}
			attributes = append(attributes, attr)
		}
		if m["type"] == "team" {
			team := Team{}
			b, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, nil, err
			}
			err = yaml.Unmarshal(b, &team)
			if err != nil {
				return nil, nil, err
			}
			teams = append(teams, team)
		}
	}

	users, err := GenerateUsers(attributes, teams)
	if err != nil {
		return nil, nil, err
	}

	return attributes, users, nil
}
