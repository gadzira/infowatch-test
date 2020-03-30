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

func filePathWalkDir(root string, c chan string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			c <- path
		}
		return nil
	})
	close(c)
	return err
}

func readFile(c chan string) {
	var lines []string
	for i := range c {
		file, err := os.Open(i)
		fmt.Println(i)
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

	fileChan := make(chan string, 100)

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
	err := filePathWalkDir(root, fileChan)
	if err != nil {
		panic(err)
	}

	readFile(fileChan)
}
