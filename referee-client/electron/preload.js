const { contextBridge, ipcRenderer } = require('electron')

contextBridge.exposeInMainWorld('electronAPI', {
  // Expose APIs here
  ping: () => ipcRenderer.invoke('ping')
})
