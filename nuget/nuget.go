package nuget

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	URLAPI = "https://api.nuget.org/v3-flatcontainer"
	timeout      = 30 * time.Second
)

type Dependency struct {
	Id      string
	Version string
}

func fetchURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Сетевой сбой: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d для URL: %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения ответа: %v", err)
	}

	return body, nil
}

func parseNuspecManual(xmlContent string) ([]Dependency, error) {
	var dependencies []Dependency

	content := strings.ToLower(xmlContent)

	depsStart := strings.Index(content, "<dependencies>")
	if depsStart == -1 {
		return dependencies, nil
	}

	depsEnd := strings.Index(content, "</dependencies>")
	if depsEnd == -1 {
		return nil, fmt.Errorf("Некоректный формат .nuspec: не найдено закрытие dependencies")
	}

	depsSection := content[depsStart:depsEnd]

	lines := strings.Split(depsSection, "<dependency")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		idStart := strings.Index(line, "id=\"")
		if idStart == -1 {
			continue
		}

		idStart += 4 
		idEnd := strings.Index(line[idStart:], "\"")
		if idEnd == -1 {
			continue
		}
		id := line[idStart:idStart+idEnd]

		versionStart := strings.Index(line, "version=\"")
		if versionStart == -1 {
			continue
		}
		versionStart += 9 // длина "version=\""
		versionEnd := strings.Index(line[versionStart:], "\"")
		if versionEnd == -1 {
			continue
		}
		version := line[versionStart : versionStart+versionEnd]
		
		// Восстанавливаем оригинальный регистр для id
		originalId := extractOriginalId(xmlContent, id)
		
		dependencies = append(dependencies, Dependency{
			Id:      originalId,
			Version: version,
		})
	}

	return dependencies, nil
}

func extractOriginalId(originalXml, lowerId string) string {
	// Ищем точное соответствие в оригинальном XML
	idStart := strings.Index(originalXml, "id=\""+lowerId+"\"")
	if idStart != -1 {
		// Пытаемся найти регистронезависимое соответствие
		lowerXml := strings.ToLower(originalXml)
		idStart = strings.Index(lowerXml, "id=\""+lowerId+"\"")
		if idStart != -1 {
			idStart += 4 // длина "id=\""
			idEnd := strings.Index(originalXml[idStart:], "\"")
			if idEnd != -1 {
				return originalXml[idStart : idStart+idEnd]
			}
		}
	}
	return lowerId // fallback
}

func GetNuspec(packageName, version string) ([]Dependency, error) {
	normalizedName := strings.ToLower(packageName)
	url := fmt.Sprintf("%s/%s/%s/%s.nuspec", URLAPI, normalizedName, version, normalizedName)
	
	fmt.Printf("Загружаем зависимости из: %s\n", url)
	
	data, err := fetchURL(url)
	if err != nil {
		return nil, err
	}

	dependencies, err := parseNuspecManual(string(data))
	if err != nil {
		return nil, err
	}

	return dependencies, nil
}

func GetPackageVersions(packageName string) ([]string, error) {
	normalizedName := strings.ToLower(packageName)
	url := fmt.Sprintf("%s/%s/index.json", URLAPI, normalizedName)
	
	data, err := fetchURL(url)
	if err != nil {
		return nil, err
	}

	// Упрощенный парсинг JSON для получения версий
	content := string(data)
	var versions []string
	
	// Ищем строки в кавычках, которые похожи на версии
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "\"") && strings.Count(trimmed, "\"") >= 2 {
			parts := strings.Split(trimmed, "\"")
			if len(parts) >= 3 {
				candidate := parts[1]
				// Проверяем что это версия (содержит цифры и точки)
				if isVersionString(candidate) {
					versions = append(versions, candidate)
				}
			}
		}
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("Не удалось найти версию пакета")
	}

	return versions, nil
}

func isVersionString(s string) bool {
	if s == "" {
		return false
	}
	hasDigit := false
	for _, char := range s {
		if char >= '0' && char <= '9' {
			hasDigit = true
		} else if char != '.' && char != '-' && char != '+' {
			return false
		}
	}
	return hasDigit
}