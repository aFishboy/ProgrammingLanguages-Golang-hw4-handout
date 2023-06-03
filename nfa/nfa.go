package nfa

import "sync"

// A nondeterministic Finite Automaton (NFA) consists of states,
// symbols in an alphabet, and a transition function.

// A state in the NFA is represented as an unsigned integer.
type state uint

// Given the current state and a symbol, the transition function
// of an NFA returns the set of next states the NFA can transition to
// on reading the given symbol.
// This set of next states could be empty.
type TransitionFunction func(st state, sym rune) []state

// Reachable returns true if there exists a sequence of transitions
// from `transitions` such that if the NFA starts at the start state
// `start` it would reach the final state `final` after reading the
// entire sequence of symbols `input`; Reachable returns false otherwise.
func Reachable(
	// `transitions` tells us what our NFA looks like
	transitions TransitionFunction,
	// `start` and `final` tell us where to start, and where we want to end up
	start, final state,
	// `input` is a (possible empty) list of symbols to apply.
	input []rune,
) bool {
	if len(input) == 0 {
		return start == final
	}

	var wg sync.WaitGroup
	ch := make(chan bool)

	for _, index := range transitions(start, input[0]) {
		wg.Add(1)
		go func(st state, sym rune) {
			defer wg.Done()
			ch <- Reachable(transitions, st, final, input[1:])
		}(index, input[0])
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		if result == true {
			return true
		}
	}

	return false
}

