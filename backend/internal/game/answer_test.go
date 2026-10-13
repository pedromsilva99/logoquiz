package game

import "testing"

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"Pingo Doce":               "pingodoce",
		"Caixa Geral de Depósitos": "caixageraldedepositos",
		"  Delta Cafés! ":          "deltacafes",
		"Super-Bock":               "superbock",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMatches(t *testing.T) {
	cases := []struct {
		guess    string
		accepted []string
		want     bool
	}{
		{"pingo doce", []string{"Pingo Doce"}, true},
		{"pingo dose", []string{"Pingo Doce"}, true}, // one typo on a long name
		{"cgd", []string{"Caixa Geral de Depósitos", "CGD"}, true},
		{"caixa geral de depositos", []string{"Caixa Geral de Depósitos"}, true},
		{"mea", []string{"MEO"}, false}, // short names must be exact
		{"nos", []string{"NOS"}, true},
		{"", []string{"NOS"}, false},
		{"continente", []string{"Pingo Doce"}, false},
	}
	for _, c := range cases {
		if got := Matches(c.guess, c.accepted); got != c.want {
			t.Errorf("Matches(%q, %v) = %v, want %v", c.guess, c.accepted, got, c.want)
		}
	}
}
