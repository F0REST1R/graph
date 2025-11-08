package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Package string `yaml:"package"`
	URL     string `yaml:"url"`
	Mode    bool   `yaml:"mode"`
	Version string `yaml:"version"`
	Filter  string `yaml:"filter"`
}

func LoadConfig(filename string) (*Config, error) {
	//Проверяем существ.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("Файл конфигурации не найден: %s", filename)
	}

	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении файла: %s", filename)
	}

	// Парсинг YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("Ошибка парсинга YAML: %v", err)
	}

	return  &cfg, nil
}

func PrintParam(cfg *Config) {
	fmt.Println("Параметры:")
	fmt.Printf("Пакет: %s\n", cfg.Package)
	fmt.Printf("Репозиторий: %s\n", cfg.URL)
	fmt.Printf("Режим работы: %t\n", cfg.Mode)
	fmt.Printf("Версия: %s\n", cfg.Version)
	fmt.Printf("Фильтрация: %s\n", cfg.Filter)
}

func Validate(cfg *Config) (error) {
	errors := make([]string, 0, 5)

	//Проверка именни пакета
	if strings.TrimSpace(cfg.Package) == "" {
		errors = append(errors, "Имя файла не может быть пустым")
	}

	if strings.TrimSpace(cfg.URL) == "" {
		errors = append(errors, "Сыллка не может быть пустой")
	}

	if cfg.Mode{
		if _, err := os.Stat(cfg.URL); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("Файл тестового репозитория не существует: %s", cfg.URL))
		}
	}

	if strings.TrimSpace(cfg.Version) != cfg.Version{
		errors = append(errors, "Версия не должна содержать только пробельные символы")
	}

	if len(errors) > 0 {
		return fmt.Errorf("Ошибки валидации конфигурации: %s", strings.Join(errors, ";"))
	}

	return nil
}


