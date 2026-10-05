# Deutsche Bildschirmtexte

Dieser Zweig ergänzt die TECHO5-Version v0.9.30 (Basis-Commit `4d9147524b24a0058f338e36fbdbf9d4f04ab11d`) um deutsche Texte auf dem Gerätebildschirm. Die Zuordnung fester Texte steht in `echod/internal/feature/display/localize_de.go`; die Wort-Uhr bildet Uhrzeiten in `clock_style.go` auf Deutsch. Die Aufrufstellen liegen im selben Display-Paket. Tests prüfen unter anderem die Uhrstil-Auswahl, gesprochene Uhrzeiten und die Darstellung längerer Texte.

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
2. Neue oder geänderte englische Bildschirmtexte in `germanScreenLabels` ergänzen und die Wort-Uhr in `clockWords` prüfen. Kontext beachten: „Sun“ bezeichnet in der Uhrstil-Auswahl die Sonne; die Abkürzung für Sonntag wird separat über `germanWeekdayAbbrev` erzeugt.
3. `go test ./...` im Verzeichnis `echod` ausführen und die Menüs sowie Sprachzustände am Gerät prüfen.

Patches aller deutschen Änderungen seit der Basisversion v0.9.30 lassen sich mit `git format-patch 4d9147524b24a0058f338e36fbdbf9d4f04ab11d..HEAD` erzeugen.

## Neue Originalversionen

Der Workflow [`update-german-fork.yml`](../.github/workflows/update-german-fork.yml) prüft stündlich die neueste veröffentlichte Show-Version von `HuskerMinion/techo5`. Er läuft auch manuell über **Actions → Prepare German fork for upstream releases → Run workflow**. Bei einer neuen Version führt er den Release-Stand in einem eigenen Zweig mit der deutschen Version zusammen, testet Show, Dot und Spot und erstellt einen Pull Request **innerhalb dieses Forks**. Der deutsche Standardzweig ändert sich erst, wenn der Pull Request zusammengeführt wird. Das Gerät selbst wird nicht aktualisiert.

Für die erstmalige Einrichtung muss dieser Workflow im deutschen Standardzweig des Forks liegen. Unter **Actions** die Workflows des Forks aktivieren; GitHub deaktiviert sie beim Forken zunächst. Unter **Settings → Actions → General → Workflow permissions** außerdem **Allow GitHub Actions to create and approve pull requests** einschalten. Der Workflow verwendet das automatisch bereitgestellte `GITHUB_TOKEN`; ein persönliches Token ist nicht nötig. Bei 60 Tagen ohne Repository-Aktivität kann GitHub geplante Workflows wieder deaktivieren; in diesem Fall unter **Actions** erneut aktivieren.

Bei Merge-Konflikten oder fehlgeschlagenen Tests stoppt der Workflow und zeigt den Fehler unter **Actions** an. Dann muss die Übersetzung an die neue Version angepasst werden. Auch bei erfolgreichen Tests neue Bildschirmtexte im Pull Request prüfen: unbekannte Texte bleiben sonst englisch. Die Prüfung erfolgt ungefähr stündlich und kann sich auf GitHub verzögern.
