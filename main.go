package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"sync"
)

// Make a Regex to say we only want letters and numbers
var reg = regexp.MustCompile("[^a-zA-Z0-9]+")

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	c.v[key]++
	c.mux.Unlock()
}

func (c *SafeCounter) Stat(out io.Writer) {
	// we need to return the result sorted by values
	wordFrequencies := c.v
	countersMap := make(map[int]string, len(wordFrequencies))
	counters := make([]int, 0, len(wordFrequencies))
	for k, v := range wordFrequencies {
		counters = append(counters, v)
		countersMap[v] = k
	}

	sort.Sort(sort.Reverse(sort.IntSlice(counters)))

	for i := 0; i < 10; i++ {
		counter := counters[i]
		fmt.Fprintf(out, "%s\t%d\n", countersMap[counter], counter)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	c := SafeCounter{v: make(map[string]int)}

	go readFile(&wg, &c, "t/data/Crime&Punishment.txt")
	go readFile(&wg, &c, "t/data/War&Peace.txt")

	wg.Wait()

	c.Stat(os.Stdout)
}

func readFile(wg *sync.WaitGroup, c *SafeCounter, filePath string) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		processedWord := reg.ReplaceAllString(word, "")
		c.Inc(processedWord)
	}
}
