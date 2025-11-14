package graph

import "fmt"

func GetLoadOrder(graph map[string][]string, root string) []string {
	visited := make(map[string]bool)
	var order []string

	var dfs func(string)
	dfs = func(node string) {
		if visited[node] {
			return
		}

		visited[node] = true

		for _, child := range graph[node] {
			dfs(child)
		}

		order = append(order, node)
	}

	dfs(root)

	return order
}

// выводит порядок загрузки в читаемом формате
func PrintLoadOrder(order []string, root string) {
	fmt.Printf("\n----- ПОРЯДОК ЗАГРУЗКИ ЗАВИСИМОСТЕЙ (корень: %s):\n", root)
	for i, pkg := range order {
		fmt.Printf("%2d. %s\n", i+1, pkg)
	}
}

