const { contextBridge, ipcRenderer } = require("electron");

const rows = ipcRenderer.invoke('db-query', 'SELECT * FROM abrechnungen');

contextBridge.exposeInMainWorld('api', {
    invoke: (channel, data) => {
        let validChannels = ['db-query'];
        if (validChannels.includes(channel)) {
            return ipcRenderer.invoke(channel, data);
        }
    }
})
