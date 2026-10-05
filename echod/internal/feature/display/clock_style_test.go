//go:build !dot && !spot

package display

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HuskerMinion/techo5/echod/internal/config"
	"github.com/HuskerMinion/techo5/echod/internal/feature/alarm"
	"github.com/HuskerMinion/techo5/echod/internal/feature/home"
	"github.com/HuskerMinion/techo5/echod/internal/feature/timer"
	"github.com/HuskerMinion/techo5/echod/internal/lib/hass"
)

// Every clock style draws on both panels, alone, with a timer running, with the glance strip and
// with the music strip; with SHOW_PREVIEW set, each is written there to look at.
func TestShowClockStylesDraw(t *testing.T) {
	config.Use(filepath.Join(t.TempDir(), "state.json"))
	t.Cleanup(func() { _ = config.Set().Screen().ClockStyle("") })
	at := time.Date(2026, 9, 16, 14, 7, 38, 0, time.Local)
	sky := home.Weather{Condition: "partlycloudy", Temp: "72°"}
	var week []hass.Day
	for i, c := range []string{"partlycloudy", "rainy", "sunny", "cloudy", "snowy"} {
		week = append(week, hass.Day{When: at.AddDate(0, 0, i), Condition: c, High: float64(72 - 3*i), Low: float64(51 - 2*i)})
	}
	facts := styleFacts{
		rise: time.Date(2026, 9, 16, 6, 48, 0, 0, time.Local), set: time.Date(2026, 9, 16, 19, 5, 0, 0, time.Local), sunOK: true,
		days: week,
		next: []hass.Event{
			{Summary: "Dentist", Start: at.Add(83 * time.Minute), End: at.Add(143 * time.Minute)},
			{Summary: "Soccer practice", Start: at.Add(233 * time.Minute), End: at.Add(293 * time.Minute)},
			{Summary: "Trash day", Start: time.Date(2026, 9, 17, 0, 0, 0, 0, time.Local), End: time.Date(2026, 9, 18, 0, 0, 0, 0, time.Local), AllDay: true},
		},
	}
	later := facts
	later.next = append([]hass.Event{{Summary: "Doctor's appointment", Start: time.Date(2026, 9, 17, 9, 45, 0, 0, time.Local), End: time.Date(2026, 9, 17, 10, 45, 0, 0, time.Local)},
		{Summary: "Still going", Start: at.Add(-30 * time.Minute), End: at.Add(time.Hour)}}, facts.next[:1]...)
	running := []timer.Countdown{{Name: "Pasta", Left: 4*time.Minute + 32*time.Second, Total: 10 * time.Minute, Active: true}}
	chips := []home.Chip{{Icon: "mdi:door", Text: "Back door open"}, {Icon: "mdi:thermometer", Text: "Upstairs 74°"}}
	scenes := map[string]scene{
		"":        {now: at, phase: "idle", weather: sky, style: facts},
		"-timer":  {now: at, phase: "idle", weather: sky, style: facts, timers: running},
		"-glance": {now: at, phase: "idle", weather: sky, style: facts, glance: chips},
		"-strip": {now: at, phase: "idle", weather: sky, style: facts, strip: true, playing: true,
			radio: home.Radio{Now: "KXYZ 101.1", Title: "Take It Easy", Artist: "Eagles"}},
		"-dawn": {now: time.Date(2026, 9, 16, 5, 30, 0, 0, time.Local), phase: "idle", weather: sky, style: facts},
		"-half": {now: time.Date(2026, 9, 16, 14, 30, 0, 0, time.Local), phase: "idle", weather: sky, style: facts},
		"-long": {now: time.Date(2026, 9, 16, 14, 31, 0, 0, time.Local), phase: "idle", weather: sky, style: facts},
		// The next alarm makes the date line longer; an event tomorrow at a time makes the when longer.
		"-alarm": {now: at, phase: "idle", weather: sky, style: later, alarms: alarm.View{Next: &alarm.Upcoming{At: at.Add(16 * time.Hour)}}},
		"-strip-timer": {now: at, phase: "idle", weather: sky, style: facts, timers: running, strip: true, playing: true,
			radio: home.Radio{Now: "KXYZ 101.1", Title: "Take It Easy", Artist: "Eagles"}},
	}
	dir := os.Getenv("SHOW_PREVIEW")
	for _, st := range clockStyles {
		if err := config.Set().Screen().ClockStyle(st.value); err != nil {
			t.Fatal(err)
		}
		for suffix, s := range scenes {
			for _, panel := range []struct {
				name       string
				wide, high int
			}{{"", showWide, showHigh}, {"-show8", show8Wide, show8High}} {
				img := image.NewRGBA(image.Rect(0, 0, panel.wide, panel.high))
				newRenderer(img).draw(s)
				if dir == "" {
					continue
				}
				name := "style-" + st.label + suffix + panel.name + ".png"
				f, err := os.Create(filepath.Join(dir, name))
				if err != nil {
					t.Fatal(err)
				}
				if err := png.Encode(f, img); err != nil {
					t.Fatal(err)
				}
				f.Close()
			}
		}
	}
}

func TestClockWords(t *testing.T) {
	for _, c := range []struct {
		h, m               int
		lead, hour, period string
	}{
		{14, 7, "sieben nach", "zwei", "nachmittags"},
		{14, 0, "", "zwei Uhr", "nachmittags"},
		{14, 15, "Viertel nach", "zwei", "nachmittags"},
		{14, 30, "halb", "drei", "nachmittags"},
		{14, 45, "Viertel vor", "drei", "nachmittags"},
		{14, 37, "dreiundzwanzig vor", "drei", "nachmittags"},
		{9, 1, "eine Minute nach", "neun", "morgens"},
		{11, 59, "eine Minute vor", "zwölf", "vormittags"},
		{12, 0, "", "Mittag", ""},
		{0, 0, "", "Mitternacht", ""},
		{23, 50, "zehn vor", "zwölf", "nachts"},
		{19, 20, "zwanzig nach", "sieben", "abends"},
		{4, 45, "Viertel vor", "fünf", "morgens"},
		{17, 45, "Viertel vor", "sechs", "abends"},
		{13, 0, "", "ein Uhr", "mittags"},
		{1, 0, "", "ein Uhr", "nachts"},
		{0, 30, "halb", "eins", "nachts"},
		{10, 21, "einundzwanzig nach", "zehn", "vormittags"},
		{14, 31, "neunundzwanzig vor", "drei", "nachmittags"},
	} {
		lead, hour, period := clockWords(time.Date(2026, 9, 16, c.h, c.m, 0, 0, time.Local))
		if lead != c.lead || hour != c.hour || period != c.period {
			t.Errorf("%02d:%02d = %q %q %q, want %q %q %q", c.h, c.m, lead, hour, period, c.lead, c.hour, c.period)
		}
	}
}
