package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// WordCount представляет пару "слово-количество"
type WordCount struct {
	Word  string
	Count int
}

func main() {
	// Парсинг аргументов командной строки
	filePath := flag.String("file", "", "Путь к текстовому файлу")
	flag.Parse()

	// Проверка обязательного параметра
	if *filePath == "" {
		fmt.Println("Ошибка: необходимо указать путь к файлу с помощью флага -file")
		flag.Usage()
		os.Exit(1)
	}

	// Чтение файла
	content, err := readFile(*filePath)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		os.Exit(1)
	}

	// Подсчет слов
	wordCounts := countWords(content)

	// Получение топ-10 слов
	topWords := getTopWords(wordCounts, 10)

	// Вывод результатов
	fmt.Printf("Топ-10 самых частых слов в файле '%s':\n", *filePath)
	fmt.Println("==========================================")
	for i, wc := range topWords {
		fmt.Printf("%d. %s: %d\n", i+1, wc.Word, wc.Count)
	}
}

// readFile читает содержимое файла
func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + " ")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

// countWords подсчитывает количество уникальных слов в тексте
func countWords(text string) map[string]int {
	// Приводим к нижнему регистру и удаляем знаки препинания
	reg := regexp.MustCompile(`[^\p{L}0-9_\s]`)
	cleanText := reg.ReplaceAllString(strings.ToLower(text), "")

	// Разбиваем на слова
	words := strings.Fields(cleanText)

	// Подсчитываем слова
	wordCounts := make(map[string]int)
	for _, word := range words {
		// Пропускаем пустые строки
		if word == "" {
			continue
		}
		wordCounts[word]++
	}

	return wordCounts
}

// getTopWords возвращает топ N самых частых слов
func getTopWords(wordCounts map[string]int, n int) []WordCount {
	// Создаем слайс для сортировки
	var wordList []WordCount
	for word, count := range wordCounts {
		wordList = append(wordList, WordCount{Word: word, Count: count})
	}

	// Сортируем по убыванию количества
	sort.Slice(wordList, func(i, j int) bool {
		if wordList[i].Count == wordList[j].Count {
			return wordList[i].Word < wordList[j].Word
		}
		return wordList[i].Count > wordList[j].Count
	})

	// Возвращаем топ N
	if len(wordList) > n {
		return wordList[:n]
	}
	return wordList
}
