package domain

import "testing"

func TestMVPPartnerCardsLoadFromCardFixture(t *testing.T) {
	cards := MVPPartnerCards()
	if len(cards) != 13 {
		t.Fatalf("partner cards=%d, want 13", len(cards))
	}

	for _, card := range cards {
		if card.Kind != "partner" {
			t.Fatalf("card %q kind=%q, want partner", card.ID, card.Kind)
		}
		if card.ColorKey != card.Function {
			t.Fatalf("card %q color=%q does not match function=%q", card.ID, card.ColorKey, card.Function)
		}
		if len(card.Icons) == 0 || len(card.Cost.Icons) != 0 {
			t.Fatalf("card %q icons=%v costIcons=%v", card.ID, card.Icons, card.Cost.Icons)
		}
		if card.Function == "channel" {
			hasInitialCashBonus := card.Cost.Cash > 0
			hasCustomerCount := card.CustomerCount["gourmet"] > 0 || card.CustomerCount["regular"] > 0
			if hasInitialCashBonus == hasCustomerCount {
				t.Fatalf("channel card %q must have exactly one function", card.ID)
			}
		}
	}
}

func TestLoadPartnerCardsRejectsMismatchedColor(t *testing.T) {
	fixture := []byte(`{"cards":[{"id":"partner-test","kind":"partner","period":null,"cost":{"cash":0,"icons":[]},"icons":["operations"],"marketChange":{"gourmet":0,"regular":0,"difficult":0},"function":"resource","colorKey":"product","source":"mvp-fixture"}]}`)
	if _, err := LoadPartnerCards(fixture); err == nil {
		t.Fatal("expected mismatched partner color to be rejected")
	}
}

func TestMVPStarterShopCardsUseTenUnitCosts(t *testing.T) {
	cards := MVPStarterShopCards()
	if len(cards) != 11 {
		t.Fatalf("starter shop cards=%d, want 11", len(cards))
	}
	for _, card := range cards {
		if card.Kind != "starter_shop" {
			t.Fatalf("card %q kind=%q, want starter_shop", card.ID, card.Kind)
		}
		if card.Cost.Cash%10 != 0 {
			t.Fatalf("card %q cost=%d is not a 10-unit value", card.ID, card.Cost.Cash)
		}
		if card.CustomerCount["gourmet"] == 0 && card.CustomerCount["regular"] == 0 {
			t.Fatalf("card %q has no customer count", card.ID)
		}
		gourmet := card.CustomerCount["gourmet"]
		regular := card.CustomerCount["regular"]
		if gourmet+regular > 2 {
			t.Fatalf("card %q has %d customers, want at most 2", card.ID, gourmet+regular)
		}
		wantCost := gourmet*20 + regular*10
		if card.Cost.Cash != wantCost {
			t.Fatalf("card %q cost=%d, want gourmet*20 + regular*10 = %d", card.ID, card.Cost.Cash, wantCost)
		}
	}
}

func TestStarterShopCustomerCountAddsCustomers(t *testing.T) {
	g, err := NewGame("starter-shop-customers", []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	p := g.Players[0]
	p.Partner = Card{}
	p.StarterShop = Card{CustomerCount: map[string]int{"gourmet": 2}}
	p.Tableau = []Card{{Demand: map[string]int{"gourmet": 1}}}

	g.DistributeCustomers(nil)

	if len(p.Customers) != 1 || p.Customers[0].Kind != "gourmet" || p.Customers[0].Count != 2 || p.Customers[0].UnitPrice == 0 {
		t.Fatalf("starter shop customers=%+v", p.Customers)
	}
}
