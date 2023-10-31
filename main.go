/* Dies ist ein Programm zur berechnung der Tagesbilanz.
* Erstellt von Patryk Hegenberg
* EMail: patrykhegenberg@gmail.com
 */
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"

	_ "github.com/mattn/go-sqlite3"
)

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' {
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *numericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func checkSqliteExistence() {
	// this function checks if the file abrechnung.db exists
	// and creates it if it doesn't
	_, err := os.Stat("./abrechnung.db")
	if os.IsNotExist(err) {
		os.Create("abrechnung.db")
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Tagesabrechnung")

	// Verbindung zur SQLite-Datenbank herstellen
	db, err := sql.Open("sqlite3", "./abrechnung.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prüfen, ob die Tabelle vorhanden ist
	tableExists := tableExists(db, "abrechnungen")
	if !tableExists {
		createTable(db)
	} else {
		fmt.Println("Tabelle existiert bereits.")
	}

	data := getTableData(db)
	lastBilanz := getLastBilanz(db)
	table := createTableWidget(data)

	mainContent := createMainContent(db, table, lastBilanz)

	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.ShowAndRun()
}

// TODO: createMainContent functions needs to be cleaned up to make it more readable and maintainable
func createMainContent(db *sql.DB, table *widget.Table, lastBilanz string) *fyne.Container {
	papierGesamt := newNumericalEntry()
	papierGesamt.SetPlaceHolder("0.00")
	barGesamt := newNumericalEntry()
	barGesamt.SetPlaceHolder("0.00")
	ecGesamt := newNumericalEntry()
	ecGesamt.SetPlaceHolder("0.00")
	summeRollen := newNumericalEntry()
	summeRollen.SetPlaceHolder("0.00")
	summeKarte := newNumericalEntry()
	summeKarte.SetPlaceHolder("0.00")
	zBon := newNumericalEntry()
	zBon.SetPlaceHolder("0.00")
	sonderAus := newNumericalEntry()
	sonderAus.SetPlaceHolder("0.00")
	sonderEin := newNumericalEntry()
	sonderEin.SetPlaceHolder("0.00")
	stornoViel := newNumericalEntry()
	stornoViel.SetPlaceHolder("0.00")
	stornoWenig := newNumericalEntry()
	stornoWenig.SetPlaceHolder("0.00")
	papierZurueck := newNumericalEntry()
	papierZurueck.SetPlaceHolder("0.00")
	einzahlung := newNumericalEntry()
	einzahlung.SetPlaceHolder("0.00")
	tagesbilanz := newNumericalEntry()
	tagesbilanz.SetPlaceHolder("0.00")

	muenzgeld := newNumericalEntry()
	muenzgeld.SetPlaceHolder("0.00")
	fuenfScheine := newNumericalEntry()
	fuenfScheine.SetPlaceHolder("0")
	zehnScheine := newNumericalEntry()
	zehnScheine.SetPlaceHolder("0")
	zwanzigScheine := newNumericalEntry()
	zwanzigScheine.SetPlaceHolder("0")
	fuenfzigScheine := newNumericalEntry()
	fuenfzigScheine.SetPlaceHolder("0")
	hundertScheine := newNumericalEntry()
	hundertScheine.SetPlaceHolder("0")
	zweiHundertScheine := newNumericalEntry()
	zweiHundertScheine.SetPlaceHolder("0")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Münzgeld", Widget: muenzgeld},
			{Text: "5€ Scheine", Widget: fuenfScheine},
			{Text: "10€ Scheine", Widget: zehnScheine},
			{Text: "20€ Scheine:", Widget: zwanzigScheine},
			{Text: "50€ Scheine:", Widget: fuenfzigScheine},
			{Text: "100€ Scheine:", Widget: hundertScheine},
			{Text: "200€ Scheine:", Widget: zweiHundertScheine},
		},
		OnSubmit: func() {
			err := updateBarGesamt(summeRollen.Text, muenzgeld.Text, fuenfScheine.Text, zehnScheine.Text, zwanzigScheine.Text, fuenfzigScheine.Text, hundertScheine.Text, zweiHundertScheine.Text, barGesamt)
			if err != nil {
				log.Fatal(err)
			}
			err = updatePapierGesamt(fuenfScheine.Text, zehnScheine.Text, zwanzigScheine.Text, fuenfzigScheine.Text, hundertScheine.Text, zweiHundertScheine.Text, papierGesamt)
			if err != nil {
				log.Fatal(err)
			}
			err = updateEcGesamt(barGesamt.Text, summeKarte.Text, ecGesamt)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	form.SubmitText = "Zwischenberechnung"

	form2 := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Summe Rollengeld", Widget: summeRollen},
			{Text: "Summe Kartenzahlung", Widget: summeKarte},
			{Text: "zBon Kassenbericht", Widget: zBon},
			{Text: "Sonder aus Kasse", Widget: sonderAus},
			{Text: "Sonder in Kasse", Widget: sonderEin},
			{Text: "Storno - zu viel", Widget: stornoViel},
			{Text: "Storno - zu wenig", Widget: stornoWenig},
		},
		OnSubmit: func() {
			err := updateBarGesamt(summeRollen.Text, muenzgeld.Text, fuenfScheine.Text, zehnScheine.Text, zwanzigScheine.Text, fuenfzigScheine.Text, hundertScheine.Text, zweiHundertScheine.Text, barGesamt)
			if err != nil {
				log.Fatal(err)
			}
			err = updatePapierGesamt(fuenfScheine.Text, zehnScheine.Text, zwanzigScheine.Text, fuenfzigScheine.Text, hundertScheine.Text, zweiHundertScheine.Text, papierGesamt)
			if err != nil {
				log.Fatal(err)
			}
			err = updateEcGesamt(barGesamt.Text, summeKarte.Text, ecGesamt)
			if err != nil {
				log.Fatal(err)
			}
			err = calcTagesbilanz(ecGesamt.Text, zBon.Text, sonderAus.Text, sonderEin.Text, stornoWenig.Text, stornoViel.Text, lastBilanz, tagesbilanz)
			if err != nil {
				log.Fatal(err)
			}
			err = zurueckKasse(muenzgeld.Text, summeRollen.Text, papierZurueck)
			if err != nil {
				log.Fatal(err)
			}
			err = Einzahlung(barGesamt.Text, papierZurueck.Text, einzahlung)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	form2.SubmitText = "Berechnen"

	form3 := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Gesamtes Papiergeld:", Widget: papierGesamt},
			{Text: "Gesamtes Bargeld:", Widget: barGesamt},
			{Text: "Gesamtes Bargeld mit EC:", Widget: ecGesamt},
			{Text: "Papiergeld zurück:", Widget: papierZurueck},
			{Text: "Einzahlung:", Widget: einzahlung},
			{Text: "Tagesbilanz:", Widget: tagesbilanz},
		},
		OnSubmit: func() {
			now := time.Now().Format("2006-01-02")
			insertQuery := `INSERT INTO abrechnungen (datum, einzahlung, tagesbilanz, bargeld) VALUES (?, ?, ?, ?)`
			_, err := db.Exec(insertQuery, now, einzahlung.Text, tagesbilanz.Text, papierZurueck.Text)
			if err != nil {
				log.Fatal(err)
			}
			updateTable(db)
			table.Refresh()

			fmt.Println("Daten gespeichert.")
		},
	}
	form3.SubmitText = "Speichern"

	// Flexibles Layout erstellen und die Tabelle hinzufügen
	forms := container.NewGridWithColumns(3, form, form2, form3)
	formBox := widget.NewCard("Neuer Eintrag", "", forms)
	tableBox := widget.NewCard("Einträge", "", table)
	mainContent := container.NewGridWithRows(2,
		formBox,
		tableBox,
	)
	return mainContent
}

