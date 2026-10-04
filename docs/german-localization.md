# Deutsche Bildschirmtexte

Dieser Zweig ergänzt die TECHO5-Version v0.9.30 (Basis-Commit `4d9147524b24a0058f338e36fbdbf9d4f04ab11d`) um deutsche Texte auf dem Gerätebildschirm. Die Zuordnung steht in `echod/internal/feature/display/localize_de.go`; die Aufrufstellen liegen im selben Display-Paket. Tests prüfen unter anderem die Uhrstil-Auswahl und die Darstellung längerer Texte.

Die Übersetzung greift nur bei der Darstellung. Interne Kennungen, gespeicherte Einstellungen und Home-Assistant-Entitäten behalten ihre ursprünglichen Werte. Nicht erfasste Texte sowie die Web-Einrichtungsseite bleiben englisch. Die Einstellung „Screen language“ betrifft weiterhin die Sprache gesprochener Bildschirmbefehle.

Zum Prüfen und Bauen im Repository:

```sh
cd echod
go test ./...
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o ../bin/echod-arm ./cmd/echod
```

Weitere Bau- und Installationshinweise stehen in [building.md](building.md). Ein eigener Daemon-Build ist erforderlich; der Quellcode allein ändert die Firmware auf dem Gerät nicht.

## Auf spätere Versionen übertragen

1. Den neuen Stand von `HuskerMinion/techo5` als Basis verwenden und diesen Commit mit `git cherry-pick` übernehmen. Bei Konflikten die Änderungen im Display-Paket an die neue Version anpassen.
2. Neue oder geänderte englische Bildschirmtexte in `germanScreenLabels` ergänzen. Kontext beachten: „Sun“ bezeichnet in der Uhrstil-Auswahl die Sonne; die Abkürzung für Sonntag wird separat über `germanWeekdayAbbrev` erzeugt.
3. `go test ./...` im Verzeichnis `echod` ausführen und die Menüs sowie Sprachzustände am Gerät prüfen.

Der eigenständige Patch zur Basisversion v0.9.30 kann bei Bedarf mit `git format-patch -1 --stdout` aus diesem Commit erzeugt werden.
