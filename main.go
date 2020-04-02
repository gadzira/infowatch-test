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
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func readFile(c chan string, l chan string, waitgroup *sync.WaitGroup) {
	var lines []string
	for i := range c {
		file, err := os.Open(i)
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
			l <- eachline
		}
	}
	close(l)
	waitgroup.Done()
}

func sortChars(l chan string, waitgroup *sync.WaitGroup) map[string]int {
	sourceMap := make(map[string]int)
	for i := range l {
		r := bufio.NewReader(strings.NewReader(i))
		for {
			if c, _, err := r.ReadRune(); err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			} else {
				sc := string(c)
				if _, found := sourceMap[sc]; found {
					sourceMap[sc]++
				} else {
					sourceMap[sc] = 1
				}
			}
		}
	}
	// for test
	for k, v := range sourceMap {
		fmt.Printf("[%s] : %d\n", k, v)
	}
	waitgroup.Done()
	return sourceMap
}

func main() {

	fileChan := make(chan string, 100)
	lineChan := make(chan string, 100)

	var waitgroup sync.WaitGroup

	var root string
	if len(os.Args) == 1 {
		log.Fatal("No path given, Please specify path.")
		return
	}

	if root = os.Args[1]; root == "" {
		log.Fatal("No path given, Please specify path.")
		return
	}

	err := filePathWalkDir(root, fileChan)
	if err != nil {
		panic(err)
	}

	waitgroup.Add(2)
	go readFile(fileChan, lineChan, &waitgroup)
	go sortChars(lineChan, &waitgroup)
	waitgroup.Wait()
}
