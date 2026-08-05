package data

import (
	_ "embed"
	"encoding/json"
)

// Each period is kept in its own fixture so card updates stay localized.
//
//go:embed mvp-fixture-period-1.json
var MVPFixturePeriod1 []byte

//go:embed mvp-fixture-period-2.json
var MVPFixturePeriod2 []byte

//go:embed mvp-fixture-period-3.json
var MVPFixturePeriod3 []byte

//go:embed mvp-fixture-partner.json
var MVPPartnerFixture []byte

type fixtureEnvelope struct {
	Source  string            `json:"source"`
	Version int               `json:"version"`
	Seed    string            `json:"seed"`
	Cards   []json.RawMessage `json:"cards"`
}

// MVPFixture preserves the original single-fixture API while combining the
// period files at startup for the existing catalog loader.
var MVPFixture = mergeFixtures(MVPFixturePeriod1, MVPFixturePeriod2, MVPFixturePeriod3)

func mergeFixtures(fixtures ...[]byte) []byte {
	merged := fixtureEnvelope{}
	for _, fixture := range fixtures {
		var part fixtureEnvelope
		if err := json.Unmarshal(fixture, &part); err != nil {
			panic(err)
		}
		if merged.Source == "" {
			merged.Source, merged.Version, merged.Seed = part.Source, part.Version, part.Seed
		}
		merged.Cards = append(merged.Cards, part.Cards...)
	}
	data, err := json.Marshal(merged)
	if err != nil {
		panic(err)
	}
	return data
}
