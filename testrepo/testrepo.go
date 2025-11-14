package testrepo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadTestRepo(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	repo := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#"){
			continue
		}

		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		pkg := strings.TrimSpace(parts[0])
		depsStr := strings.TrimSpace(parts[1])

		var deps []string
		if depsStr != "" {
			deps = strings.Fields(depsStr)
		}

		repo[pkg] = deps
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла: %v", err)
	}

	return repo, nil
}

func GetDepsFromTestRepo(repo map[string][]string) func(string, string) map[string]string {
	return func(packageName, version string) map[string]string {
		depsList, exists := repo[packageName]
		if !exists {
			return make(map[string]string)
		}

		deps := make(map[string]string)
		for _, dep := range depsList {
			deps[dep] = "" 
		}

		return deps
	}
}