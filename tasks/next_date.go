package tasks

import (
	"errors"
	"fmt"
	"go_final_project/str"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func normalizeToDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	nextDate, err := time.ParseInLocation(str.DateFormat, date, time.UTC)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	now = normalizeToDate(now.UTC())
	nextDate = normalizeToDate(nextDate.UTC())

	repeatRule := strings.Fields(repeat)

	switch repeatRule[0] {
	case "d":
		if len(repeatRule) < 2 {
			return "", errors.New("не указано количество дней")
		}

		days, err := strconv.Atoi(repeatRule[1])
		if err != nil {
			return "", fmt.Errorf("неверное значение для дней: %v", err)
		}

		if days > 400 {
			return "", fmt.Errorf("количество дней не должно превышать 400")
		}

		nextDate = nextDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format(str.DateFormat), nil

	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		if nextDate.Month() == time.February && nextDate.Day() == 29 && !isLeapYear(nextDate.Year()) {
			nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
		}
		return nextDate.Format(str.DateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", repeatRule[0])
	}
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "Отсутствуют обязательные параметры", http.StatusBadRequest)
		return
	}

	now, err := time.ParseInLocation(str.DateFormat, nowStr, time.UTC)
	if err != nil {
		http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
