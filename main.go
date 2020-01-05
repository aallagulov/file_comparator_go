package main

import (
	"bufio"
	"flag"
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

func (c *SafeCounter) Stat(out io.Writer, outLimit int) {
	// we need to return the result sorted by values
	wordFrequencies := c.v
	countersMap := make(map[int]string, len(wordFrequencies))
	counters := make([]int, 0, len(wordFrequencies))
	for k, v := range wordFrequencies {
		counters = append(counters, v)
		countersMap[v] = k
	}

	sort.Sort(sort.Reverse(sort.IntSlice(counters)))

	wordsAmount := len(counters)
	var limit int
	if outLimit > 0 && wordsAmount > outLimit {
		limit = outLimit
	} else {
		limit = wordsAmount
	}
	for i := 0; i < limit; i++ {
		counter := counters[i]
		fmt.Fprintf(out, "%s\t%d\n", countersMap[counter], counter)
	}
}

func main() {
	var filePath1, filePath2, outLimit = getCLIParams()

	var wg sync.WaitGroup
	wg.Add(2)

	c := SafeCounter{v: make(map[string]int)}

	go readFile(&wg, &c, filePath1)
	go readFile(&wg, &c, filePath2)

	wg.Wait()

	c.Stat(os.Stdout, outLimit)
}

func getCLIParams() (string, string, int) {
	filePath1 := flag.String("f1", "", "path to the 1st file")
	filePath2 := flag.String("f2", "", "path to the 2nd file")
	outLimit := flag.Int("n", 10, "limit for outputted lines")

	flag.Parse()

	if *filePath1 == "" || *filePath2 == "" {
		fmt.Println("Please specify both -f1 and -f2 paths to comparing files")
		os.Exit(1)
	}

	return *filePath1, *filePath2, *outLimit
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
