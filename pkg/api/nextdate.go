package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("incorrect repeat")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("incorrect dstart")
	}

	repeatParse := strings.Split(repeat, " ")
	switch repeatParse[0] {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	case "d":
		if len(repeatParse) != 2 {
			return "", fmt.Errorf("incorrect format repeat")
		}

		if repeatParse[1] == "" {
			return "", fmt.Errorf("incorrect format repeat")
		}

		interval, err := strconv.Atoi(repeatParse[1])
		if err != nil {
			return "", err
		}

		if interval > 400 {
			return "", fmt.Errorf("incorrect format repeat")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", fmt.Errorf("incorrect format repeat")
	}

	return date.Format(dateFormat), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	now, err := time.Parse(dateFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(nextDate))
}
