package domain

import (
	"testing"

	fixturedata "cafestartups/data"
)

func TestLoadCatalogRejectsDuplicateCards(t *testing.T) {
	if _, err := LoadCatalog([]byte(`{"cards":[{"id":"x","period":1,"cost":{"cash":0,"icons":[]}}, {"id":"x","period":1,"cost":{"cash":0,"icons":[]}}]}`)); err == nil {
		t.Fatal("expected duplicate card error")
	}
}

func TestCatalogDealsPeriodCards(t *testing.T) {
	catalog, err := LoadCatalog(fixturedata.MVPFixture)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame("catalog-seed", []string{"a", "b", "c", "d"})
	if err != nil {
		t.Fatal(err)
	}
	g.SetCatalog(catalog)
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	for _, p := range g.Players {
		for _, card := range p.Hand {
			if card.Period != PeriodOne || card.Name == "Management Card" {
				t.Fatalf("unexpected dealt card: %+v", card)
			}
		}
	}
}

func TestCatalogPreservesExplicitCostIcons(t *testing.T) {
	catalog, err := LoadCatalog(fixturedata.MVPFixture)
	if err != nil {
		t.Fatal(err)
	}

	for _, card := range catalog {
		if len(card.Cost.Icons) == 0 {
			continue
		}
		for index, icon := range card.Cost.Icons {
			if icon == "" {
				t.Fatalf("card %q cost icon[%d] is empty", card.ID, index)
			}
		}
	}
}

func TestCatalogKeepsBlankCostIconsBlank(t *testing.T) {
	catalog, err := LoadCatalog([]byte(`{"cards":[{"id":"multi-icon","name":"Multi icon","kind":"product","period":3,"cost":{"cash":35,"icons":[]},"icons":["coffee","operations"],"marketChange":{},"source":"mvp-fixture"}]}`))
	if err != nil {
		t.Fatal(err)
	}

	got := catalog[0].Cost.Icons
	if len(got) != 0 {
		t.Fatalf("cost icons=%v, want blank cost", got)
	}
}

func TestCatalogRejectsMarketingAndChannelCostIcons(t *testing.T) {
	for _, icon := range []string{"marketing", "channel"} {
		fixture := []byte(`{"cards":[{"id":"invalid-cost","name":"Invalid cost","kind":"resource","period":1,"cost":{"cash":10,"icons":["` + icon + `"]},"icons":["operations"],"marketChange":{},"source":"mvp-fixture"}]}`)
		if _, err := LoadCatalog(fixture); err == nil {
			t.Fatalf("expected %q cost icon to be rejected", icon)
		}
	}
}
