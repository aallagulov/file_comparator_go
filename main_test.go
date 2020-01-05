package main

import (
	"bytes"
	"sync"
	"testing"
)

const expectedWarAndPeaceResult = `the	36634
and	20739
to	17737
of	16207
in	10111
a	9976
was	8070
that	7976
his	7882
he	7669
`

func TestReadFileWarAndPeace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	c := SafeCounter{v: make(map[string]int)}

	go readFile(&wg, &c, "t/data/War&Peace.txt")

	wg.Wait()

	out := new(bytes.Buffer)
	c.Stat(out)

	result := out.String()
	if result != expectedWarAndPeaceResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expectedWarAndPeaceResult)
	}
}

const expectedCrimeAndPunishmentResult = `the	8105
and	6712
to	5625
a	4406
of	3672
in	3304
he	3295
I	3168
you	2875
was	2870
`

func TestReadFileCrimeAndPunishment(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	c := SafeCounter{v: make(map[string]int)}

	go readFile(&wg, &c, "t/data/Crime&Punishment.txt")

	wg.Wait()

	out := new(bytes.Buffer)
	c.Stat(out)

	result := out.String()
	if result != expectedCrimeAndPunishmentResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expectedCrimeAndPunishmentResult)
	}
}

const expectedBothBooksResult = `the	44739
and	27451
to	23362
of	19879
a	14382
in	13415
he	10964
was	10940
that	10809
his	9963
`

func TestReadBothBooks(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	c := SafeCounter{v: make(map[string]int)}

	go readFile(&wg, &c, "t/data/War&Peace.txt")
	go readFile(&wg, &c, "t/data/Crime&Punishment.txt")

	wg.Wait()

	out := new(bytes.Buffer)
	c.Stat(out)

	result := out.String()
	if result != expectedBothBooksResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expectedBothBooksResult)
	}
}
