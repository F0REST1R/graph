package graph

import (
	"fmt"
	"strings"
)

func BuildDependencyGraph(packageName, version, filter string, getDepsFunc func(string, string) (map[string]string, error) ) (map[string][]string, error) {
	graph := make(map[string][]string)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(string, string) error
	dfs = func(name, ver string) error {
		node := name
		if ver != "" {
			node = fmt.Sprintf("%s %s", name, ver)
		}

		if recStack[node] {
			fmt.Printf("Обнаружен стек зависимостей для узла: %s\n", node)
			fmt.Printf("Цикл прерван (Избегание бесконечной рекурсии)\n")
			return nil
		}

		if filter != "" && name != packageName && strings.Contains(strings.ToLower(name), strings.ToLower(filter)){
			fmt.Printf("Пропускаем пакет (фильтр '%s'): %s\n", filter, node)
			visited[node] = true
			return nil
		}

		if visited[node] {
			return nil
		}

		visited[node] = true
		recStack[node] = true

		deps, err := getDepsFunc(name, ver)
		if err != nil{
			return  fmt.Errorf("Ошибка получения зависимостей для %s: %v", node, err)
		}

		var children []string
		for depName, depVer := range deps {
			if filter != "" && strings.Contains(strings.ToLower(depName), strings.ToLower(filter)) {
				fmt.Printf("Пропускаем пакет (фильтр '%s'): %s\n", filter, depName)
				continue
			}
			
			childNode := depName
			if depVer != "" {
				childNode = fmt.Sprintf("%s %s", depName, depVer)
			}

			children = append(children, childNode)

			if err := dfs(depName, depVer); err != nil {
				return err
			}
		}
		
		graph[node] = children

		recStack[node] = false

		return nil
	}

	if err := dfs(packageName, version); err != nil {
		return nil, err
	}

	return graph, nil
}

func PrintGraph(graph map[string][]string, root string) {
	fmt.Printf("\n📊 Полный граф зависимостей (корень: %s):\n", root)
	for node, deps := range graph {
		if len(deps) == 0 {
			fmt.Printf("%s -> (нет зависимостей)\n", node)
		} else {
			fmt.Printf("%s -> %v\n", node, deps)
		}
	}
}