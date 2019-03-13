package main

import (
    "sort"
    "sync"
    "math/rand"
    // "strconv"
)

// Structure to store key/value pair
type Pair struct {
    key    string
    value  int
}

// Implement sort interface for decreasing order
type DescendingPairs []Pair
func (p DescendingPairs) Len() int           { return len(p) }
func (p DescendingPairs) Less(i, j int) bool { return p[i].value > p[j].value }
func (p DescendingPairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

/* Actual Markov model data structure, thread safe */
type MarkovModel struct {
    sync.RWMutex
    // map of string slice to a map of string to an int
    m map[string]map[string]int
}

func (mm *MarkovModel) get(key string) (map[string]int, bool) {
    mm.RLock()
    value, hasKey := mm.m[key]
    mm.RUnlock()
    return value, hasKey
}

func (mm *MarkovModel) insert(keySlice []string, elem string) {
    key := ""
    for i := 0; i < settings.order; i++ {
        key += keySlice[i] + " "
    }
    // Check if current element exists at key
    // If so update the count
    // Else, add the key
    currentValue, hasKey := mm.get(key)
    mm.Lock()
    if hasKey {
        oldCount := currentValue[elem]
        // debug("Replacing: (" + key + ", " + strconv.Itoa(oldCount) + ") with (" + elem + ", " + (strconv.Itoa(oldCount + 1)) + ")")
        currentValue[elem] = oldCount + 1
    } else {
        // debug("Creating key: " + key + " with (" + elem + ", " + strconv.Itoa(1) + ")")
        mm.m[key] = map[string]int{elem: 1}
    }
    mm.Unlock()
}

/// Given a 'key', randomly choose the next element based on previous state
func (mm *MarkovModel) getNext(keySlice []string) string {
    key := ""
    for i := 0; i < settings.order; i++ {
        key += keySlice[i] + " "
    }
    // Get the possible next states based on the key
    internalMap := make(map[string]int)
    debug(key)

    mm.RLock()
    for k,v := range mm.m[key] {
      internalMap[k] = v
    }
    mm.RUnlock()

    // Initialize the internal data structure for easy sorting
    descendingPairs := make(DescendingPairs, len(internalMap))
    pairIndex := 0
    // Populate an array containing the key, value pairs for later sorting 
    for key, value := range internalMap {
        descendingPairs[pairIndex]  = Pair{key, value}
        pairIndex++
    }

    // Sort the pairs in decreasing order according to the count of a given letter
    sort.Sort(descendingPairs)
    
    // Create an array with cumulative values which will be used to decide on a move
    sum := 0
    cumulatives := make([]int, len(descendingPairs))
    for i := 0; i < len(descendingPairs); i++ {
        cumulatives[i] = descendingPairs[i].value + sum
        sum += descendingPairs[i].value
    }
    
    // Generate a random number, given that number:
    // Find the largest index of 'cumulatives' that is larger than that number 
    // This should be replaced with a binary search to speed things up
    random := rand.Intn(sum + 1)
    var nextIndex int
    for i := 0; i < len(cumulatives); i++ {
        if cumulatives[i] < random {
            nextIndex = i
        } else {
            break
        }
    }
    // The index that we matched our random number to is the same index
    // of the next value in our pairs index
    return descendingPairs[nextIndex].key
}