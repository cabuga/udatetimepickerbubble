package datetimepicker

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	fieldCalendarWeekLength = 2
	fieldISOWeekdayLength   = 1
	fieldYearLength         = 4
	fieldMonthLength        = 2
	fieldDayLength          = 2
	fieldHourLength         = 2
	fieldMinuteLength       = 2
)

const tutorialText = "Use left/right to select, up/down or digits to edit, " +
	"t for now, enter to confirm, q/esc to cancel.\n\n"

type Field int

const (
	FieldCalendarWeek Field = iota
	FieldISOWeekday
	FieldYear
	FieldMonth
	FieldDay
	FieldHour
	FieldMinute
)

type ExitKey string

const (
	ExitKeyNone   ExitKey = ""
	ExitKeyEnter  ExitKey = "enter"
	ExitKeyQuit   ExitKey = "q"
	ExitKeyCancel ExitKey = "esc"
)

type Config struct {
	InitialTime  time.Time
	InitialField Field
	ShowTutorial bool
	Title        string
	Now          func() time.Time
}

type Result struct {
	Time     time.Time
	ExitKey  ExitKey
	Canceled bool
}

type Model struct {
	current      time.Time
	initialField Field
	field        Field
	tempInput    string
	showTutorial bool
	title        string
	now          func() time.Time
	result       Result
	done         bool
}

