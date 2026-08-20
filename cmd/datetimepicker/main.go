package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	datetimepicker "github.com/carlos/udatetimepickerbubble"
)

const defaultTimeLayout = "2006-01-02 15:04"
const timeUsage = "initial datetime in YYYY-MM-DD HH:MM format"

func main() {
	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	var initial string
	var fieldName string
	var showDoc bool
	var showAIGenerationDoc bool
	var showTutorial bool
	var title string

	flag.StringVar(&initial, "time", "", timeUsage)
	flag.StringVar(&fieldName, "field", "cw", "initial selected field")
	flag.BoolVar(
		&showAIGenerationDoc,
		"aigeneration",
		false,
		"print the AI usage document",
	)
	flag.BoolVar(&showDoc, "doc", false, "print the AI usage document")
	flag.BoolVar(&showTutorial, "tutorial", true, "show control hints")
	flag.StringVar(&title, "title", "Select datetime", "picker title")
	flag.Parse()

	if showAIGenerationDoc || showDoc {
		fmt.Print(datetimepicker.AIUsageDocument())
		return 0
	}

	initialTime, parseErr := parseInitialTime(initial)
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, parseErr)
		return 2
	}

	initialField, fieldErr := parseField(fieldName)
	if fieldErr != nil {
		fmt.Fprintln(os.Stderr, fieldErr)
		return 2
	}

	picker := datetimepicker.New(datetimepicker.Config{
		InitialTime:  initialTime,
		InitialField: initialField,
		ShowTutorial: showTutorial,
		Title:        title,
	})

	result, runErr := picker.Run()
	if runErr != nil {
		fmt.Fprintln(os.Stderr, runErr)
		return 1
	}
	if result.Canceled {
		return 130
	}

	fmt.Println(datetimepicker.Format(result.Time))
	return 0
}

func parseInitialTime(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}

	value, err := time.ParseInLocation(defaultTimeLayout, raw, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse -time: %w", err)
	}

	return value, nil
}

func parseField(raw string) (datetimepicker.Field, error) {
	switch strings.ToLower(raw) {
	case "cw", "calendar-week", "week":
		return datetimepicker.FieldCalendarWeek, nil
	case "cw-day", "iso-weekday", "weekday":
		return datetimepicker.FieldISOWeekday, nil
	case "year":
		return datetimepicker.FieldYear, nil
	case "month":
		return datetimepicker.FieldMonth, nil
	case "day":
		return datetimepicker.FieldDay, nil
	case "hour":
		return datetimepicker.FieldHour, nil
	case "minute":
		return datetimepicker.FieldMinute, nil
	}

	return datetimepicker.FieldCalendarWeek, fmt.Errorf("unknown field %q", raw)
}
