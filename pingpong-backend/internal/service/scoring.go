package service

import (
	"sort"
	"strings"

	"pingpong-backend/internal/model"
)

func CalculateServer(scoreP1, scoreP2 int, firstServer, player1ID, player2ID *uint64) string {
	if firstServer == nil || player1ID == nil || player2ID == nil {
		return "Unknown"
	}

	totalScore := scoreP1 + scoreP2
	turns := 0
	if scoreP1 >= 10 && scoreP2 >= 10 {
		turns = 10 + (totalScore - 20)
	} else {
		turns = totalScore / 2
	}

	currentServer := *firstServer
	if turns%2 != 0 {
		if *firstServer == *player1ID {
			currentServer = *player2ID
		} else {
			currentServer = *player1ID
		}
	}

	if currentServer == *player1ID {
		return "player1"
	}
	if currentServer == *player2ID {
		return "player2"
	}
	return "Unknown"
}

func CalculateNextScore(rallies []model.Rally, scorer string) (int, int) {
	scoreP1, scoreP2 := LatestScore(rallies)
	if strings.EqualFold(scorer, "player1") {
		scoreP1++
	} else if strings.EqualFold(scorer, "player2") {
		scoreP2++
	}
	return scoreP1, scoreP2
}

func LatestScore(rallies []model.Rally) (int, int) {
	if len(rallies) == 0 {
		return 0, 0
	}
	latest := LatestRally(rallies)
	scoreP1, scoreP2 := 0, 0
	if latest.ScoreP1 != nil {
		scoreP1 = *latest.ScoreP1
	}
	if latest.ScoreP2 != nil {
		scoreP2 = *latest.ScoreP2
	}
	return scoreP1, scoreP2
}

func LatestRally(rallies []model.Rally) model.Rally {
	sort.Slice(rallies, func(i, j int) bool {
		return rallies[i].RallyNumber > rallies[j].RallyNumber
	})
	return rallies[0]
}

func MaxConsecutiveScore(rallies []model.Rally, player string) int {
	sort.Slice(rallies, func(i, j int) bool {
		return rallies[i].RallyNumber < rallies[j].RallyNumber
	})

	maxScore, current := 0, 0
	for _, rally := range rallies {
		if strings.EqualFold(rally.Scorer, player) {
			current++
			if current > maxScore {
				maxScore = current
			}
			continue
		}
		current = 0
	}
	return maxScore
}
