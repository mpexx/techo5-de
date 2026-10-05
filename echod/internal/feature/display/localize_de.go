//go:build !dot

package display

import (
	"strconv"
	"strings"
)

// germanScreenText is a small, display-only German overlay. It changes text as it
// is drawn, so setting IDs, picker choices, Home Assistant entities and saved
// configuration retain their original values. Unmapped text passes through
// unchanged. Add new upstream screen labels here when rebasing.
func germanScreenText(s string) string {
	if translated, ok := germanScreenLabels[s]; ok {
		return translated
	}
	if translated, ok := germanDateText(s); ok {
		return translated
	}
	if strings.HasPrefix(s, "Snoozed until ") {
		return "Schlummern bis " + strings.TrimPrefix(s, "Snoozed until ")
	}
	if strings.HasPrefix(s, "High ") && strings.Contains(s, "°") && strings.Contains(s, "  Low ") {
		s = strings.Replace(s, "High ", "Höchst ", 1)
		s = strings.Replace(s, "Low ", "Tiefst ", 1)
		s = strings.Replace(s, "Rain ", "Regen ", 1)
		return s
	}
	if strings.HasPrefix(s, "Rain ") && strings.HasSuffix(s, "%") {
		return "Regen " + strings.TrimPrefix(s, "Rain ")
	}
	return s
}

// Dates arrive as Go's English time.Format output. Translate only known date
// shapes, leaving Home Assistant names and spoken replies untouched.
func germanDateText(s string) (string, bool) {
	date, suffix, _ := strings.Cut(s, "  ·  ")
	if suffix != "" {
		suffix = "  ·  " + germanScreenText(suffix)
	}
	if day, rest, ok := strings.Cut(date, ", "); ok {
		if deDay, found := germanWeekdays[day]; found {
			if translated, valid := germanMonthDay(rest); valid {
				return deDay + ", " + translated + suffix, true
			}
		}
	}
	if parts := strings.Fields(date); len(parts) == 2 {
		if day, found := germanWeekdays[parts[0]]; found && len(parts[0]) == 3 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n >= 1 && n <= 31 {
				return day + " " + strconv.Itoa(n) + "." + suffix, true
			}
		}
	}
	if translated, ok := germanWeekdays[date]; ok && len(date) > 3 {
		return translated + suffix, true
	}
	if translated, ok := germanMonthDay(date); ok {
		return translated + suffix, true
	}
	return "", false
}

func germanMonthDay(s string) (string, bool) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return "", false
	}
	month, ok := germanMonths[parts[0]]
	if !ok {
		return "", false
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", false
	}
	if n >= 1 && n <= 31 {
		return strconv.Itoa(n) + ". " + month, true
	}
	if n >= 2000 && n <= 9999 {
		return month + " " + parts[1], true
	}
	return "", false
}

var germanWeekdays = map[string]string{
	"Monday": "Montag", "Tuesday": "Dienstag", "Wednesday": "Mittwoch",
	"Thursday": "Donnerstag", "Friday": "Freitag", "Saturday": "Samstag", "Sunday": "Sonntag",
	"Mon": "Mo.", "Tue": "Di.", "Wed": "Mi.", "Thu": "Do.",
	"Fri": "Fr.", "Sat": "Sa.", "Sun": "So.",
}

// Bare weekday abbreviations are only translated where the caller knows the
// text is a date. "Sun" in the clock-style picker means the sun.
func germanWeekdayAbbrev(s string) string {
	if translated, ok := germanWeekdays[s]; ok {
		return translated
	}
	return s
}

var germanMonths = map[string]string{
	"January": "Januar", "February": "Februar", "March": "März", "April": "April",
	"May": "Mai", "June": "Juni", "July": "Juli", "August": "August",
	"September": "September", "October": "Oktober", "November": "November", "December": "Dezember",
	"Jan": "Jan.", "Feb": "Feb.", "Mar": "Mär.", "Apr": "Apr.",
	"Jun": "Jun.", "Jul": "Jul.", "Aug": "Aug.", "Sep": "Sep.",
	"Oct": "Okt.", "Nov": "Nov.", "Dec": "Dez.",
}

