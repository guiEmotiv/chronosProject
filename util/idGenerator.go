package util

import (
	"math/rand"
	"time"
)

const dict = "0123456789abcdefghijklmnopqrstvwxyz"

func init()  {
	rand.Seed(time.Now().Unix())
}

func internalGen(size int) []rune {
	bricks := make([]rune, 0)
	for i:=0;i<size;i++{
		nPos := rand.Intn(len(dict))
		bricks = append(bricks, rune(dict[nPos]))
	}
	return bricks
}

func NextId(size int) string {
	runes := internalGen(size)
	finalId := ""
	for _, r := range runes {
		finalId += string(r)
	}
	return finalId
}

