package domain

import "testing"

func TestMVPPartnerCardsLoadFromCardFixture(t *testing.T) {
	cards := MVPPartnerCards()
	if len(cards) != 10 {
		t.Fatalf("partner cards=%d, want 10", len(cards))
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
