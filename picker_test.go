package datetimepicker

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestFormatInitialExample(t *testing.T) {
	value := time.Date(2026, time.August, 20, 22, 21, 44, 0, time.UTC)

	got := Format(value)
	want := "CW 34.4  2026/08/20 22:21"

	if got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}
}

func TestDayAdjustUpdatesISOWeekday(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		InitialField: FieldDay,
	})

	updated := updateWithKey(t, model, "up")

	got := Format(updated.CurrentTime())
	want := "CW 34.5  2026/08/21 22:21"
	if got != want {
		t.Fatalf("after day up = %q, want %q", got, want)
	}
}

func TestCalendarWeekAdjustUpdatesDate(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		InitialField: FieldCalendarWeek,
	})

	updated := updateWithKey(t, model, "3")
	updated = updateWithKey(t, updated, "6")

	got := Format(updated.CurrentTime())
	want := "CW 36.4  2026/09/03 22:21"
	if got != want {
		t.Fatalf("after week 36 = %q, want %q", got, want)
	}
}

func TestMonthChangeClampsDay(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.March, 31, 10, 30, 0, 0, time.UTC),
		InitialField: FieldMonth,
	})

	updated := updateWithKey(t, model, "down")

	got := Format(updated.CurrentTime())
	want := "CW 09.6  2026/02/28 10:30"
	if got != want {
		t.Fatalf("after month down = %q, want %q", got, want)
	}
}

func TestDayDownMovesToPreviousMonth(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.March, 1, 10, 30, 0, 0, time.UTC),
		InitialField: FieldDay,
	})

	updated := updateWithKey(t, model, "down")

	got := Format(updated.CurrentTime())
	want := "CW 09.6  2026/02/28 10:30"
	if got != want {
		t.Fatalf("after day down = %q, want %q", got, want)
	}
}

func TestFormatDateOmitsTime(t *testing.T) {
	value := time.Date(2026, time.August, 20, 23, 21, 0, 0, time.UTC)

	got := FormatDate(value)
	want := "CW 34.4  2026/08/20"

	if got != want {
		t.Fatalf("FormatDate() = %q, want %q", got, want)
	}
}

func TestDateModeSkipsTimeFields(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 23, 21, 0, 0, time.UTC),
		InitialField: FieldDay,
		DateOnly:     true,
	})

	updated := updateWithKey(t, model, "right")

	if updated.SelectedField() != FieldCalendarWeek {
		t.Fatalf(
			"selected field = %v, want %v",
			updated.SelectedField(),
			FieldCalendarWeek,
		)
	}

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 00:00"
	if got != want {
		t.Fatalf("date mode current time = %q, want %q", got, want)
	}
}

func TestNextScheduleUsesDailySchedule(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(2026, time.August, 20, 20, 30, 0, 0, time.UTC),
		Schedules:   []Schedule{NewDailySchedule(21, 0)},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 21:00"
	if got != want {
		t.Fatalf("next schedule = %q, want %q", got, want)
	}
}

func TestNextScheduleUsesNearestSchedule(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(2026, time.August, 20, 20, 30, 0, 0, time.UTC),
		Schedules: []Schedule{
			NewDailySchedule(23, 0),
			NewDailySchedule(21, 0),
		},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 21:00"
	if got != want {
		t.Fatalf("nearest schedule = %q, want %q", got, want)
	}
}

func TestNextScheduleCanUseCustomKey(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(
			2026,
			time.August,
			20,
			21,
			0,
			0,
			0,
			time.UTC,
		),
		Schedules: []Schedule{NewDailySchedule(21, 0)},
		NextScheduleKeys: []string{
			"x",
		},
	})

	updated := updateWithKey(t, model, "x")

	got := Format(updated.CurrentTime())
	want := "CW 34.5  2026/08/21 21:00"
	if got != want {
		t.Fatalf("custom key next schedule = %q, want %q", got, want)
	}
}

