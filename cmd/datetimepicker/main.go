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
const defaultDateLayout = "2006-01-02"
const timeUsage = "initial datetime in YYYY-MM-DD HH:MM format"
const cronUsage = "cron schedule as 'M H DOM MON DOW'"

func main() {
	exitCode := run()
	os.Exit(exitCode)
}

type cliConfig struct {
	initial             string
	fieldName           string
	modeName            string
	dailySchedules      dailyFlag
	cronSchedules       dailyFlag
	showDoc             bool
	showAIGenerationDoc bool
	showTutorial        bool
	title               string
}

func run() int {
	config := cliConfig{}
	registerFlags(&config)
	flag.Parse()

	if config.showAIGenerationDoc || config.showDoc {
		fmt.Print(datetimepicker.AIUsageDocument())
		return 0
	}

	dateOnly, modeErr := parseMode(config.modeName)
	if modeErr != nil {
		fmt.Fprintln(os.Stderr, modeErr)
		return 2
	}

	initialTime, parseErr := parseInitialTime(config.initial, dateOnly)
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, parseErr)
		return 2
	}

	initialField, fieldErr := parseField(config.fieldName)
	if fieldErr != nil {
		fmt.Fprintln(os.Stderr, fieldErr)
		return 2
	}

	schedules, scheduleErr := parseSchedules(config)
	if scheduleErr != nil {
		fmt.Fprintln(os.Stderr, scheduleErr)
		return 2
	}

	picker := datetimepicker.New(datetimepicker.Config{
		InitialTime:  initialTime,
		InitialField: initialField,
		DateOnly:     dateOnly,
		Schedules:    schedules,
		ShowTutorial: config.showTutorial,
		Title:        config.title,
	})

	result, runErr := picker.Run()
	if runErr != nil {
		fmt.Fprintln(os.Stderr, runErr)
		return 1
	}
	if result.Canceled {
		return 130
	}

	if dateOnly {
		fmt.Println(datetimepicker.FormatDate(result.Time))
		return 0
	}

	fmt.Println(datetimepicker.Format(result.Time))
	return 0
}

func registerFlags(config *cliConfig) {
	flag.StringVar(&config.initial, "time", "", timeUsage)
	flag.StringVar(&config.fieldName, "field", "cw", "initial selected field")
	flag.StringVar(
		&config.modeName,
		"mode",
		"datetime",
		"picker mode: datetime or date",
	)
	flag.Var(&config.dailySchedules, "daily", "daily schedule in HH:MM format")
	flag.Var(&config.cronSchedules, "cron", cronUsage)
	flag.BoolVar(
		&config.showAIGenerationDoc,
		"aigeneration",
		false,
		"print the AI usage document",
	)
	flag.BoolVar(&config.showDoc, "doc", false, "print the AI usage document")
	flag.BoolVar(&config.showTutorial, "tutorial", true, "show control hints")
	flag.StringVar(&config.title, "title", "Select datetime", "picker title")
}

func parseInitialTime(
	raw string,
	dateOnly bool,
) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}

	if dateOnly {
		value, err := time.ParseInLocation(defaultDateLayout, raw, time.Local)
		if err == nil {
			return value, nil
		}
	}

	value, err := time.ParseInLocation(defaultTimeLayout, raw, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse -time: %w", err)
	}

	return value, nil
}

func parseMode(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "datetime", "date-time":
		return false, nil
	case "date":
		return true, nil
	}

	return false, fmt.Errorf("unknown mode %q", raw)
}

func parseSchedules(config cliConfig) ([]datetimepicker.Schedule, error) {
	return parseScheduleValues(config.dailySchedules, config.cronSchedules)
}

func parseScheduleValues(
	dailyValues []string,
	cronValues []string,
) ([]datetimepicker.Schedule, error) {
	size := len(dailyValues) + len(cronValues)
	schedules := make([]datetimepicker.Schedule, 0, size)
	for _, value := range dailyValues {
		parsed, err := time.Parse("15:04", value)
		if err != nil {
			return nil, fmt.Errorf("parse -daily %q: %w", value, err)
		}

		schedules = append(
			schedules,
			datetimepicker.NewDailySchedule(parsed.Hour(), parsed.Minute()),
		)
	}
	for _, value := range cronValues {
		schedule, err := datetimepicker.NewCronSchedule(value)
		if err != nil {
			return nil, fmt.Errorf("parse -cron %q: %w", value, err)
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
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

type dailyFlag []string

func (f *dailyFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *dailyFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}