func createTableWidget(data [][]string) *widget.Table {
	createEmptyTableCell := func() fyne.CanvasObject {
		return widget.NewLabel(".....................")
	}

	// Funktion zur Erstellung der Zelleninhalte für die Tabelle
	createTableCell := func(id widget.TableCellID, content fyne.CanvasObject) {
		content.(*widget.Label).SetText(data[id.Row][id.Col])
	}

	// Tabelle erstellen
	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		createEmptyTableCell, // Leere CanvasObject-Funktion
		createTableCell,      // Funktion zur Erstellung der Zelleninhalte
	)
	return table
}

func getLastBilanz(db *sql.DB) string {
	lastBilanz := ""
	last, err := db.Query("SELECT tagesbilanz FROM abrechnungen ORDER BY datum DESC Limit 1")
	if err != nil {
		log.Fatal("Can't query last entry ", err)
	}
	defer last.Close()
	for last.Next() {
		var tagesbilanz float64
		err := last.Scan(&tagesbilanz)
		if err != nil {
			log.Fatal(err)
		}
		lastBilanz = fmt.Sprintf("%.2f", tagesbilanz)
	}
	return lastBilanz
}

// Funktion zur Überprüfung, ob eine Tabelle in der Datenbank existiert
func tableExists(db *sql.DB, tableName string) bool {
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func createTable(db *sql.DB) {
	// Wenn die Tabelle nicht vorhanden ist, sie erstellen
	createTableQuery := `
			CREATE TABLE abrechnungen (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				datum TEXT,
				einzahlung REAL,
				tagesbilanz REAL,
				bargeld REAL
			);
		`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Datum des Vortags im Format "YYYY-MM-DD" erhalten
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// Mit einem ersten Eintrag befüllen
	insertDataQuery := `
			INSERT INTO abrechnungen (datum, einzahlung, tagesbilanz, bargeld)
			VALUES (? , 0.00, 0.00, 300.00);
		`
	_, err = db.Exec(insertDataQuery, yesterday)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tabelle erstellt und mit Daten befüllt.")
}

func getTableData(db *sql.DB) [][]string {
	rows, err := db.Query("SELECT * FROM abrechnungen")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Daten aus der Datenbank in eine Liste konvertieren
	var data [][]string
	for rows.Next() {
		var id sql.NullInt64
		var datum string
		var einzahlung, tagesbilanz, bargeld float64
		err := rows.Scan(&id, &datum, &einzahlung, &tagesbilanz, &bargeld)
		if err != nil {
			log.Fatal(err)
		}

		var idString string
		if id.Valid {
			idString = fmt.Sprintf("%d", id.Int64)
		} else {
			idString = "NULL"
		}

		fmt.Println(idString, datum, fmt.Sprintf("%.2f", einzahlung), fmt.Sprintf("%.2f", tagesbilanz), fmt.Sprintf("%.2f", bargeld))
		data = append(data, []string{idString, datum, fmt.Sprintf("%.2f", einzahlung), fmt.Sprintf("%.2f", tagesbilanz), fmt.Sprintf("%.2f", bargeld)})
	}
	return data
}

func updateBarGesamt(rollen, muenz, fuenf, zehn, zwanzig, fuenfzig, hundert, zweihundert string, barGesamt fyne.Widget) error {
	rollenf, err := convertToFloat(rollen)
	if err != nil {
		log.Fatal("Error converting Rollengeld", err)
	}
	muenzf, err := convertToFloat(muenz)
	if err != nil {
		log.Fatal("Error converting Münzgeld", err)
	}
	fuenff, err := convertToInt(fuenf)
	if err != nil {
		log.Fatal("Error converting Fünf", err)
	}
	zehnf, err := convertToInt(zehn)
	if err != nil {
		log.Fatal("Error converting Zehn", err)
	}
	zwanzigf, err := convertToInt(zwanzig)
	if err != nil {
		log.Fatal("Error converting Zwanzig", err)
	}
	fuenfzigf, err := convertToInt(fuenfzig)
	if err != nil {
		log.Fatal("Error converting Fünfzig", err)
	}
	hundertf, err := convertToInt(hundert)
	if err != nil {
		log.Fatal("Error converting Hundert", err)
	}
	zweihundertf, err := convertToInt(zweihundert)
	if err != nil {
		log.Fatal("Error converting Zweihundert", err)
	}
	result := rollenf + muenzf + float64(fuenff*5) + float64(zehnf*10) + float64(zwanzigf*20) + float64(fuenfzigf*50) + float64(hundertf*100) + float64(zweihundertf*200)

	barGesamt.(*numericalEntry).SetText(fmt.Sprintf("%.2f", result))
	return nil
}

func updatePapierGesamt(fuenf, zehn, zwanzig, fuenfzig, hundert, zweihundert string, papierGesamt fyne.Widget) error {
	fuenff, err := convertToInt(fuenf)
	if err != nil {
		log.Fatal("Error converting Fünf", err)
	}
	zehnf, err := convertToInt(zehn)
	if err != nil {
		log.Fatal("Error converting Zehn", err)
	}
	zwanzigf, err := convertToInt(zwanzig)
	if err != nil {
		log.Fatal("Error converting Zwanzig", err)
	}
	fuenfzigf, err := convertToInt(fuenfzig)
	if err != nil {
		log.Fatal("Error converting Fünfzig", err)
	}
	hundertf, err := convertToInt(hundert)
	if err != nil {
		log.Fatal("Error converting Hundert", err)
	}
	zweihundertf, err := convertToInt(zweihundert)
	if err != nil {
		log.Fatal("Error converting Zweihundert", err)
	}
	result := float64(fuenff*5) + float64(zehnf*10) + float64(zwanzigf*20) + float64(fuenfzigf*50) + float64(hundertf*100) + float64(zweihundertf*200)

	papierGesamt.(*numericalEntry).SetText(fmt.Sprintf("%.2f", result))
	return nil
}

func updateEcGesamt(bar, ec string, ecGesamt fyne.Widget) error {
	barf, err := convertToFloat(bar)
	if err != nil {
		log.Fatal("Error converting BarGesamt", err)
	}
	ecf, err := convertToFloat(ec)
	if err != nil {
		log.Fatal("Error converting EC", err)
	}
	result := barf + ecf
	ecGesamt.(*numericalEntry).SetText(fmt.Sprintf("%.2f", result))
	return nil
}

func convertToFloat(text string) (float64, error) {
	if text != "" {
		result, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Fatal(err)
		}
		return result, nil
	}
	return 0.00, nil
}

