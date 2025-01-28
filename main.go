package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"github.com/invopop/validation"
)

const (
	FULL_TIME    = 9*time.Hour + 48*time.Minute
	HALF_TIME    = 9*time.Hour + 30*time.Minute
	SHORT_TIME   = 9 * time.Hour
	MAX_OVERTIME = 11*time.Hour + 48*time.Minute
)

type ClockIn struct {
	StartHour   int
	StartMinute int
}

func (c ClockIn) GetTime() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		c.StartHour,
		c.StartMinute,
		0,
		0,
		time.Local,
	)
}

func (c ClockIn) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.StartHour, validation.NotNil.Error("não pode ser vazio")),
		validation.Field(&c.StartMinute, validation.NotNil.Error("não pode ser vazio")),
	)
}

func main() {
	log.SetReportTimestamp(false)

	clockIn := ClockIn{}

	if len(os.Args) < 2 {
		log.Fatal("Necessário informar horário de entrada no formato <h> <m>")
	}

	startHour, errHour := strconv.Atoi(os.Args[1])
	startMinute, errMinute := strconv.Atoi(os.Args[2])

	if errHour != nil || errMinute != nil {
		log.Fatal("Os horários de entrada e de saída são obrigatórios.")
	}

	clockIn.StartHour = startHour
	clockIn.StartMinute = startMinute

	if err := clockIn.Validate(); err != nil {
		log.Fatal(err.Error())
	}

	fullTime := clockIn.GetTime().Add(FULL_TIME)
	halfTime := clockIn.GetTime().Add(HALF_TIME)
	shortTime := clockIn.GetTime().Add(SHORT_TIME)
	maxOvertime := clockIn.GetTime().Add(MAX_OVERTIME)

	bold := lipgloss.NewStyle().Bold(true).Render
	center := lipgloss.NewStyle().Align(lipgloss.Center)

	t := table.New()
	t.StyleFunc(func(row, col int) lipgloss.Style {
		if col == 1 {
			return center
		}

		return table.DefaultStyles(row, col)
	})
	t.Row(bold("Jornada"), bold("Saída"))
	t.Row(bold("10:48 (+2h)"), maxOvertime.Format("15:04"))
	t.Row(bold("08:48"), fullTime.Format("15:04"))
	t.Row("08:30", halfTime.Format("15:04"))
	t.Row("08:00", shortTime.Format("15:04"))

	t.Border(lipgloss.HiddenBorder())

	fmt.Println(t.Render())
}
