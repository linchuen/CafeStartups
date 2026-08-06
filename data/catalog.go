package data

import (
	_ "embed"
	"encoding/json"
)

// Each period and function is kept in its own fixture so card balance can be
// reviewed by period and function without mixing unrelated cards.
//
//go:embed mvp-fixture-period-1-resource.json
var MVPFixturePeriod1Resource []byte

//go:embed mvp-fixture-period-1-product.json
var MVPFixturePeriod1Product []byte

//go:embed mvp-fixture-period-1-value.json
var MVPFixturePeriod1Value []byte

//go:embed mvp-fixture-period-1-channel.json
var MVPFixturePeriod1Channel []byte

//go:embed mvp-fixture-period-1-marketing.json
var MVPFixturePeriod1Marketing []byte

//go:embed mvp-fixture-period-2-resource.json
var MVPFixturePeriod2Resource []byte

//go:embed mvp-fixture-period-2-product.json
var MVPFixturePeriod2Product []byte

//go:embed mvp-fixture-period-2-value.json
var MVPFixturePeriod2Value []byte

//go:embed mvp-fixture-period-2-channel.json
var MVPFixturePeriod2Channel []byte

//go:embed mvp-fixture-period-2-marketing.json
var MVPFixturePeriod2Marketing []byte

//go:embed mvp-fixture-period-3-resource.json
var MVPFixturePeriod3Resource []byte

//go:embed mvp-fixture-period-3-product.json
var MVPFixturePeriod3Product []byte

//go:embed mvp-fixture-period-3-value.json
var MVPFixturePeriod3Value []byte

//go:embed mvp-fixture-period-3-channel.json
var MVPFixturePeriod3Channel []byte

//go:embed mvp-fixture-period-3-marketing.json
var MVPFixturePeriod3Marketing []byte

//go:embed mvp-fixture-partner.json
var MVPPartnerFixture []byte

type fixtureEnvelope struct {
	Source  string            `json:"source"`
	Version int               `json:"version"`
	Seed    string            `json:"seed"`
	Cards   []json.RawMessage `json:"cards"`
}

// Period fixtures preserve the existing catalog API while making each
// period/function distribution independently reviewable.
var MVPFixturePeriod1 = mergeFixtures(MVPFixturePeriod1Resource, MVPFixturePeriod1Product, MVPFixturePeriod1Value, MVPFixturePeriod1Channel, MVPFixturePeriod1Marketing)
var MVPFixturePeriod2 = mergeFixtures(MVPFixturePeriod2Resource, MVPFixturePeriod2Product, MVPFixturePeriod2Value, MVPFixturePeriod2Channel, MVPFixturePeriod2Marketing)
var MVPFixturePeriod3 = mergeFixtures(MVPFixturePeriod3Resource, MVPFixturePeriod3Product, MVPFixturePeriod3Value, MVPFixturePeriod3Channel, MVPFixturePeriod3Marketing)

// MVPFixture preserves the original single-fixture API while combining all
// period/function fixtures for the existing catalog loader.
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
