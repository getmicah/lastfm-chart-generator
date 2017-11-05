package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	user, period, size, err := parse(args)
	if err != nil {
		fmt.Printf("Error: %s\n%s\n", err, usage())
	} else {
		chart(user, period, size)
	}
}

func parse(args []string) (string, string, int, error) {
	if len(args) != 3 {
		return "", "", 0, errors.New("Invalid arguments")
	}
	user := args[0]
	var period string
	switch args[1] {
	case "week":
		period = "7day"
	case "month":
		period = "1month"
	case "3month":
		period = "3month"
	case "6month":
		period = "6month"
	case "year":
		period = "12month"
	case "overall":
		period = "overall"
	default:
		return "", "", 0, errors.New("Invalid period")
	}
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return "", "", 0, err
	}
	if !(size == 3 || size == 4 || size == 5) {
		return "", "", 0, errors.New("Invalid size")
	}
	return user, period, size, nil
}

func usage() string {
	return "Usage: Usage: lcg <user> <period> <size>\nParams:\nuser\t<last.fm username>\nperiod\tweek, month, 3month, 6month, year, overall\nsize\t3, 4, 5"
}
