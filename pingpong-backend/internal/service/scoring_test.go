package service

import (
	"testing"

	"pingpong-backend/internal/model"
)

func TestCalculateServerUsesJavaCompatibleRotation(t *testing.T) {
	player1 := uint64(1)
	player2 := uint64(2)
	firstServer := player1

	tests := []struct {
		name     string
		scoreP1  int
		scoreP2  int
		expected string
	}{
		{name: "zero zero", scoreP1: 0, scoreP2: 0, expected: "player1"},
		{name: "after one point", scoreP1: 1, scoreP2: 0, expected: "player1"},
		{name: "after two points", scoreP1: 2, scoreP2: 0, expected: "player2"},
		{name: "before deuce by parity", scoreP1: 9, scoreP2: 10, expected: "player2"},
		{name: "ten all", scoreP1: 10, scoreP2: 10, expected: "player1"},
		{name: "eleven ten", scoreP1: 11, scoreP2: 10, expected: "player2"},
		{name: "eleven all", scoreP1: 11, scoreP2: 11, expected: "player1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := CalculateServer(test.scoreP1, test.scoreP2, &firstServer, &player1, &player2)
			if actual != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}

func TestCalculateNextScore(t *testing.T) {
	p1Score := 3
	p2Score := 2
	rallies := []model.Rally{
		{RallyNumber: 1, ScoreP1: intPtr(1), ScoreP2: intPtr(0)},
		{RallyNumber: 5, ScoreP1: &p1Score, ScoreP2: &p2Score},
	}

	scoreP1, scoreP2 := CalculateNextScore(nil, "player1")
	if scoreP1 != 1 || scoreP2 != 0 {
		t.Fatalf("first rally expected 1-0, got %d-%d", scoreP1, scoreP2)
	}

	scoreP1, scoreP2 = CalculateNextScore(rallies, "player2")
	if scoreP1 != 3 || scoreP2 != 3 {
		t.Fatalf("later rally expected 3-3, got %d-%d", scoreP1, scoreP2)
	}

	scoreP1, scoreP2 = CalculateNextScore(rallies, "invalid")
	if scoreP1 != 3 || scoreP2 != 2 {
		t.Fatalf("invalid scorer expected unchanged 3-2, got %d-%d", scoreP1, scoreP2)
	}
}

func TestMaxConsecutiveScore(t *testing.T) {
	rallies := []model.Rally{
		{RallyNumber: 4, Scorer: "player2"},
		{RallyNumber: 1, Scorer: "player1"},
		{RallyNumber: 2, Scorer: "player1"},
		{RallyNumber: 3, Scorer: "player2"},
		{RallyNumber: 5, Scorer: "player2"},
		{RallyNumber: 6, Scorer: "player1"},
	}

	if actual := MaxConsecutiveScore(rallies, "player1"); actual != 2 {
		t.Fatalf("player1 expected 2, got %d", actual)
	}
	if actual := MaxConsecutiveScore(rallies, "player2"); actual != 3 {
		t.Fatalf("player2 expected 3, got %d", actual)
	}
}

func intPtr(value int) *int {
	return &value
}
