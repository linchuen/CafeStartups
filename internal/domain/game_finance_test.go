package domain

import "testing"

func TestInterestCanTakeLoanWhenCashIsShort(t *testing.T) {
	g, _ := NewGame("x", []string{"a", "b"})
	p := g.Players[0]
	p.Loans, p.Cash = 1, 0
	if err := g.SettleInterest(); err != nil {
		t.Fatal(err)
	}
	if p.Loans != 2 || p.Cash != 40 {
		t.Fatalf("loans=%d cash=%d", p.Loans, p.Cash)
	}
}

func TestSeedAndRevenueAreDeterministic(t *testing.T) {
	a := gameForTest(t)
	b := gameForTest(t)
	for i := range a.Players {
		if a.Players[i].Hand[0].ID != b.Players[i].Hand[0].ID {
			t.Fatal("same seed produced different deal")
		}
	}
	a.Players[0].Tableau = []Card{{Demand: map[string]int{"gourmet": 1}}}
	a.DistributeCustomers([]Customer{{Kind: "gourmet", Demand: "gourmet", Count: 2}, {Kind: "regular", Demand: "gourmet", Count: 1}})
	a.SettleRevenue()
	if a.Players[0].Revenue != 40 || a.Players[0].Cash != InitialCash-a.Players[0].StarterShop.Cost.Cash+40 {
		t.Fatalf("revenue=%d cash=%d", a.Players[0].Revenue, a.Players[0].Cash)
	}
}
