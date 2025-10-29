package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Запрос пути к JSON-файлу
	fmt.Print("Введите путь к JSON-файлу: ")
	var filePath string
	fmt.Scanln(&filePath)

	// Чтение файла
	data, err := os.ReadFile(filePath)
	//data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		return
	}

	// Декодирование JSON в интерфейс
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Printf("Ошибка парсинга JSON: %v\n", err)
		return
	}

	fmt.Println("JSON успешно загружен!")
	fmt.Println("Введите ключи для поиска (через точку для вложенных полей)")
	fmt.Println("Для выхода введите 'exit'")

	// Цикл обработки пользовательского ввода
	for {
		fmt.Print("\nВведите ключ: ")
		var input string
		fmt.Scanln(&input)

		input = strings.TrimSpace(input)
		if input == "exit" {
			fmt.Println("Выход из программы")
			break
		}

		if input == "" {
			continue
		}

		// Поиск значения по ключу
		value, found := findValue(jsonData, input)
		if found {
			fmt.Printf("Значение для ключа '%s': %v\n", input, value)
		} else {
			fmt.Printf("Ключ '%s' не найден\n", input)
		}
	}
}

// findValue ищет значение по ключу в структуре JSON
// Поддерживает вложенные ключи через точку (например: "user.address.city")
func findValue(data interface{}, key string) (interface{}, bool) {
	// Разделяем ключ на части если есть вложенность
	keys := strings.Split(key, ".")

	return findNestedValue(data, keys)
}

// findNestedValue рекурсивно ищет значение по пути ключей
func findNestedValue(data interface{}, keys []string) (interface{}, bool) {
	if len(keys) == 0 {
		return data, true
	}

	currentKey := keys[0]
	remainingKeys := keys[1:]

	switch v := data.(type) {
	case map[string]interface{}:
		// Если это объект - ищем ключ
		if value, exists := v[currentKey]; exists {
			if len(remainingKeys) == 0 {
				return value, true
			}
			return findNestedValue(value, remainingKeys)
		}
	case []interface{}:
		// Если это массив - пытаемся преобразовать ключ в индекс
		if index, err := stringToIndex(currentKey); err == nil && index >= 0 && index < len(v) {
			if len(remainingKeys) == 0 {
				return v[index], true
			}
			return findNestedValue(v[index], remainingKeys)
		}
	}

	return nil, false
}

// stringToIndex пытается преобразовать строку в индекс массива
func stringToIndex(s string) (int, error) {
	var index int
	_, err := fmt.Sscanf(s, "%d", &index)
	if err != nil {
		return -1, err
	}
	return index, nil
}

// Вспомогательная функция для красивого вывода (опционально)
func prettyPrint(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case nil:
		return "null"
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		// JSON числа всегда float64 в Go
		return fmt.Sprintf("%g", v)
	default:
		// Для сложных структур форматируем как JSON
		jsonBytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	}
}
