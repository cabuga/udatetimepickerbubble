package datetimepicker

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

func updateWithKey(t *testing.T, model Model, key string) Model {
	t.Helper()

	var updated tea.Model
	switch key {
	case "up":
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	case "down":
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	default:
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
		updated, _ = model.Update(keyMsg)
	}

	pickerModel, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want datetimepicker.Model", updated)
	}

	return pickerModel
}
