package build

import (
	"fmt"
	"math"
	"strings"
)

func rgb(i int) (r, g, b int) {
	var f = 0.275

	return int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128),
		45,
		int(math.Sin(f*float64(i)+0)*127 + 128)
}

func rainbow(text string) string {
	rainbowStr := make([]string, 0, len(text))
	for index, value := range text {
		r, g, b := rgb(index)
		str := fmt.Sprintf("\033[1m\033[38;2;%d;%d;%dm%c\033[0m\033[0;1m", r, g, b, value)
		rainbowStr = append(rainbowStr, str)
	}

	return strings.Join(rainbowStr, "")
}
