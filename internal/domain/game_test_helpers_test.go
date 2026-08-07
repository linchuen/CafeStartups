package domain

import "testing"

func gameForTest(t *testing.T) *Game {
	t.Helper()
	g, err := NewGame("fixed-seed", []string{"a", "b", "c", "d"})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	return g
}
