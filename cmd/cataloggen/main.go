package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cafestartups/internal/domain"
)

type distributionFile struct {
	Periods map[string]map[string]int `json:"periods"`
}

type fixture struct {
	Source   string        `json:"source"`
	Version  int           `json:"version"`
	Seed     string        `json:"seed"`
	Function string        `json:"function"`
	Cards    []domain.Card `json:"cards"`
}

var functions = []string{"resource", "product", "value", "channel", "marketing"}

var effectPools = map[string][]string{
	"resource":  {"data", "procurement", "operations", "marketing_resource"},
	"product":   {"coffee", "dessert", "beans"},
	"value":     {"value", "taste", "service"},
	"channel":   {"channel"},
	"marketing": {"marketing"},
}

var costPool = []string{"data", "procurement", "operations", "marketing_resource", "coffee", "dessert", "beans", "value", "taste", "service"}

func main() {
	distribution := distributionFile{}
	readJSON(filepath.Join("data", "card-distribution.json"), &distribution)

	for period := 1; period <= 3; period++ {
		fixtures := make([]fixture, 0, len(functions))
		for _, function := range functions {
			path := filepath.Join("data", fmt.Sprintf("mvp-fixture-period-%d-%s.json", period, function))
			current := fixture{}
			readJSON(path, &current)
			wanted := distribution.Periods[fmt.Sprint(period)][function]
			current.Cards = resizeCards(current.Cards, wanted, function)
			balanceEffectIcons(current.Cards, function)
			fixtures = append(fixtures, current)
		}
		balanceCostIcons(fixtures)
		for index, function := range functions {
			path := filepath.Join("data", fmt.Sprintf("mvp-fixture-period-%d-%s.json", period, function))
			current := fixtures[index]
			writeJSON(path, current)
		}
	}
}

func resizeCards(cards []domain.Card, wanted int, function string) []domain.Card {
	if wanted < 1 || len(cards) == 0 {
		panic(fmt.Sprintf("function %s must contain at least one source card", function))
	}
	if len(cards) > wanted {
		return append([]domain.Card(nil), cards[:wanted]...)
	}
	for len(cards) < wanted {
		clone := cards[len(cards)%len(cards)]
		clone.ID = fmt.Sprintf("%s-balanced-%d", clone.ID, len(cards)+1)
		cards = append(cards, clone)
	}
	return cards
}

func balanceEffectIcons(cards []domain.Card, function string) {
	effectPool := effectPools[function]
	effectCounts := make([]int, len(effectPool))
	for cardIndex := range cards {
		for iconIndex := range cards[cardIndex].Icons {
			choice := leastUsed(effectCounts, effectPool, cards[cardIndex].Icons[:iconIndex])
			cards[cardIndex].Icons[iconIndex] = choice.value
			effectCounts[choice.index]++
		}
	}
}

func balanceCostIcons(fixtures []fixture) {
	costCounts := make([]int, len(costPool))
	for fixtureIndex := range fixtures {
		for cardIndex := range fixtures[fixtureIndex].Cards {
			for iconIndex := range fixtures[fixtureIndex].Cards[cardIndex].Cost.Icons {
				icons := fixtures[fixtureIndex].Cards[cardIndex].Cost.Icons
				choice := leastUsed(costCounts, costPool, icons[:iconIndex])
				fixtures[fixtureIndex].Cards[cardIndex].Cost.Icons[iconIndex] = choice.value
				costCounts[choice.index]++
			}
		}
	}
}

type iconChoice struct {
	index int
	value string
}

func leastUsed(counts []int, pool, current []string) iconChoice {
	best := iconChoice{index: -1}
	for index, value := range pool {
		if contains(current, value) {
			continue
		}
		if best.index == -1 || counts[index] < counts[best.index] {
			best = iconChoice{index: index, value: value}
		}
	}
	if best.index >= 0 {
		return best
	}
	return iconChoice{index: 0, value: pool[0]}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func readJSON(path string, target any) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		panic(fmt.Errorf("decode %s: %w", path, err))
	}
}

func writeJSON(path string, value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		panic(err)
	}
}
