package main

import (
	"flag"
	"fmt"
	"os"
	"reversesearch"
	"strconv"
	"sync"
)

type RSGraph struct {
	Root   []int
	Vert   []int
	Maxdeg int
}

type IndexError struct {
	Index  int
	Maxdeg int
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("Index %v out of bounds: 0 to %v", e.Index, e.Maxdeg-1)
}

func (v *RSGraph) Adjacent(i int) error {
	var maxdeg = v.MaxDeg()

	if (i < 0) || (i > maxdeg-1) {
		return &IndexError{i, maxdeg}
	}

	v.Vert[i], v.Vert[i+1] = v.Vert[i+1], v.Vert[i]

	return nil
}

func (v *RSGraph) Localsearch() {
	var maxdeg = v.MaxDeg()

	for i := 0; i < maxdeg; i++ {
		if v.Vert[i] > v.Vert[i+1] {
			v.Vert[i], v.Vert[i+1] = v.Vert[i+1], v.Vert[i]
			break
		}
	}
	return
}

func (v *RSGraph) Equal(t reversesearch.ReverseSearch) bool {
	w := t.(*RSGraph)

	if v.MaxDeg() != w.MaxDeg() {
		return false
	}

	for i, val := range v.Vert {
		if val != w.Vert[i] {
			return false
		}
	}

	return true
}

func (v *RSGraph) Copy() reversesearch.ReverseSearch {
	var w RSGraph

	w.Maxdeg = v.Maxdeg
	w.Root = make([]int, len(v.Root))
	copy(w.Root, v.Root)
	w.Vert = make([]int, len(v.Vert))
	copy(w.Vert, v.Vert)

	return &w
}

func (v *RSGraph) IsRoot() bool {

	for i, val := range v.Vert {
		if val != v.Root[i] {
			return false
		}
	}

	return true
}

func (v *RSGraph) MaxDeg() int {
	return v.Maxdeg
}

func (v *RSGraph) Output(t interface{}) {
	c := t.(chan []int)
	newSlice := make([]int, len(v.Vert))
	copy(newSlice, v.Vert)
	c <- newSlice
}

func permUsage() {
	fmt.Println("Usage: ./perm length-of-permutation (length should be an integer <= 20,")
	fmt.Println("                                     pending move to BigInt")
	flag.PrintDefaults()
}

func main() {
	var length int
	var count int = 0
	var w RSGraph
	var wg sync.WaitGroup

	flag.Usage = permUsage
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		permUsage()
		os.Exit(1)
	}

	if len, err := strconv.Atoi(args[0]); err != nil {
		permUsage()
		os.Exit(1)
	} else {
		length = len
	}

	if length > 20 {
		permUsage()
		os.Exit(1)
	}

	// Initialize permutation

	w.Root = make([]int, length)
	w.Vert = make([]int, length)
	w.Maxdeg = length - 1

	for i := 0; i < length; i++ {
		w.Root[i] = i + 1
		w.Vert[i] = i + 1
	}

	c := make(chan []int, 100)

	go reversesearch.Enumerate(&w, c)

	wg.Add(1)
	go func() {
		for {
			vert, ok := <-c
			if ok == false {
				break
			} else {
				count++
				fmt.Println(vert)
			}
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Printf("Number of permutations is %v\n", count)
}
