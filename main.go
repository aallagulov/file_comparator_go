package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

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

// func (c *SafeCounter) Value(key string) int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
// 	return c.v[key]
// }

func (c *SafeCounter) Stat() string {
	// we need to return the result sorted by values
	wordFrequencies := c.v
	countersMap := make(map[int]string, len(wordFrequencies))
	counters := make([]int, 0, len(wordFrequencies))
	for k, v := range wordFrequencies {
		counters = append(counters, v)
		countersMap[v] = k
	}

	sort.Sort(sort.Reverse(sort.IntSlice(counters)))

	b := new(bytes.Buffer)
	for i := 0; i < 10; i++ {
		counter := counters[i]
		fmt.Fprintf(b, "%s\t%d\n", countersMap[counter], counter)
	}
	return b.String()
}

func main() {
	// f1_path := flag.String("f1", "", "path to the 1st file")
	// f2_path := flag.String("f2", "", "path to the 2nd file")

	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("t/data/Crime&Punishment.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	sl1 := make([]string, 0, 0)
	for scanner.Scan() {
		word := scanner.Text()
		processedWord := reg.ReplaceAllString(word, "")
		sl1 = append(sl1, processedWord)
	}

	str2 := "asd cat  tac asd"
	sl2 := strings.Fields(str2)

	var wg sync.WaitGroup
	wg.Add(2)

	c := SafeCounter{v: make(map[string]int)}

	countFunc := func(words []string) {
		defer wg.Done()
		for i := 0; i < len(words); i++ {
			c.Inc(words[i])
		}
	}

	go countFunc(sl1)
	go countFunc(sl2)

	wg.Wait()

	fmt.Println(c.Stat())
}
