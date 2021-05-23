package lib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func Input(message string) string {
	var inp string

	fmt.Println(message)
	fmt.Print(" > ")
	fmt.Scanln(&inp)

	return inp
}

func SetupConditionOptions() ConditionOptions {
	pass := []rune(Input("Set the color pass you wish to compare to [R, G, B, A, and V (where V is the value in HSV)]"))
	thresholdStr := Input("Set the threshold for brightness (0 - 255)")
	inverseStr := Input("Should the bot click when it is less bright than threshold (yes / no)")

	inverse := strings.ToLower(inverseStr) == "yes"

	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		log.Fatal("Did not input valid integer for threshold")
	}

	return ConditionOptions{
		Pass:      pass[0],
		Threshold: threshold,
		Inverse:   inverse,
	}
}