func convertToInt(text string) (int, error) {
	if text != "" {
		result, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal(err)
			return 0.00, err
		}
		return result, nil
	}
	return 0.00, nil
}

func calcTagesbilanz(ecGesamt, zBon, sonderOut, sonderIn, stornoWenig, stornoViel, last string, tagesbilanz fyne.Widget) error {
	ecGesamtf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	zBonf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	sonderOutf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	sonderInf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	stornoWenigf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	stornoVielf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	lastf, err := convertToFloat(ecGesamt)
	if err != nil {
		return err
	}
	result := ecGesamtf + lastf - zBonf + sonderOutf - sonderInf - stornoWenigf + stornoVielf
	tagesbilanz.(*numericalEntry).SetText(fmt.Sprintf("%.2f", result))

	return nil
}

func zurueckKasse(muenz, rollen string, zurueck fyne.Widget) error {
	muenzf, err := convertToFloat(muenz)
	if err != nil {
		return err
	}
	rollenf, err := convertToFloat(rollen)
	if err != nil {
		return err
	}
	result := 300.00 - muenzf - rollenf
	zurueckKasse := float64(math.Floor((result / 5)) * 5)
	zurueck.(*numericalEntry).SetText(fmt.Sprintf("%.2f", zurueckKasse))
	return nil
}

