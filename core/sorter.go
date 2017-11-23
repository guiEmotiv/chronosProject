package core

import "strconv"

type SortString []string

func (a SortString) Len() int           { return len(a) }

func (a SortString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortString) Less(i, j int) bool {
	x1, err := strconv.Atoi(a[i])
	if err != nil {
		panic(err)
	}
	x2, err := strconv.Atoi(a[j])
	if err != nil {
		panic(err)
	}
	return x1 < x2
}
