package main

// 1. Собрать список файлов +
// 2. Открыть файл +
// 3. Создать мапу куда заносить новый символ или обновлять значение существующего
// 4. После закрытия файла отрисовать гистогрмму

// get full path to file -> send result to channel ->
// read file line by line and send to channel ->
// scan every line and add symbol to map as key and iterate value

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func readFile(path string) {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	for _, eachline := range lines {
		fmt.Println(eachline)
	}
}

// goal: add every symbol to map
// r := bufio.NewReader(strings.NewReader(text))
// for {
// if c, _, err := r.ReadRune(); err != nil {
// if err == io.EOF {
// break
// } else {
// log.Fatal(err)
// }
// } else {
// fmt.Println(string(c))
// }
// }

func main() {

	fileCh := make(chan string, 100)
	lineCh := make(chan string, 100)

	var root string
	if len(os.Args) == 1 {
		log.Fatal("No path given, Please specify path.")
		return
	}

	if root = os.Args[1]; root == "" {
		log.Fatal("No path given, Please specify path.")
		return
	}

	// filepath.Walk
	files, err := FilePathWalkDir(root)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		readFile(file)
	}
}
