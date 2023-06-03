package smash

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

type word string

func Smash(reader io.Reader, smasher func(word) uint32) map[uint32]uint {
	smashMap := make(map[uint32]uint)
	scanner := bufio.NewScanner(reader)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	lines := make(chan string) 

	go func() {
		defer close(lines)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	for line := range lines {
		wg.Add(1)

		go func(line string) {
			defer wg.Done()

			words := strings.Fields(line)

			for _, w := range words {
				word := word(w)
				hash := smasher(word)

				mutex.Lock()
				smashMap[hash]++
				mutex.Unlock()
			}
		}(line)
	}

	wg.Wait()

	return smashMap
}