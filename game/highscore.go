package main

import (
	"encoding/gob"
	"os"
	"slices"

	"github.com/bennicholls/tyumi/log"
)

var highscore_filename string = "high.scores"

type HighScores struct {
	Scores []HighScoreEntry
}

func (hs HighScores) IsHighScore(score int) bool {
	if score == 0 {
		return false
	}

	if len(hs.Scores) < 10 {
		return true
	}

	return score > hs.Scores[len(hs.Scores)-1].Score
}

func (hs *HighScores) AddEntry(entry HighScoreEntry) {
	hs.Scores = append(hs.Scores, entry)
	slices.SortFunc(hs.Scores, func(e1, e2 HighScoreEntry) int {
		if e1.Score > e2.Score {
			return -1
		} else if e1.Score < e2.Score {
			return 1
		} else {
			return 0
		}
	})

	if len(hs.Scores) > 10 {
		hs.Scores = hs.Scores[0:10]
	}
}

func (hs *HighScores) WriteToDisk() {
	if len(hs.Scores) == 0 {
		return
	}

	file, err := os.Create(highscore_filename)
	if err != nil {
		log.Error("Could not open high score file: ", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(hs)
	log.Info("Wrote high scores to file high.scores")
}

func (hs *HighScores) LoadFromDisk() {
	hs.Scores = make([]HighScoreEntry, 0)

	file, err := os.Open(highscore_filename)
	if err != nil {
		log.Info("Could not open high score (maybe because it wasn't there?)")
		return
	}

	decoder := gob.NewDecoder(file)
	dhs := HighScores{}
	err = decoder.Decode(&dhs)
	if err != nil {
		log.Warning("Could not decode high score file. How bizarre ;)")
	}

	if len(dhs.Scores) > 0 {
		hs.Scores = dhs.Scores
	}
}

type HighScoreEntry struct {
	Name  string
	Score int
}
