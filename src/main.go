package main

import (
    "bufio"     // Used for input/ouput functions
    "fmt"       // Used for printing to console
    "math/rand" // Used for generating random numbers
    "os"        // Used for os file access
    // "sort"      // Used to output prettily
    "strconv"   // Used for parsing int from string
    "strings"   // Used for cleaning up strings
    "sync"      // Used for concurrency syncing
    "time"      // Used for getting time
)

/* global variable declaration */
var markov MarkovModel
var padding string
var settings Settings
var verbose bool
var baseQuery []string

/* Concurrent Friendly Data Structures */
type Settings struct {
    fileName  string
    minLength int
    maxLength int
    order     int
    reload    bool
}

type ConcurrentMap struct {
    counter struct {
        sync.RWMutex
        values map[string]int
    }
}

func (cs *ConcurrentMap) Append(value string) {
    cs.counter.Lock()
    cs.counter.values[value] = 1
    cs.counter.Unlock()
}

func (cs *ConcurrentMap) Exists(value string) bool {
    cs.counter.RLock()
    _, exists := cs.counter.values[value]
    cs.counter.RUnlock()
    return exists
}

func (cs *ConcurrentMap) Length() int {
    cs.counter.RLock()
    length := len(cs.counter.values)
    cs.counter.RUnlock()
    return length
}

func debug(msg string) {
    if verbose {
        fmt.Println(msg)
    }
}

/* Request the parameters for our name generation from the end-user */
func getSettings() {
    // Get user input:

    // Ask for gender
    input := bufio.NewReader(os.Stdin)
    fmt.Print("Path to a file with tweets: ")

    fileName, _ := input.ReadString('\n')

    // Ask for model order
    fmt.Print("What should be the order of the Markov model (default 2)? ")

    rawOrder, _ := input.ReadString('\n')

    // Ask for number of names to generate
    fmt.Print("Should the tweet database be reloaded? ")

    rawReload, _ := input.ReadString('\n')

    // Process inputs
    fileName      = strings.Trim(strings.Trim(fileName, "\n"), " ")
    order, _     := strconv.Atoi(strings.Trim(rawOrder, "\n"))
    intReload, _  := strconv.Atoi(strings.Trim(rawReload, "\n"))
    var reload bool
    if intReload == 1 {
        reload = true
    } else {
        reload = false
    }

    settings = Settings{fileName, 20, 140, order, reload}
}

/* Format everything prettily */
func prettify(rawTweet []string) {
    tweet := ""
    for _, word := range rawTweet {
        if word != `\` {
            tweet += word + " "
        }
    }
    fmt.Println("\n" + tweet + "\n")
}

/// Generate a tweet unil we get to the end of a tweet "\"
func generateTweet() []string {
    query := make([]string, settings.order)
    copy(query, baseQuery)
    tweet := make([]string, settings.order)
    copy(tweet, query)
    var next string
    for {
        next = markov.getNext(query)
        if next == `\` {
            break
        }
        tweet = append(tweet, next)
        query = tweet[len(tweet)-settings.order:]
    }

    return tweet
}

// Get input from the end user and generate names accordingly
func main() {
    // Get user input:
    // getSettings()
    settings = Settings{"/Users/tlagran/Documents/Code/Go/tweet-generator/data/trumpTweets.txt", 20, 140, 2, true}
    verbose = false

    // Initialize markov model
    markov = MarkovModel{}
    markov.m = make(map[string]map[string]int)

    // Add backslashes(s) as padding to detect start and end of tweet
    padding = ""
    baseQuery = make([]string, settings.order)
    for i := 0; i < settings.order; i++ {
        padding = padding + `\ `
        baseQuery[i] = `\`
    }

    // Open selected corpus file
    corpusFileName := settings.fileName
    corpusFile, error := os.Open(corpusFileName)
    if error != nil {
        fmt.Println(error)
    }
    corpus := bufio.NewScanner(corpusFile)

    // Read in selected corpus line by line (one tweet per line)
    var tokenWaitGroup sync.WaitGroup
    for corpus.Scan() {
        tokenWaitGroup.Add(1)
        go func(tweet string) {
            tokenize(tweet)
            tokenWaitGroup.Done()
        }(corpus.Text())
    }
    tokenWaitGroup.Wait()
    corpusFile.Close()

    // Generate new tweets
    rand.Seed(time.Now().UnixNano())
    
    var generateWaitGroup sync.WaitGroup
    for {
        generateWaitGroup.Add(1)
        go func() {
            newTweet := generateTweet()
            prettify(newTweet)
            generateWaitGroup.Done()
        }()
        fmt.Scanln()
    }
    generateWaitGroup.Wait()
}