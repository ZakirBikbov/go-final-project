package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	rule := parts[0]

	switch rule {
	case "y":
		return handleYearly(date, now)
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("d rule requires interval")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("invalid interval for d rule")
		}
		return handleDaily(date, now, interval)
	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("w rule requires days")
		}
		return handleWeekly(date, now, parts[1])
	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("m rule requires days")
		}
		months := ""
		if len(parts) >= 3 {
			months = parts[2]
		}
		return handleMonthly(date, now, parts[1], months)
	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func handleYearly(date, now time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

func handleDaily(date, now time.Time, interval int) (string, error) {
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

func handleWeekly(date, now time.Time, daysStr string) (string, error) {
	days := make(map[int]bool)
	for _, d := range strings.Split(daysStr, ",") {
		d = strings.TrimSpace(d)
		day, err := strconv.Atoi(d)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("invalid day of week")
		}
		days[day] = true
	}

	for i := 0; i < 365*2; i++ {
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		if days[weekday] && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}
	return "", fmt.Errorf("no valid date found")
}

func handleMonthly(date, now time.Time, daysStr, monthsStr string) (string, error) {
	dayList := []int{}
	for _, d := range strings.Split(daysStr, ",") {
		d = strings.TrimSpace(d)
		day, err := strconv.Atoi(d)
		if err != nil {
			return "", fmt.Errorf("invalid day of month")
		}
		if day < -2 || day == 0 || day > 31 {
			return "", fmt.Errorf("invalid day of month")
		}
		dayList = append(dayList, day)
	}

	monthMap := make(map[int]bool)
	if monthsStr != "" {
		for _, m := range strings.Split(monthsStr, ",") {
			m = strings.TrimSpace(m)
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("invalid month")
			}
			monthMap[month] = true
		}
	}

	var dayArr [32]bool
	for _, d := range dayList {
		if d > 0 {
			dayArr[d] = true
		}
	}

	var monthArr [13]bool
	if len(monthMap) > 0 {
		for m := range monthMap {
			monthArr[m] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthArr[i] = true
		}
	}

	startDate := date
	if startDate.Before(now) {
		startDate = now
	}
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	currentYear := startDate.Year()
	currentMonth := int(startDate.Month())

	for monthOffset := 0; monthOffset < 24; monthOffset++ {
		year := currentYear
		month := currentMonth + monthOffset
		if month > 12 {
			year += (month - 1) / 12
			month = ((month - 1) % 12) + 1
		}

		if !monthArr[month] {
			continue
		}

		var candidateDates []time.Time

		nextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, date.Location())
		lastDayOfMonth := nextMonth.AddDate(0, 0, -1)
		prevDayOfMonth := lastDayOfMonth.AddDate(0, 0, -1)

		for _, targetDay := range dayList {
			var checkDate time.Time
			switch targetDay {
			case -1:
				checkDate = lastDayOfMonth
			case -2:
				checkDate = prevDayOfMonth
			default:
				checkDate = time.Date(year, time.Month(month), targetDay, 0, 0, 0, 0, date.Location())
				if checkDate.Month() != time.Month(month) {
					continue
				}
			}

			if !checkDate.Before(startDate) && afterNow(checkDate, now) {
				candidateDates = append(candidateDates, checkDate)
			}
		}

		if len(candidateDates) > 0 {
			earliest := candidateDates[0]
			for _, d := range candidateDates[1:] {
				if d.Before(earliest) {
					earliest = d
				}
			}
			return earliest.Format(DateFormat), nil
		}
	}

	return "", fmt.Errorf("no valid date found")
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr != "" {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now parameter: %v", err), http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}

	if dateStr == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}
