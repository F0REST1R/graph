package graph

import (
	"fmt"
	"sort"
	"strings"
)

func GeneratePlantUML(graph map[string][]string, title string) string {
	var builder strings.Builder

	builder.WriteString("@startuml\n")
	builder.WriteString("title Граф зависимостей: " + title + "\n")
	builder.WriteString("skinparam componentStyle rectangle\n")
	builder.WriteString("skinparam nodesep 20\n")
	builder.WriteString("skinparam ranksep 30\n\n")

	nodes := make([]string, 0, len(graph))
	for node := range graph {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	for _, node := range nodes {
		safeNode := strings.ReplaceAll(node, " ", "\\n")
		builder.WriteString(fmt.Sprintf("component \"%s\" as %s\n", 
			safeNode, 
			generateNodeID(node)))
	}
	
	builder.WriteString("\n")

	for _, node := range nodes {
		deps := graph[node]
		sort.Strings(deps)
		
		for _, dep := range deps {
			// Проверяем что зависимость существует в графе
			if _, exists := graph[dep]; exists {
				builder.WriteString(fmt.Sprintf("%s --> %s\n",
					generateNodeID(node),
					generateNodeID(dep)))
			}
		}
	}

	builder.WriteString("@enduml")
	return builder.String()
}

//создает уникальный идентификатор для узла в PlantUML
func generateNodeID(node string) string {
	// Заменяем пробелы и специальные символы на подчеркивания
	id := strings.ReplaceAll(node, " ", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, "@", "_")
	return "Node_" + id
}

func PrintPlantUML(plantUMLCode string, title string) {
	fmt.Printf("\n🎨 PlantUML ДИАГРАММА: %s\n", title)
	fmt.Println("═" + strings.Repeat("═", 50))
	fmt.Println(plantUMLCode)
	fmt.Println("═" + strings.Repeat("═", 50))
	
	fmt.Println("\n📋 ИНСТРУКЦИЯ ДЛЯ ВИЗУАЛИЗАЦИИ:")
	fmt.Println("1. Скопируйте код выше")
	fmt.Println("2. Перейдите на сайт: https://www.plantuml.com/plantuml/")
	fmt.Println("3. Вставьте код в текстовое поле")
	fmt.Println("4. Нажмите 'Submit' для генерации диаграммы")
	fmt.Println("5. Или используйте локальную установку PlantUML")
}