func New(config Config) Model {
	now := config.Now
	if now == nil {
		now = time.Now
	}

	current := config.InitialTime
	if current.IsZero() {
		current = now()
	}
	current = truncateToMinute(current)

	field := config.InitialField
	if field < FieldCalendarWeek || field > FieldMinute {
		field = FieldCalendarWeek
	}

	return Model{
		current:      current,
		initialField: field,
		field:        field,
		showTutorial: config.ShowTutorial,
		title:        config.Title,
		now:          now,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	key := keyMsg.String()
	switch key {
	case "left":
		m.applyTemp()
		m.field = previousField(m.field)
	case "right", "tab":
		m.applyTemp()
		m.field = nextField(m.field)
	case "up":
		m.applyTemp()
		m.adjust(1)
	case "down":
		m.applyTemp()
		m.adjust(-1)
	case "t":
		m.tempInput = ""
		m.current = truncateToMinute(m.now())
	case "enter":
		m.applyTemp()
		m.done = true
		m.result = Result{Time: m.current, ExitKey: ExitKeyEnter}
		return m, tea.Quit
	case "q":
		m.done = true
		m.result = Result{ExitKey: ExitKeyQuit, Canceled: true}
		return m, tea.Quit
	case "esc", "ctrl+c":
		m.done = true
		m.result = Result{ExitKey: ExitKeyCancel, Canceled: true}
		return m, tea.Quit
	default:
		if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
			m.tempInput += key
			if len(m.tempInput) >= fieldLength(m.field) {
				m.applyTemp()
				m.field = nextField(m.field)
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var builder strings.Builder

	if m.title != "" {
		builder.WriteString(m.title)
		builder.WriteByte('\n')
	}
	if m.showTutorial {
		builder.WriteString(tutorialText)
	}

	builder.WriteString(m.renderPicker())
	if m.tempInput != "" {
		builder.WriteString("  ")
		builder.WriteString(m.tempInput)
	}
	builder.WriteByte('\n')

	return builder.String()
}

func (m Model) CurrentTime() time.Time {
	return m.current
}

func (m Model) SelectedField() Field {
	return m.field
}

func (m Model) Result() Result {
	if m.done {
		return m.result
	}

	return Result{Time: m.current, ExitKey: ExitKeyNone}
}

func (m Model) Run() (Result, error) {
	program := tea.NewProgram(m)
	finalModel, err := program.Run()
	if err != nil {
		return Result{}, err
	}

	pickerModel, ok := finalModel.(Model)
	if !ok {
		return Result{}, unexpectedModelError(finalModel)
	}

	return pickerModel.Result(), nil
}

func unexpectedModelError(model tea.Model) error {
	return fmt.Errorf("unexpected final model type %T", model)
}

func (m Model) renderPicker() string {
	_, week := m.current.ISOWeek()
	isoDay := isoWeekday(m.current)
	values := []string{
		fmt.Sprintf("%02d", week),
		fmt.Sprintf("%d", isoDay),
		fmt.Sprintf("%04d", m.current.Year()),
		fmt.Sprintf("%02d", int(m.current.Month())),
		fmt.Sprintf("%02d", m.current.Day()),
		fmt.Sprintf("%02d", m.current.Hour()),
		fmt.Sprintf("%02d", m.current.Minute()),
	}

	for field := FieldCalendarWeek; field <= FieldMinute; field++ {
		if field == m.field {
			values[field] = "\x1b[7m" + values[field] + "\x1b[0m"
		}
	}

	return fmt.Sprintf(
		"CW %s.%s  %s/%s/%s %s:%s",
		values[FieldCalendarWeek],
		values[FieldISOWeekday],
		values[FieldYear],
		values[FieldMonth],
		values[FieldDay],
		values[FieldHour],
		values[FieldMinute],
	)
}

func (m *Model) applyTemp() {
	if m.tempInput == "" {
		return
	}

	var value int
	_, scanErr := fmt.Sscanf(m.tempInput, "%d", &value)
	if scanErr == nil {
		m.setField(value)
	}
	m.tempInput = ""
}

func (m *Model) adjust(delta int) {
	switch m.field {
	case FieldCalendarWeek:
		m.current = m.current.AddDate(0, 0, delta*7)
	case FieldISOWeekday:
		m.setISOWeekday(isoWeekday(m.current) + delta)
	case FieldYear:
		m.current = setYear(m.current, m.current.Year()+delta)
	case FieldMonth:
		m.current = setMonth(m.current, int(m.current.Month())+delta)
	case FieldDay:
		m.current = m.current.AddDate(0, 0, delta)
	case FieldHour:
		m.current = m.current.Add(time.Duration(delta) * time.Hour)
	case FieldMinute:
		m.current = m.current.Add(time.Duration(delta) * time.Minute)
	}
	m.current = truncateToMinute(m.current)
}

func (m *Model) setField(value int) {
	switch m.field {
	case FieldCalendarWeek:
		_, week := m.current.ISOWeek()
		m.current = m.current.AddDate(0, 0, (value-week)*7)
	case FieldISOWeekday:
		m.setISOWeekday(value)
	case FieldYear:
		m.current = setYear(m.current, clamp(value, 1, 9999))
	case FieldMonth:
		m.current = setMonth(m.current, clamp(value, 1, 12))
	case FieldDay:
		lastDay := daysInMonth(m.current.Year(), m.current.Month())
		m.current = setDay(m.current, clamp(value, 1, lastDay))
	case FieldHour:
		m.current = setClock(m.current, clamp(value, 0, 23), m.current.Minute())
	case FieldMinute:
		m.current = setClock(m.current, m.current.Hour(), clamp(value, 0, 59))
	}
	m.current = truncateToMinute(m.current)
}

func (m *Model) setISOWeekday(value int) {
	year, week := m.current.ISOWeek()
	day := clamp(value, 1, 7)
	date := isoWeekStart(year, week).AddDate(0, 0, day-1)
	m.current = time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		m.current.Hour(),
		m.current.Minute(),
		0,
		0,
		m.current.Location(),
	)
}

func Format(value time.Time) string {
	value = truncateToMinute(value)
	_, week := value.ISOWeek()
	return fmt.Sprintf(
		"CW %02d.%d  %04d/%02d/%02d %02d:%02d",
		week,
		isoWeekday(value),
		value.Year(),
		int(value.Month()),
		value.Day(),
		value.Hour(),
		value.Minute(),
	)
}

func previousField(field Field) Field {
	if field == FieldCalendarWeek {
		return FieldMinute
	}

	return field - 1
}

func nextField(field Field) Field {
	if field == FieldMinute {
		return FieldCalendarWeek
	}

	return field + 1
}

func fieldLength(field Field) int {
	switch field {
	case FieldCalendarWeek:
		return fieldCalendarWeekLength
	case FieldISOWeekday:
		return fieldISOWeekdayLength
	case FieldYear:
		return fieldYearLength
	case FieldMonth:
		return fieldMonthLength
	case FieldDay:
		return fieldDayLength
	case FieldHour:
		return fieldHourLength
	case FieldMinute:
		return fieldMinuteLength
	}

	return fieldMinuteLength
}

func isoWeekday(value time.Time) int {
	weekday := int(value.Weekday())
	if weekday == 0 {
		return 7
	}

	return weekday
}

func isoWeekStart(year int, week int) time.Time {
	janFourth := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)
	weekday := isoWeekday(janFourth)
	return janFourth.AddDate(0, 0, -(weekday-1)).AddDate(0, 0, (week-1)*7)
}

func setYear(value time.Time, year int) time.Time {
	day := min(value.Day(), daysInMonth(year, value.Month()))
	return time.Date(
		year,
		value.Month(),
		day,
		value.Hour(),
		value.Minute(),
		0,
		0,
		value.Location(),
	)
}

func setMonth(value time.Time, month int) time.Time {
	year := value.Year()
	for month < 1 {
		year--
		month += 12
	}
	for month > 12 {
		year++
		month -= 12
	}

	nextMonth := time.Month(month)
	day := min(value.Day(), daysInMonth(year, nextMonth))
	return time.Date(
		year,
		nextMonth,
		day,
		value.Hour(),
		value.Minute(),
		0,
		0,
		value.Location(),
	)
}

func setDay(value time.Time, day int) time.Time {
	return time.Date(
		value.Year(),
		value.Month(),
		day,
		value.Hour(),
		value.Minute(),
		0,
		0,
		value.Location(),
	)
}

func setClock(value time.Time, hour int, minute int) time.Time {
	return time.Date(
		value.Year(),
		value.Month(),
		value.Day(),
		hour,
		minute,
		0,
		0,
		value.Location(),
	)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func truncateToMinute(value time.Time) time.Time {
	return value.Truncate(time.Minute)
}

func clamp(value int, minimum int, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}

	return value
}
