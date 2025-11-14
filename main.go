package main

import (
	"fmt"
	"graph/config"
	"graph/graph"
	"graph/nuget"
	"graph/testrepo"
	"os"
)

func main() {

	//Загрузка конфигураций
	cfg, err := config.LoadConfig("config_test.yaml") //тут указываем файл yaml с которым будем работать
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигураций: %v", err)
		os.Exit(1)
	}

	config.PrintParam(cfg)

	if err := config.Validate(cfg); err != nil {
		fmt.Printf("\n❌ Ошибка валидации: %v", err)
	}

	fmt.Println("\n\n✅ Конфигурация загружена и проверена успешно!")

	if cfg.Mode {
		fmt.Println("\n\nТЕСТОВЫЙ РЕЖИМ")
		processTestMode(cfg)
	} else {
		processRealMode(cfg)
	}

	fmt.Println("Программа завершенна успешно")
}

func processRealMode(cfg *config.Config) {
	
	//Получение прямых зависимостей
	deps, err := nuget.FetchPackageDependencies(cfg.Package, cfg.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка получения зависимостей: %v\n", err)
		os.Exit(1)
	}

	nuget.PrintDependencies(cfg.Package, cfg.Version, deps)

	//Построение полного графа зависимостей
	getDepsFunc := func (packageName, version string) (map[string]string, error) {
		deps, err := nuget.FetchPackageDependencies(cfg.Package, cfg.Version)
		if err != nil {
			return make(map[string]string), fmt.Errorf("⚠️  Предупреждение: не удалось получить зависимости для %s: %v\n", packageName, err)
		}

		return deps, nil
	}

	dependencyGraph, err := graph.BuildDependencyGraph(cfg.Package, cfg.Version, cfg.Filter, getDepsFunc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка построения графа: %v\n", err)
		os.Exit(1)
	}

	graph.PrintGraph(dependencyGraph, cfg.Package)
}

func processTestMode(cfg *config.Config) {
		// Загружаем тестовый репозиторий из файла
	fmt.Printf("Загружаем тестовый репозиторий из: %s\n", cfg.URL)
	
	repo, err := testrepo.LoadTestRepo(cfg.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка загрузки тестового репозитория: %v\n", err)
		os.Exit(1)
	}

	// Показываем доступные пакеты
	fmt.Printf("Доступные пакеты в репозитории: ")
	for pkg := range repo {
		fmt.Printf("%s ", pkg)
	}
	fmt.Println()

	// Проверяем что запрошенный пакет существует
	if _, exists := repo[cfg.Package]; !exists {
		fmt.Fprintf(os.Stderr, "❌ Пакет '%s' не найден в тестовом репозитории\n", cfg.Package)
		os.Exit(1)
	}

	getDepsFunc := func(packageName, version string) (map[string]string, error) {
		deps := testrepo.GetDepsFromTestRepo(repo)(packageName, version)
		return deps, nil
	}

	deps, err := getDepsFunc(cfg.Package, "") 
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	
	if len(deps) == 0 {
		fmt.Printf("Пакет %s не имеет прямых зависимостей\n", cfg.Package)
	} else {
		fmt.Printf("Прямые зависимости %s:\n", cfg.Package)
		for name := range deps {
			fmt.Printf("   • %s\n", name)
		}
	}

	dependencyGraph, err := graph.BuildDependencyGraph(cfg.Package, "", cfg.Filter, getDepsFunc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка построения графа: %v\n", err)
		os.Exit(1)
	}

	graph.PrintGraph(dependencyGraph, cfg.Package)
}