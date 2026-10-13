// Package game holds the game rules that don't depend on HTTP or the database.
package game

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Normalize lowercases s, strips accents and drops everything that isn't a letter
// or digit, so "Pingo-Doce", "pingo doce" and "PINGODOCE" all compare equal.
func Normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	stripped, _, err := transform.String(t, s)
	if err != nil {
		stripped = s
	}
	var b strings.Builder
	for _, r := range strings.ToLower(stripped) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Matches reports whether guess matches any of the accepted answers, allowing a
// few typos on longer names.
func Matches(guess string, accepted []string) bool {
	g := Normalize(guess)
	if g == "" {
		return false
	}
	for _, a := range accepted {
		n := Normalize(a)
		if levenshtein(g, n) <= typoTolerance(len([]rune(n))) {
			return true
		}
	}
	return false
}

// typoTolerance is how many edits we forgive for an answer of length n. Short
// names (MEO, NOS, CTT) must be exact, otherwise almost anything would match.
func typoTolerance(n int) int {
	switch {
	case n <= 4:
		return 0
	case n <= 8:
		return 1
	default:
		return 2
	}
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	cur := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(rb)]
}
