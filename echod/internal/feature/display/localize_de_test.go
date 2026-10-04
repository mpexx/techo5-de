//go:build !dot

package display

import "testing"

func TestGermanDynamicScreenText(t *testing.T) {
	cases := map[string]string{
		"Monday, January 2":                  "Montag, 2. Januar",
		"Mon, Jan 2  ·  Snoozed until 08:30": "Mo., 2. Jan.  ·  Schlummern bis 08:30",
		"October 2026":                       "Oktober 2026",
		"High 21°   Low 12°":                 "Höchst 21°   Tiefst 12°",
		"High 21°  Low 12°  Rain 35%":        "Höchst 21°  Tiefst 12°  Regen 35%",
		"Rain 35%":                           "Regen 35%",
		"Listening…":                         "Ich höre zu…",
		"Living Room":                        "Living Room",
		"Clock style":                        "Uhrstil",
		"How the clock looks all day":        "Aussehen der Uhr am Tag",
		"Sun":                                "Sonne",
		"Sunday, October 4":                  "Sonntag, 4. Oktober",
		"Sun 4  ·  21°":                      "So. 4.  ·  21°",
	}
	for input, want := range cases {
		if got := germanScreenText(input); got != want {
			t.Errorf("%q: got %q, want %q", input, got, want)
		}
	}
	if got := germanWeekdayAbbrev("Sun"); got != "So." {
		t.Errorf("Sunday forecast heading: got %q, want So.", got)
	}
}