func Einzahlung(gesamtBar, zurueckKasse string, einzahlung fyne.Widget) error {
	barf, err := convertToFloat(gesamtBar)
	if err != nil {
		return err
	}
	kassef, err := convertToFloat(zurueckKasse)
	if err != nil {
		return err
	}
	result := math.Floor((barf-kassef)/5) * 5
	if result < 0 {
		einzahlung.(*numericalEntry).SetText("0.00")
	} else {
		einzahlung.(*numericalEntry).SetText(fmt.Sprintf("%.2f", result))
	}
	return nil
}

func updateTable(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM abrechnungen")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data [][]string
	for rows.Next() {
		var id sql.NullInt64
		var datum string
		var einzahlung, tagesbilanz, bargeld float64
		err := rows.Scan(&id, &datum, &einzahlung, &tagesbilanz, &bargeld)
		if err != nil {
			log.Fatal(err)
		}

		var idString string
		if id.Valid {
			idString = fmt.Sprintf("%d", id.Int64)
		} else {
			idString = "NULL"
		}

		fmt.Println(idString, datum, fmt.Sprintf("%.2f", einzahlung), fmt.Sprintf("%.2f", tagesbilanz), fmt.Sprintf("%.2f", bargeld))
		data = append(data, []string{idString, datum, fmt.Sprintf("%.2f", einzahlung), fmt.Sprintf("%.2f", tagesbilanz), fmt.Sprintf("%.2f", bargeld)})
	}
}