var germanScreenLabels = map[string]string{
	// Navigation and category headings.
	"Settings":           "Einstellungen",
	"Display":            "Anzeige",
	"Sound":              "Ton",
	"Alarms":             "Wecker",
	"Connections":        "Netzwerk",
	"Privacy":            "Privat",
	"General":            "Allgemein",
	"Sound & Voice":      "Ton & Sprache",
	"Alarms & Timers":    "Wecker & Timer",
	"Privacy & Security": "Datenschutz",
	"Brightness, night, theme, home screen and weather": "Helligkeit, Nacht, Design und Wetter",
	"Brightness, night, clock and photos":               "Helligkeit, Nacht, Uhr und Fotos",
	"Volume, voice, quiet hours and music":              "Lautstärke, Sprache, Ruhezeit und Musik",
	"Alarms on this device and how they ring":           "Wecker und ihre Klingeltöne",
	"Wi-Fi, Bluetooth and the Bluetooth proxy":          "WLAN, Bluetooth und Bluetooth-Proxy",
	"Remote access and how this device is reached":      "Fernzugriff und Gerätesicherheit",
	"Remote access and security":                        "Fernzugriff und Sicherheit",
	"Name, weather, updates and restart":                "Name, Wetter, Updates und Neustart",

	// Display card.
	"Brightness":               "Helligkeit",
	"Auto-brightness":          "Automatische Helligkeit",
	"Follows the room's light": "Passt sich dem Raumlicht an",
	"Dimmest":                  "Mindesthelligkeit",
	"How dark auto-brightness goes in a dark room": "Helligkeit bei Dunkelheit",
	"Night hours": "Nachtzeiten",
	"Night mode":  "Nachtmodus",
	"At night":    "Nachts",
	"Dark, a faint glow, or a clock alone until touched": "Dunkel, Nachtlicht oder Uhr",
	"Night clock":                          "Nachtuhr",
	"How the night clock looks":            "Darstellung der Nachtuhr",
	"Look":                                 "Aussehen",
	"Theme":                                "Design",
	"Custom colors":                        "Eigene Farben",
	"Make the theme your own":              "Farben selbst festlegen",
	"Clock format":                         "Uhrzeitformat",
	"Home screen":                          "Startbildschirm",
	"Slideshow":                            "Diashow",
	"Photos from Home Assistant":           "Fotos aus Home Assistant",
	"Now playing":                          "Aktuelle Wiedergabe",
	"Full page, or a strip over the clock": "Vollbild oder Leiste über der Uhr",
	"Call button":                          "Anruftaste",
	"On the home screen: devices and contacts": "Geräte und Kontakte auf der Startseite",
	"Weather":           "Wetter",
	"Weather animation": "Wetteranimation",
	"Rain, snow and storms move on the forecast": "Regen, Schnee und Gewitter animieren",
	"Weather alerts":                               "Wetterwarnungen",
	"The NWS's alerts for home, in the U.S.":       "NWS-Warnungen in den USA",
	"Radar source":                                 "Radarquelle",
	"Automatic uses the NWS in the lower 48":       "Automatisch: NWS in den 48 US-Staaten",
	"Pop-ups":                                      "Einblendungen",
	"Camera time":                                  "Kameradauer",
	"How long a camera opened here stays up":       "Anzeigedauer einer geöffneten Kamera",
	"Answer time":                                  "Antwortdauer",
	"How long an answer stays up; a tap clears it": "Anzeigedauer einer Antwort",
	"Turn screen":                                  "Sprachanzeige",
	"Classic, or a wave or bars that move with the voice": "Klassisch, Welle oder Balken",
	"Clock style":                        "Uhrstil",
	"How the clock looks all day":        "Aussehen der Uhr am Tag",
	"Clock position":                     "Uhrposition",
	"Out of a photo's way":               "Damit Fotos sichtbar bleiben",
	"Date color":                         "Datumsfarbe",
	"The date and AM/PM":                 "Für Datum und AM/PM",
	"Classic":                            "Klassisch",
	"Big":                                "Groß",
	"Flip":                               "Klappuhr",
	"LED":                                "LED",
	"Analog":                             "Analog",
	"Words":                              "Wortuhr",
	"Sun":                                "Sonne",
	"Dashboard":                          "Übersicht",
	"Center":                             "Mitte",
	"Bottom left":                        "Unten links",
	"Bottom right":                       "Unten rechts",
	"Default":                            "Standard",
	"White":                              "Weiß",
	"Gold":                               "Gold",
	"Sky":                                "Himmelblau",
	"Mint":                               "Mintgrün",
	"Rose":                               "Rosa",
	"Setup page":                         "Einrichtungsseite",
	"Photo folder":                       "Fotoordner",
	"Time per photo":                     "Zeit pro Foto",
	"Shuffle photos":                     "Fotos mischen",
	"Include subfolders":                 "Unterordner einbeziehen",
	"Every folder inside the one chosen": "Alle Unterordner des gewählten Ordners",
	"Show whole photo":                   "Ganzes Foto anzeigen",
	"All of it, with blurred sides; off fills the screen": "Ganzes Foto mit unscharfen Rändern; aus füllt den Bildschirm",
	"Weather art": "Wetterbild",
	"A landscape for the weather, in place of photos": "Wetterlandschaft statt Fotos",
	"Calendars": "Kalender",
	"Shown on the calendar page, from Home Assistant": "Kalender aus Home Assistant anzeigen",
	"Event pop-ups":                              "Terminhinweise",
	"An event coming up, over the screen":        "Bevorstehende Termine einblenden",
	"Pop up":                                     "Einblenden",
	"Pop-up chime":                               "Hinweiston",
	"Silent at night and in quiet hours":         "Nachts und während der Ruhezeit stumm",
	"All-day events":                             "Ganztägige Termine",
	"Pop-up calendars":                           "Kalender für Hinweise",
	"Speaking voice":                             "Sprechstimme",
	"How answers sound":                          "Stimme für Antworten",
	"Home Assistant: Settings, Voice assistants": "Home Assistant: Einstellungen, Sprachassistenten",
	"Set there":                                  "Dort einstellen",
	"AirPlay":                                    "AirPlay",
	"Play to it from an iPhone, iPad or Mac":     "Vom iPhone, iPad oder Mac abspielen",
	"Spotify Connect":                            "Spotify Connect",
	"Play to it from the Spotify app (Premium)":  "Aus der Spotify-App abspielen (Premium)",
	"Wi-Fi":                                "WLAN",
	"SSH":                                  "SSH",
	"No alarms yet":                        "Noch keine Wecker",
	"Add one here, or from Home Assistant": "Hier oder in Home Assistant hinzufügen",
	"+ Add alarm":                          "+ Wecker hinzufügen",
	"Delete alarm":                         "Wecker löschen",
	"Delete":                               "Löschen",
	"Snooze length":                        "Schlummerdauer",
	"Ring volume":                          "Weckerlautstärke",
	"Alarms and timers will make no sound": "Wecker und Timer bleiben stumm",
	"Alarms and timers, not the music":     "Für Wecker und Timer, nicht für Musik",
	"Plays once when you choose it":        "Bei Auswahl einmal abspielen",
	"Sun with a face":                      "Sonne mit Gesicht",
	"A face on it, for whoever has to look at it": "Sonne mit Gesicht beim Lichtwecker",
	"Rings from this device itself":               "Klingelt auf diesem Gerät",
	"Asks twice":                                  "Fragt zweimal",
	"Days":                                        "Tage",
	"Hour":                                        "Stunde",
	"Minute":                                      "Minute",
	"Set":                                         "Übernehmen",
	"Check for updates":                           "Nach Updates suchen",
	"Say something":                               "Etwas sagen",
	"It sends when you stop talking":              "Wird gesendet, wenn du aufhörst zu sprechen",
	"Speak":                                       "Sprechen",
	"Speak now":                                   "Jetzt sprechen",
	"Asking Home Assistant for stations…":         "Frage Sender bei Home Assistant ab…",
	"Full page":                                   "Ganze Seite",
	"Strip after 10 s":                            "Leiste nach 10 s",
	"Strip after 30 s":                            "Leiste nach 30 s",
	"Strip at once":                               "Leiste sofort",
	"15 seconds":                                  "15 Sekunden",
	"30 seconds":                                  "30 Sekunden",
	"1 minute":                                    "1 Minute",
	"2 minutes":                                   "2 Minuten",
	"3 minutes":                                   "3 Minuten",
	"5 minutes":                                   "5 Minuten",
	"10 minutes":                                  "10 Minuten",
	"15 minutes":                                  "15 Minuten",
	"1 hour":                                      "1 Stunde",
	"Once":                                        "Einmal",
	"Every day":                                   "Jeden Tag",
	"Weekdays":                                    "Wochentags",
	"Weekends":                                    "Wochenends",
	"As it starts":                                "Bei Beginn",
	"5 minutes before":                            "5 Minuten vorher",
	"10 minutes before":                           "10 Minuten vorher",
	"15 minutes before":                           "15 Minuten vorher",
	"30 minutes before":                           "30 Minuten vorher",
	"1 hour before":                               "1 Stunde vorher",
	"In the morning":                              "Morgens",
	"Never":                                       "Nie",
	"All calendars shown":                         "Alle angezeigten Kalender",
	"None in Home Assistant yet":                  "Noch keine Kalender in Home Assistant",
	"None chosen":                                 "Keiner ausgewählt",
	"On this device":                              "Auf diesem Gerät",
	"Opening…":                                    "Öffne…",
	"Couldn't open this folder":                   "Ordner konnte nicht geöffnet werden",
	"‹ Back":                                      "‹ Zurück",
	"‹  Back":                                     "‹  Zurück",
	"Running":                                     "Läuft",
	"Silent":                                      "Stumm",
	"Tap a color for each part; the theme becomes Custom": "Farben wählen; Design wird benutzerdefiniert",
	"5 seconds":    "5 Sekunden",
	"10 seconds":   "10 Sekunden",
	"20 seconds":   "20 Sekunden",
	"Until tapped": "Bis zum Antippen",
	"Wave":         "Welle",
	"Bars":         "Balken",
	"Screen off":   "Bildschirm aus",
	"Night light":  "Nachtlicht",
	"Red":          "Rot",
	"Red LED":      "Rote LED",
	"Custom":       "Benutzerdefiniert",
	"Ground":       "Hintergrund",
	"Accent":       "Akzent",
	"Text":         "Text",
	"Dim text":     "Gedämpfter Text",
	"Rules":        "Trennlinien",
	"Walnut":       "Walnuss",
	"Slate":        "Schiefer",
	"Midnight":     "Mitternacht",
	"Forest":       "Wald",
	"Plum":         "Pflaume",
	"Ocean":        "Ozean",
	"Ember":        "Glut",
	"Lavender":     "Lavendel",
	"Graphite":     "Graphit",
	"Cherry":       "Kirsche",
	"Paper":        "Papier",
	"Linen":        "Leinen",

	// Sound card.
	"Volume":       "Lautstärke",
	"Bass":         "Bass",
	"Treble":       "Höhen",
	"Audio output": "Audioausgabe",
	"Where the sound goes with headphones in": "Ausgabe bei angeschlossenen Kopfhörern",
	"Voice":                              "Sprache",
	"Microphone":                         "Mikrofon",
	"The mute button does this too":      "Auch über die Stummschalttaste",
	"Listening":                          "Aktiv",
	"Muted":                              "Stumm",
	"Wake word":                          "Aktivierungswort",
	"Wake word sensitivity":              "Wort-Empfindlichkeit",
	"Higher wakes by mistake less often": "Höher: weniger Fehlauslösungen",
	"Wake sound":                         "Aktivierungston",
	"Home Assistant sounds":              "Home-Assistant-Töne",
	"For muting and timers":              "Für Stummschalten und Timer",
	"Quiet":                              "Ruhe",
	"Quiet hours":                        "Ruhezeiten",
	"Sleep timer":                        "Einschlaftimer",
	"Do not disturb":                     "Nicht stören",
	"Intercom calls from other rooms are turned away": "Rufe aus anderen Räumen abweisen",
	"Music":                               "Musik",
	"Music Assistant player":              "Music-Assistant-Player",
	"Play music in sync with other rooms": "Musik synchron in mehreren Räumen",
	"Cameras":                             "Kameras",
	"Camera sound":                        "Kameraton",
	"A camera's own audio, while its view is up": "Ton der gerade geöffneten Kamera",

	// Network, privacy and general cards.
	"Bluetooth audio":                         "Bluetooth-Audio",
	"Bluetooth proxy":                         "Bluetooth-Proxy",
	"Lets Home Assistant hear nearby devices": "Geräte in der Nähe an Home Assistant",
	"Not available on this build":             "In diesem Build nicht verfügbar",
	"Not connected":                           "Nicht verbunden",
	"None yet":                                "Noch keines",
	"Earbuds or a speaker":                    "Kopfhörer oder Lautsprecher",
	"Pair a new device":                       "Neues Gerät koppeln",
	"Put it in pairing mode first":            "Zuerst Kopplungsmodus einschalten",
	"No address yet":                          "Noch keine Adresse",
	"No login":                                "Ohne Anmeldung",
	"SSH keys":                                "SSH-Schlüssel",
	"Sent from Home Assistant":                "Von Home Assistant gesendet",
	"Closed":                                  "Geschlossen",
	"Port 22, keys only":                      "Port 22, nur Schlüssel",
	"Not managed here":                        "Hier nicht verwaltet",
	"Allow Drop In":                           "Drop In erlauben",
	"Intercom calls connect by themselves, after a chime": "Hausrufe nach Signalton annehmen",
	"Home Assistant link":                             "Home-Assistant-Verbindung",
	"Encrypted with this device's key":                "Mit Geräteschlüssel verschlüsselt",
	"Encrypted":                                       "Verschlüsselt",
	"Not encrypted":                                   "Nicht verschlüsselt",
	"Certificate checks":                              "Zertifikatsprüfung",
	"For downloads; updates always check":             "Für Downloads; Updates werden stets geprüft",
	"Turned off in Home Assistant":                    "In Home Assistant ausgeschaltet",
	"Camera on the network":                           "Kamera im Netzwerk",
	"Screen on the network":                           "Bildschirm im Netzwerk",
	"Talk through cameras":                            "Über Kameras sprechen",
	"Talk on the camera page":                         "Sprechen auf der Kameraseite",
	"Name":                                            "Name",
	"Change it on the setup page":                     "Auf der Einrichtungsseite ändern",
	"Shown with the clock":                            "Zusammen mit der Uhr anzeigen",
	"Time zone":                                       "Zeitzone",
	"Screen language":                                 "Bildschirmsprache",
	"Its dates, its weather, and what it listens for": "Datums- und Wettertexte sowie Sprachbefehle",
	"What this screen listens for, not what the assistant speaks": "Nur Sprachbefehle am Bildschirm",
	"Updates":                  "Updates",
	"Checking…":                "Prüfe…",
	"Restart":                  "Neustart",
	"Back in about a minute":   "In etwa einer Minute zurück",
	"Tap again to restart now": "Zum Neustart erneut tippen",
	"About":                    "Info",

	// Settings and pages added in upstream v1.0.0.
	"Play to it from music apps and servers":           "Wiedergabe über Musik-Apps und Server",
	"Now playing follows":                              "Wiedergabe folgen",
	"Another speaker's music, while this one is quiet": "Musik eines anderen Lautsprechers, wenn dieser still ist",
	"Lyrics":                                 "Liedtext",
	"The words in time, looked up at LRCLIB": "Synchroner Liedtext von LRCLIB",
	"Tap on the clock":                       "Uhr antippen",
	"Start Assist, open the dashboard, or nothing": "Sprachassistent, Übersicht oder keine Aktion",
	"Assist":                           "Sprachassistent",
	"Nothing":                          "Keine Aktion",
	"Settings lock":                    "Einstellungssperre",
	"A PIN before these settings open": "PIN zum Öffnen der Einstellungen",
	"On: turn off to remove the PIN":   "Zum Entfernen der PIN ausschalten",
	"Presence":                         "Anwesenheit",
	"Presence detection":               "Anwesenheitserkennung",
	"The camera notices somebody near; nothing is kept": "Kamera erkennt Personen in der Nähe; keine Aufzeichnung",
	"Screen off when nobody is near":                    "Bildschirm ohne Personen ausschalten",
	"Lights again when somebody comes near":             "Schaltet sich bei Annäherung wieder ein",
	"30 minutes":                                        "30 Minuten",
	"60 minutes":                                        "60 Minuten",
	"Enter the PIN":                                     "PIN eingeben",
	"A new PIN, 4 to 8 digits":                          "Neue PIN mit 4 bis 8 Ziffern",
	"The same PIN again":                                "PIN erneut eingeben",
	"Wrong PIN":                                         "Falsche PIN",
	"They did not match":                                "PINs stimmen nicht überein",
	"Could not save it":                                 "Speichern fehlgeschlagen",
	"Subtle mute ring":                                  "Dezenter Stumm-Ring",
	"Now":                                               "Jetzt",

	// Common actions and picker choices. Values are translated only for drawing.
	"Back":            "Zurück",
	"Done":            "Fertig",
	"Edit":            "Bearbeiten",
	"Save":            "Speichern",
	"Cancel":          "Abbrechen",
	"Show":            "Anzeigen",
	"Change":          "Ändern",
	"Check now":       "Jetzt prüfen",
	"Install":         "Installieren",
	"Confirm":         "Bestätigen",
	"Connect":         "Verbinden",
	"Disconnect":      "Trennen",
	"Pair":            "Koppeln",
	"Forget":          "Vergessen",
	"Off":             "Aus",
	"On":              "Ein",
	"None":            "Keine",
	"Automatic":       "Automatisch",
	"Background":      "Hintergrund",
	"Screensaver":     "Bildschirmschoner",
	"Match all":       "Alle erkennen",
	"Repeat":          "Wiederholen",
	"Alarm sound":     "Weckton",
	"New timer":       "Neuer Timer",
	"Wake with light": "Mit Licht wecken",
	"Night starts":    "Nacht beginnt",
	"Night ends":      "Nacht endet",
	"Night starts at": "Nacht beginnt um",
	"Night ends at":   "Nacht endet um",

	// Frequently used pages outside Settings.
	"Radar":                        "Radar",
	"Forecast":                     "Vorhersage",
	"Radio":                        "Radio",
	"Announce":                     "Durchsage",
	"Call":                         "Anrufen",
	"Stations":                     "Sender",
	"Could not list the stations":  "Senderliste nicht verfügbar",
	"No stations within 100 km":    "Keine Sender im Umkreis von 100 km",
	"Try Popular worldwide":        "Weltweit beliebte Sender versuchen",
	"No stations in this list yet": "Noch keine Sender in dieser Liste",
	"Stop the radio":               "Radio stoppen",
	"Stop":                         "Stopp",
	"Play":                         "Abspielen",
	"Playing":                      "Wiedergabe läuft",
	"Playing now":                  "Läuft gerade",
	"Paused":                       "Pausiert",
	"Starting":                     "Starte",
	"Starting…":                    "Starte…",
	"Phone":                        "Telefon",
	"Intercom, in the house":       "Hausruf",
	"Alarm":                        "Wecker",
	"Timer and alarm":              "Timer und Wecker",
	"Silenced":                     "Stummgeschaltet",
	"Listening…":                   "Ich höre zu…",
	"Thinking…":                    "Ich denke nach…",
	"Thinking":                     "Denke nach",
	"microphone off":               "Mikrofon aus",
	"♪ playing":                    "♪ spielt",
	"♪ paused":                     "♪ pausiert",
	"setup: a browser is asking to be let in": "Einrichtung: Browser bittet um Zugriff",
	"MICROPHONE OFF":    "MIKROFON AUS",
	"NEXT":              "ALS NÄCHSTES",
	"Nothing coming up": "Keine Termine",
	"Snoozed until":     "Schlummern bis",
	"Clear":             "Klar",
	"Partly cloudy":     "Teilweise bewölkt",
	"Thunderstorms":     "Gewitter",
	"Sleet":             "Schneeregen",
	"Severe":            "Unwetter",
	"Windy":             "Windig",
	"Sunny":             "Sonnig",
	"Cloudy":            "Bewölkt",
	"Rainy":             "Regnerisch",
	"Pouring":           "Starkregen",
	"Fog":               "Nebel",
	"Hail":              "Hagel",
	"Snowy":             "Schnee",
	"Lightning":         "Blitze",
}
