const { app, BrowserWindow, ipcMain } = require('electron');
require('electron-reloader')(module);
const sqlite = require('sqlite-electron');

const path = require('path');
const sqlite3 = require('sqlite3');
const db = new sqlite3.Database('./abrechnung.db');

const createWindow = () => {
    const win = new BrowserWindow({
        width: 1000,
        height: 900,
        webPreferences: {
            nodeIntegration: true,
            preload: path.join(__dirname, 'preload.js')
        }
    });
    ipcMain.handle('db-query', async (event, sqlQuery) => {
        return new Promise(res => {
            db.all(sqlQuery, (err, rows) => {
                if (err) throw err;
                res(rows);
            });
        });
    });
    win.webContents.openDevTools();
    win.loadFile('index.html');
}

app.whenReady().then(() => {
    createWindow();

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) createWindow()
    })
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') app.quit();
})

ipcMain.handle('databasePath', async (event, dbPath) => {
    return await sqlite.setdbPath(dbPath)
});

ipcMain.handle('executeQuery', async (event, query, fetch, value) => {
  return await sqlite.executeQuery(query, fetch, value);
});

db.serialize(() => {
    db.run("CREATE TABLE IF NOT EXISTS abrechnungen (id INTEGER PRIMARY KEY, datum TEXT, einzahlung REAL, tagesbilanz REAL, bargeld REAL)");
});

