package main

import (
    // "database/sql"
    // _ "github.com/go-sql-driver/mysql"
    "strings"
)

// Given a corpus, break text down into tokens
// Tokens can be a:
//  -- Start/End
//  -- Word
//  -- Puncuation mark
//      -- Marking end of sentence
//      -- Other
//  -- 'Blocks'
//      -- Quote
//      -- Parenthetical
//  -- Hashtag
//  -- @
//  -- Hyperlink
//      -- Domain (including protocol)
//          -- 
//  -- Emoji

// Takes a single tweet and breaks it down into tokens
func tokenize(tweet string) {
    tweet = strings.Trim(tweet, " ")

    tweet = padding + " " + tweet + " " + padding

    tweetSlice := strings.Fields(tweet)

    for i := 0; i < len(tweetSlice) - settings.order; i++ {
        key := tweetSlice[i:i+settings.order]
        markov.insert(key, tweetSlice[i+settings.order])
        // debug(key[0] + ", " + key[1] + ": " + tweetSlice[i+settings.order])
    }
}

