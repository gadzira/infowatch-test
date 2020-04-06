package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chenjiandongx/go-echarts/charts"
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

func sortChars(l chan string, m chan map[string]int, waitgroup *sync.WaitGroup) map[string]int {
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
	m <- sourceMap
	// for test
	// for k, v := range sourceMap {
	// fmt.Printf("[%s] : %d\n", k, v)
	// }
	waitgroup.Done()
	return sourceMap
}

func renderTheHistogram(mapChan chan map[string]int, waitgroup *sync.WaitGroup) {
	keyItems := []string{}
	valueItems := []int{}
	for m := range mapChan {
		for k, v := range m {
			keyItems = append(keyItems, k)
			valueItems = append(valueItems, v)
		}

		bar := charts.NewBar()
		bar.SetSeriesOptions(
			charts.BarOpts{BarCategoryGap: "170%"})
		bar.SetGlobalOptions(
			charts.TitleOpts{Title: "For infowatch", Right: "80%"},
			charts.InitOpts{Width: "1900px", Height: "900px"},
		)
		bar.AddXAxis(keyItems).
			AddYAxis("Symbols", valueItems,
				charts.ColorOpts{"lightblue"})
		f, err := os.Create("bar.html")
		if err != nil {
			log.Println(err)
		}
		bar.Render(f)
		fmt.Println("The chart was rendered")
		os.Exit(0)
	}
	waitgroup.Done()
}

func main() {

	fileChan := make(chan string, 100)
	lineChan := make(chan string, 100)
	mapChan := make(chan map[string]int, 100)

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

	waitgroup.Add(3)
	go readFile(fileChan, lineChan, &waitgroup)
	go sortChars(lineChan, mapChan, &waitgroup)
	go renderTheHistogram(mapChan, &waitgroup)
	waitgroup.Wait()
}
