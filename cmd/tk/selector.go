package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/tanka"
)

func environmentSelector(path string, selector labels.Selector) (string, string, error) {
	paths := make(map[string]string, 0)

	envs, err := tanka.FindEnvironments(path, selector)
	if err != nil {
		return "", "", err
	}

	if len(envs) == 0 {
		return "", "", fmt.Errorf("no environment found")
	}

	for path, envs := range envs {
		for _, env := range envs {
			if env != nil {
				paths[env.Metadata.Name] = path
			}
		}
	}

	names := []string{}
	for name, _ := range paths {
		names = append(names, name)
	}
	sort.Strings(names)

	prompt := promptui.Select{
		Label: "Select Environment",
		Items: names,
		Size:  10,
		Searcher: func(input string, index int) bool {
			return strings.Contains(names[index], input)
		},
	}

	_, selected, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("prompt failed %v\n", err)
	}
	return selected, paths[selected], nil
}
