package datetimepicker

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const cronFieldCount = 5
const cronSearchYears = 5

func NewCronSchedule(expression string) (Schedule, error) {
	parsed, err := parseCronSchedule(expression)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

type cronSchedule struct {
	minutes       map[int]struct{}
	hours         map[int]struct{}
	daysOfMonth   map[int]struct{}
	months        map[int]struct{}
	daysOfWeek    map[int]struct{}
	dayOfMonthAll bool
	dayOfWeekAll  bool
}

func (s cronSchedule) Next(after time.Time) time.Time {
	candidate := after.Truncate(time.Minute).Add(time.Minute)
	limit := after.AddDate(cronSearchYears, 0, 0)
	for !candidate.After(limit) {
		if s.matches(candidate) {
			return candidate
		}
		candidate = candidate.Add(time.Minute)
	}

	return time.Time{}
}

func (s cronSchedule) matches(candidate time.Time) bool {
	if !setContains(s.minutes, candidate.Minute()) {
		return false
	}
	if !setContains(s.hours, candidate.Hour()) {
		return false
	}
	if !setContains(s.months, int(candidate.Month())) {
		return false
	}

	dayMatches := setContains(s.daysOfMonth, candidate.Day())
	weekMatches := setContains(s.daysOfWeek, cronWeekday(candidate))
	if s.dayOfMonthAll && s.dayOfWeekAll {
		return true
	}
	if s.dayOfMonthAll {
		return weekMatches
	}
	if s.dayOfWeekAll {
		return dayMatches
	}

	return dayMatches || weekMatches
}

func parseCronSchedule(expression string) (cronSchedule, error) {
	fields := strings.Fields(expression)
	if len(fields) != cronFieldCount {
		return cronSchedule{}, cronFieldCountError(expression)
	}

	minutes, _, minuteErr := parseCronField(fields[0], 0, 59, false)
	if minuteErr != nil {
		return cronSchedule{}, fmt.Errorf("parse cron minute: %w", minuteErr)
	}
	hours, _, hourErr := parseCronField(fields[1], 0, 23, false)
	if hourErr != nil {
		return cronSchedule{}, fmt.Errorf("parse cron hour: %w", hourErr)
	}
	daysOfMonth, dayOfMonthAll, dayErr := parseCronField(
		fields[2],
		1,
		31,
		false,
	)
	if dayErr != nil {
		return cronSchedule{}, fmt.Errorf("parse cron day-of-month: %w", dayErr)
	}
	months, _, monthErr := parseCronField(fields[3], 1, 12, false)
	if monthErr != nil {
		return cronSchedule{}, fmt.Errorf("parse cron month: %w", monthErr)
	}
	daysOfWeek, dayOfWeekAll, weekErr := parseCronField(fields[4], 0, 7, true)
	if weekErr != nil {
		return cronSchedule{}, fmt.Errorf("parse cron day-of-week: %w", weekErr)
	}

	return cronSchedule{
		minutes:       minutes,
		hours:         hours,
		daysOfMonth:   daysOfMonth,
		months:        months,
		daysOfWeek:    daysOfWeek,
		dayOfMonthAll: dayOfMonthAll,
		dayOfWeekAll:  dayOfWeekAll,
	}, nil
}

func cronFieldCountError(expression string) error {
	return fmt.Errorf("cron %q must have 5 fields", expression)
}

func parseCronField(
	raw string,
	minimum int,
	maximum int,
	weekday bool,
) (map[int]struct{}, bool, error) {
	values := make(map[int]struct{})
	all := raw == "*"
	for _, part := range strings.Split(raw, ",") {
		partValues, err := parseCronFieldPart(part, minimum, maximum, weekday)
		if err != nil {
			return nil, false, err
		}
		for _, value := range partValues {
			values[value] = struct{}{}
		}
	}

	return values, all, nil
}

func parseCronFieldPart(
	raw string,
	minimum int,
	maximum int,
	weekday bool,
) ([]int, error) {
	base, step, err := splitCronStep(raw)
	if err != nil {
		return nil, err
	}

	start, end, err := parseCronRange(base, minimum, maximum)
	if err != nil {
		return nil, err
	}

	values := make([]int, 0, ((end-start)/step)+1)
	for value := start; value <= end; value += step {
		values = append(values, normalizeCronValue(value, weekday))
	}

	return values, nil
}

func splitCronStep(raw string) (string, int, error) {
	parts := strings.Split(raw, "/")
	if len(parts) > 2 {
		return "", 0, fmt.Errorf("invalid step %q", raw)
	}
	if len(parts) == 1 {
		return raw, 1, nil
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil || step < 1 {
		return "", 0, fmt.Errorf("invalid step %q", raw)
	}

	return parts[0], step, nil
}

func parseCronRange(
	raw string,
	minimum int,
	maximum int,
) (int, int, error) {
	if raw == "*" {
		return minimum, maximum, nil
	}

	parts := strings.Split(raw, "-")
	if len(parts) > 2 {
		return 0, 0, fmt.Errorf("invalid range %q", raw)
	}

	start, err := parseCronNumber(parts[0], minimum, maximum)
	if err != nil {
		return 0, 0, err
	}
	if len(parts) == 1 {
		return start, start, nil
	}

	end, err := parseCronNumber(parts[1], minimum, maximum)
	if err != nil {
		return 0, 0, err
	}
	if end < start {
		return 0, 0, fmt.Errorf("invalid descending range %q", raw)
	}

	return start, end, nil
}

func parseCronNumber(raw string, minimum int, maximum int) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q", raw)
	}
	if value < minimum || value > maximum {
		return 0, fmt.Errorf("value %d outside %d-%d", value, minimum, maximum)
	}

	return value, nil
}

func normalizeCronValue(value int, weekday bool) int {
	if weekday && value == 7 {
		return 0
	}

	return value
}

func cronWeekday(value time.Time) int {
	return int(value.Weekday())
}

func setContains(values map[int]struct{}, value int) bool {
	_, ok := values[value]
	return ok
}