func TestCronScheduleNextMonday(t *testing.T) {
	schedule, err := NewCronSchedule("5 4 * * 1")
	if err != nil {
		t.Fatalf("NewCronSchedule() error = %v", err)
	}

	model := New(Config{
		InitialTime: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
		Schedules: []Schedule{
			schedule,
		},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 35.1  2026/08/24 04:05"
	if got != want {
		t.Fatalf("cron next monday = %q, want %q", got, want)
	}
}

func TestOffsetSchedule(t *testing.T) {
	after := time.Date(2026, time.August, 20, 10, 30, 15, 0, time.UTC)
	for _, test := range []struct {
		expression string
		duration   time.Duration
	}{
		{expression: "offset:+72h", duration: 72 * time.Hour},
		{expression: "offset:72h", duration: 72 * time.Hour},
		{expression: "offset:72m", duration: 72 * time.Minute},
		{expression: "offset:2h", duration: 2 * time.Hour},
		{expression: "offset:3d", duration: 72 * time.Hour},
		{expression: "offset:3d2h", duration: 74 * time.Hour},
		{expression: "offset:2w", duration: 14 * 24 * time.Hour},
	} {
		schedule, err := NewOffsetSchedule(test.expression)
		if err != nil {
			t.Fatalf("NewOffsetSchedule(%q) error = %v", test.expression, err)
		}
		if got, want := schedule.Next(after), after.Add(test.duration); !got.Equal(want) {
			t.Errorf("Next(%q) = %q, want %q", test.expression, got, want)
		}
	}
}

func TestNewOffsetScheduleRejectsInvalidExpressions(t *testing.T) {
	for _, expression := range []string{"", "offset:", "offset:0s", "offset:-1h", "offset:bogus", "offset:3x", "72h"} {
		if _, err := NewOffsetSchedule(expression); err == nil {
			t.Errorf("NewOffsetSchedule(%q) error = nil, want error", expression)
		}
	}
}

func TestNewScheduleDispatchesOffsetAndCron(t *testing.T) {
	location := time.FixedZone("test", 3*60*60)
	after := time.Date(2026, time.August, 20, 10, 30, 15, 0, location)
	offset, err := NewSchedule("  offset:72h  ")
	if err != nil {
		t.Fatalf("NewSchedule(offset) error = %v", err)
	}
	if got, want := offset.Next(after), after.Add(72*time.Hour); !got.Equal(want) {
		t.Fatalf("offset Next() = %q, want %q", got, want)
	}
	if got := offset.Next(after).Location(); got != location {
		t.Fatalf("offset location = %v, want input location %v", got, location)
	}
	cron, err := NewSchedule("  5 4 * * 1  ")
	if err != nil {
		t.Fatalf("NewSchedule(cron) error = %v", err)
	}
	if got, want := cron.Next(after), time.Date(2026, time.August, 24, 4, 5, 0, 0, location); !got.Equal(want) {
		t.Fatalf("cron Next() = %q, want %q", got, want)
	}
}

func TestNextScheduleUsesNearestCronSchedule(t *testing.T) {
	monday, mondayErr := NewCronSchedule("5 4 * * 1")
	if mondayErr != nil {
		t.Fatalf("NewCronSchedule(monday) error = %v", mondayErr)
	}
	thursday, thursdayErr := NewCronSchedule("7 5 * * 4")
	if thursdayErr != nil {
		t.Fatalf("NewCronSchedule(thursday) error = %v", thursdayErr)
	}

	model := New(Config{
		InitialTime: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
		Schedules: []Schedule{
			monday,
			thursday,
		},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 05:07"
	if got != want {
		t.Fatalf("nearest cron = %q, want %q", got, want)
	}
}

func TestViewWidthLimitsRenderedLines(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(
			2026,
			time.August,
			20,
			22,
			21,
			0,
			0,
			time.UTC,
		),
		ShowTutorial: true,
		Title:        "Edit start datetime with a title that needs wrapping",
		Width:        66,
	})

	for lineNumber, line := range strings.Split(model.View().Content, "\n") {
		if width := ansi.StringWidth(line); width > 66 {
			t.Fatalf("line %d width = %d, want <= 66: %q", lineNumber+1, width, line)
		}
	}
}

func TestViewWrapsCompleteTutorialText(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		ShowTutorial: true,
		Width:        66,
		WrapTutorial: true,
	})

	view := model.View().Content
	if !strings.Contains(view, "q/esc to cancel.") {
		t.Fatalf("view does not contain end of wrapped tutorial: %q", view)
	}
	for lineNumber, line := range strings.Split(view, "\n") {
		if width := ansi.StringWidth(line); width > 66 {
			t.Fatalf("line %d width = %d, want <= 66: %q", lineNumber+1, width, line)
		}
	}
}

func TestViewUsesCompactTutorial(t *testing.T) {
	model := New(Config{
		InitialTime:     time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		ShowTutorial:    true,
		CompactTutorial: true,
		Width:           66,
	})

	view := model.View().Content
	if !strings.Contains(view, compactTutorialText) {
		t.Fatalf("view does not contain compact tutorial %q: %q", compactTutorialText, view)
	}
	if strings.Contains(view, defaultTutorialText) {
		t.Fatalf("view contains default tutorial when compact tutorial was requested: %q", view)
	}
}

func TestViewUsesCustomTutorialText(t *testing.T) {
	customTutorial := "left/right field | up/down adjust | enter save"
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		ShowTutorial: true,
		TutorialText: customTutorial,
		Width:        66,
	})

	view := model.View().Content
	if !strings.Contains(view, customTutorial) {
		t.Fatalf("view does not contain custom tutorial %q: %q", customTutorial, view)
	}
	if strings.Contains(view, defaultTutorialText) {
		t.Fatalf("view contains default tutorial when custom tutorial was set: %q", view)
	}
}

func TestViewClipsTutorialTextToWidthWithoutWrapping(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		ShowTutorial: true,
		Width:        66,
	})

	view := model.View().Content
	if strings.Contains(view, "q/esc to cancel.") {
		t.Fatalf("view contains end of tutorial without wrapping enabled: %q", view)
	}
	for lineNumber, line := range strings.Split(view, "\n") {
		if width := ansi.StringWidth(line); width > 66 {
			t.Fatalf("line %d width = %d, want <= 66: %q", lineNumber+1, width, line)
		}
	}
}

func TestViewWithoutWidthKeepsTutorialCompatible(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		ShowTutorial: true,
		Title:        "Select datetime",
	})

	got := model.View().Content
	want := "Select datetime\n" +
		defaultTutorialText +
		"\n\n" +
		"CW \x1b[7m34\x1b[0m.4  2026/08/20 22:21\n"

	if got != want {
		t.Fatalf("View() = %q, want %q", got, want)
	}
}

func updateWithKey(t *testing.T, model Model, key string) Model {
	t.Helper()

	var updated tea.Model
	switch key {
	case "left":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyLeft}))
	case "right":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyRight}))
	case "up":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyUp}))
	case "down":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	default:
		keyMsg := tea.KeyPressMsg(tea.Key{Code: []rune(key)[0], Text: key})
		updated, _ = model.Update(keyMsg)
	}

	pickerModel, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want datetimepicker.Model", updated)
	}

	return pickerModel
}
