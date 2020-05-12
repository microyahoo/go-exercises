package main

import (
	"fmt"
)

func lcs(a, b string) (maxLen int) {
	var m = len(a)
	var n = len(b)
	c := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		c[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		c[i][0] = 0
	}
	for i := 0; i <= n; i++ {
		c[0][i] = 0
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				c[i][j] = c[i-1][j-1] + 1
			} else if c[i][j-1] >= c[i-1][j] {
				c[i][j] = c[i][j-1]
			} else if c[i][j-1] < c[i-1][j] {
				c[i][j] = c[i-1][j]
			}
		}
	}
	fmt.Println(c)
	return c[m][n]
}

func main() {
	fmt.Println("LCS(abcddab, bdcaba)=", lcs("abcddab", "bdcaba"))
}
