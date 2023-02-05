let abrechnungen = window.api.invoke('db-query', "SELECT * FROM abrechnungen").then(function(res) {
    console.log(res);
    return res;
})
    .catch(function(err) {
        console.error(err);
    });

const eintragen = () => {
    window.api.invoke('db-query', "INSERT INTO abrechnungen (datum, Einzahlung, Tagesbilanz, Bargeld) VALUES (date(), "
        + parseFloat(document.getElementById("aufKonto").value) + ","
        + parseFloat(document.getElementById("tagesbilanz").value) + ","
        + parseFloat(document.getElementById("zurueckKasse").value) + ");").then(function(res) {
            console.log(res);
        })
        .catch(function(err) {
            console.error(err);
        });
}


document.getElementById("save-button").addEventListener('click', () => {
    eintragen()
});
document.getElementById("show-button").addEventListener('click', () => {
    document.getElementById("main").style.display = "none"
    document.getElementById("Ausgaben").style.display = "block"
    anzeige();
});
document.getElementById("in-button").addEventListener('click', () => {
    document.getElementById("Eingabe").style.display = "block"
    document.getElementById("main").style.display = "none"
});
document.getElementById("back-button").addEventListener('click', () => {
    document.getElementById("Eingabe").style.display = "none"
    document.getElementById("main").style.display = "block"
});
document.getElementById("back2-button").addEventListener('click', () => {
    document.getElementById("Ausgaben").style.display = "none"
    document.getElementById("main").style.display = "block"
});

const updateGesamtBargeld = () => {
    const muenzen = parseFloat(document.getElementById("gesamtMuenzen").value);
    const papier = parseFloat(document.getElementById("gesamtPapier").value);
    const rollengeld = parseFloat(document.getElementById("Rollengeld").value);
    document.getElementById("gesamtBar").value = (muenzen + papier + rollengeld).toFixed(2);
    getGesamtMitEC();
    zurueckKasse();
    Einzahlung();
    Tagesbilanz();
}

const getPapierGeld = () => {
    const euro5 = document.getElementById("5Euro").value * 5;
    const euro10 = document.getElementById("10Euro").value * 10;
    const euro20 = document.getElementById("20Euro").value * 20;
    const euro50 = document.getElementById("50Euro").value * 50;
    const euro100 = document.getElementById("100Euro").value * 100;
    const euro200 = document.getElementById("200Euro").value * 200;

    document.getElementById("gesamtPapier").value = euro5 + euro10 + euro20 + euro50 + euro100 + euro200;
    updateGesamtBargeld();
}
const getGesamtMitEC = () => {
    const ec = parseFloat(document.getElementById("Kartenzahlung").value);
    const gesamtBar = parseFloat(document.getElementById("gesamtBar").value);
    document.getElementById("gesamtBarEC").value = (ec + gesamtBar).toFixed(2);
}

const Tagesbilanz = () => {
    const gesamtMitEc = parseFloat(document.getElementById("gesamtBarEC").value);
    const zBonBericht = parseFloat(document.getElementById("zBon-Bericht").value);
    const sonderOut = parseFloat(document.getElementById("Sonderausgabe").value);
    const sonderIn = parseFloat(document.getElementById("Sondereingabe").value);
    const stornoWenig = parseFloat(document.getElementById("Storno-wenig").value);
    const stornoViel = parseFloat(document.getElementById("Storno-viel").value);

    const result = gesamtMitEc - zBonBericht + sonderOut - sonderIn - stornoWenig + stornoViel;
    document.getElementById("tagesbilanz").value = result.toFixed(2);
}
const zurueckKasse = () => {
    const muenzgeld = parseFloat(document.getElementById("gesamtMuenzen").value);
    const rollengeld = parseFloat(document.getElementById("Rollengeld").value);
    const gesamt = 300 - muenzgeld - rollengeld;
    document.getElementById("zurueckKasse").value = (Math.floor(gesamt / 5) * 5).toFixed(2);
}

const Einzahlung = () => {
    const gesamtBar = parseFloat(document.getElementById("gesamtBar").value);
    const zurueckKasse = parseFloat(document.getElementById("zurueckKasse").value);
    document.getElementById("aufKonto").value = (gesamtBar - zurueckKasse).toFixed(2);
}

const createListeners = () => {
    document.getElementById("5Euro").addEventListener("change", getPapierGeld);
    document.getElementById("10Euro").addEventListener("change", getPapierGeld);
    document.getElementById("20Euro").addEventListener("change", getPapierGeld);
    document.getElementById("50Euro").addEventListener("change", getPapierGeld);
    document.getElementById("100Euro").addEventListener("change", getPapierGeld);
    document.getElementById("200Euro").addEventListener("change", getPapierGeld);

    document.getElementById("gesamtMuenzen").addEventListener("change", updateGesamtBargeld);
    document.getElementById("gesamtPapier").addEventListener("change", updateGesamtBargeld);
    document.getElementById("Rollengeld").addEventListener("change", updateGesamtBargeld);

    document.getElementById("Kartenzahlung").addEventListener("change", getGesamtMitEC);

    document.getElementById("Kartenzahlung").addEventListener("change", Tagesbilanz);
    document.getElementById("zBon-Bericht").addEventListener("change", Tagesbilanz);
    document.getElementById("Sonderausgabe").addEventListener("change", Tagesbilanz);
    document.getElementById("Sondereingabe").addEventListener("change", Tagesbilanz);
    document.getElementById("Storno-wenig").addEventListener("change", Tagesbilanz);
    document.getElementById("Storno-viel").addEventListener("change", Tagesbilanz);

}

const anzeige = () => {
    console.log(abrechnungen);
    let table = "<tr><th>Datum</th><th>Einzahlung</th><th>Tagesbilanz</th><th>Bargeld</th></tr>";
    abrechnungen.then(function(data) {
        for (let i = 0; i < data.length; i++) {
            table += "<tr><td>" + data[i].datum + "</td><td>"
                + data[i].einzahlung + "</td><td>" + data[i].tagesbilanz
                + "</td><td>" + data[i].bargeld + "</td></tr>";
        }
        document.getElementById("table").innerHTML = table;
    });
}

createListeners();
