package nuget

import (
	"fmt"
)

func FetchPackageDependencies(packageName, version string) (map[string]string, error) {
	if version == "" {
		versions, err := GetPackageVersions(packageName)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить версии пакета: %v", err)
		}
		version = versions[len(versions)-1]
		fmt.Printf("Версия не указана, используем последнюю: %s\n", version)
	}

	dependenciesList, err := GetNuspec(packageName, version)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить зависимости: %v", err)
	}

	dependencies := make(map[string]string)
	for _, dep := range dependenciesList {
		dependencies[dep.Id] = dep.Version
	}

	return dependencies, nil
}

func PrintDependencies(packageName, version string, dependencies map[string]string) {
	if len(dependencies) == 0 {
		fmt.Printf("Пакет %s %s не имеет прямых зависимостей\n", packageName, version)
		return
	}

	fmt.Printf("Прямые зависимости %s %s:\n", packageName, version)
	for name, ver := range dependencies {
		fmt.Printf("%s %s\n", name, ver)
	}
}
