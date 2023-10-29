package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	_ "github.com/mattn/go-sqlite3"
)

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
	} else {
		fmt.Println("Tabelle existiert bereits.")
	}

	rows, err := db.Query("SELECT * FROM abrechnungen")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	/*if !rows.Next() {
		log.Println("Die Datenbank ist leer.")
		return
	}*/

	// Daten aus der Datenbank in eine Liste konvertieren
	var data [][]string
	for rows.Next() {
		fmt.Println("In der Schleife")
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

	// Funktion zur Erstellung leerer CanvasObject für die Tabelle
	createEmptyTableCell := func() fyne.CanvasObject {
		return widget.NewLabel(".....................")
	}

	// Funktion zur Erstellung der Zelleninhalte für die Tabelle
	createTableCell := func(id widget.TableCellID, content fyne.CanvasObject) {
		content.(*widget.Label).SetText(data[id.Row][id.Col])
	}

	// Tabelle erstellen
	fmt.Println(data)
	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		createEmptyTableCell, // Leere CanvasObject-Funktion
		createTableCell,      // Funktion zur Erstellung der Zelleninhalte
	)
	muenzgeld := widget.NewEntry()
	fuenfScheine := widget.NewEntry()
	zehnScheine := widget.NewEntry()
	zwanzigScheine := widget.NewEntry()
	fuenfzigScheine := widget.NewEntry()
	hundertScheine := widget.NewEntry()
	zweiHundertScheine := widget.NewEntry()
	papierGesamt := widget.NewEntry()
	barGesamt := widget.NewEntry()
	ecGesamt := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Gesamtes Münzgeld", Widget: muenzgeld},
			{Text: "5€ Scheine", Widget: fuenfScheine},
			{Text: "10€ Scheine", Widget: zehnScheine},
			{Text: "20€ Scheine:", Widget: zwanzigScheine},
			{Text: "50€ Scheine:", Widget: fuenfzigScheine},
			{Text: "100€ Scheine:", Widget: hundertScheine},
			{Text: "200€ Scheine:", Widget: zweiHundertScheine},
			{Text: "Gesamtes Papiergeld:", Widget: papierGesamt},
			{Text: "Gesamtes Bargeld:", Widget: barGesamt},
			{Text: "Gesamtes Bargeld mit EC:", Widget: ecGesamt},
		},
		OnSubmit: func() {
			log.Println("Form submitted:", muenzgeld.Text)
		},
	}

	summeRollen := widget.NewEntry()
	summeKarte := widget.NewEntry()
	zBon := widget.NewEntry()
	sonderAus := widget.NewEntry()
	sonderEin := widget.NewEntry()
	stornoViel := widget.NewEntry()
	stornoWenig := widget.NewEntry()

	form2 := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Summe Rollengeld", Widget: summeRollen},
			{Text: "Summe Kartenzahlung", Widget: summeKarte},
			{Text: "zBon Kassenbericht", Widget: zBon},
			{Text: "Sonderausgabe aus Kasse", Widget: sonderAus},
			{Text: "Sondereingabe in Kasse", Widget: sonderEin},
			{Text: "Storno - zu viel", Widget: stornoViel},
			{Text: "Storno - zu wenig", Widget: stornoWenig},
		},
	}

	papierZurueck := widget.NewEntry()
	einzahlung := widget.NewEntry()
	tagesbilanz := widget.NewEntry()

	form3 := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Papiergeld zurück in die Kasse", Widget: papierZurueck},
			{Text: "Einzahlung auf Konto", Widget: einzahlung},
			{Text: "Tagesbilanz vgl- Einnahme zu zBon", Widget: tagesbilanz},
		},
	}

	// Flexibles Layout erstellen und die Tabelle hinzufügen
	forms := container.NewHBox(form, form2, form3)
	formBox := widget.NewCard("Neuer Eintrag", "", forms)
	tableBox := widget.NewCard("Einträge", "", table)
	mainContent := container.NewVBox(
		formBox,
		tableBox,
	)

	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
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
