package moderationhelpers

// similarityMin returns the minimum of three integers
func similarityMin(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

func levenshteinDistance(s1, s2 []rune) int {
	m, n := len(s1), len(s2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = similarityMin(
					dp[i-1][j-1]+1,
					dp[i-1][j]+1,
					dp[i][j-1]+1,
				)
			}
		}
	}

	return dp[m][n]
}

func computeTextSimilarity(text1, text2 string) float64 {
	runes1 := []rune(text1)
	runes2 := []rune(text2)

	distance := levenshteinDistance(runes1, runes2)

	maxLen := len(runes1)
	if len(runes2) > maxLen {
		maxLen = len(runes2)
	}

	if maxLen == 0 {
		return 100.0
	}

	similarity := (1.0 - float64(distance)/float64(maxLen)) * 100.0
	if similarity < 0 {
		similarity = 0
	}
	return similarity
}

func (c *ModerationHelpers) ContainsSimilar(
	text string,
	list []string,
	maxSimilarity float64,
) bool {
	for _, item := range list {
		similarity := computeTextSimilarity(text, item)
		if similarity >= maxSimilarity {
			return true
		}
	}

	return false
}
