package graph

import (
	"fmt"
	"strings"
)

func BuildDependencyGraph(packageName, version, filter string, getDepsFunc func(string, string) map[string]string) (map[string][]string, error) {
	graph := make(map[string][]string)
	visited := make(map[string]bool)

	var dfs func(string, string) error
	dfs = func(name, ver string) error {
		node := name
		if ver != "" {
			node = fmt.Sprintf("%s %s", name, ver)
		}

		if visited[node] {
			return nil
		}

		if filter != "" && name != packageName && strings.Contains(strings.ToLower(name), strings.ToLower(filter)) {
			visited[node] = true
			return nil
		}

		visited[node] = true

		deps := getDepsFunc(name, ver)

		var children []string
		for depName, depVer := range deps {
			if filter != "" && strings.Contains(strings.ToLower(depName), strings.ToLower(filter)) {
				continue
			}

			childNode := depName
			if depVer != "" {
				childNode = fmt.Sprintf("%s %s", depName, depVer)
			}

			if !visited[childNode] {
				if err := dfs(depName, depVer); err != nil {
					return err
				}
			}

			if _, exists := graph[childNode]; exists {
				children = append(children, childNode)
			}
		}

		graph[node] = children
		return nil
	}

	if err := dfs(packageName, version); err != nil {
		return nil, err
	}

	return graph, nil
}

func PrintGraph(graph map[string][]string, root string) {
	fmt.Printf("\nПолный граф зависимостей (корень: %s):\n", root)
	for node, deps := range graph {
		if len(deps) == 0 {
			fmt.Printf("%s -> (нет зависимостей)\n", node)
		} else {
			fmt.Printf("%s -> %v\n", node, deps)
		}
	}
}