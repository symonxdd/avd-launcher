import { defineStore } from 'pinia'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { AvdState } from '../enums/avdState'

export const useAvdStore = defineStore('avdStore', {
  state: () => ({
    avds: [],
    logs: '',
    _watcherStarted: false,
    _logListener: null,
    _bootedListener: null,
    _shutdownListener: null
  }),

  actions: {
    appendLog(line) {
      this.logs += line + '\n'
    },

    updateAvdStatus(name, update) {
      const index = this.avds.findIndex(a => a.name === name)
      if (index !== -1) {
        this.avds[index] = {
          ...this.avds[index], // preserve existing props
          ...update            // override with updates
        }
      }
    },

    startAvdWatcher() {
      if (this._watcherStarted) {
        console.log(`[STORE EVENT] Watcher already started`)
        return
      }
      else {
        console.log(`[STORE EVENT] Starting watcher...`)
      }
      this._watcherStarted = true

      this._logListener = (line) => {
        this.appendLog(line)
      }

      this._bootedListener = (name) => {
        this.updateAvdStatus(name, {
          state: AvdState.RUNNING
        })
      }

      this._shutdownListener = (name) => {
        this.updateAvdStatus(name, {
          state: AvdState.POWERED_OFF
        })
      }

      EventsOn('avd-log', this._logListener)
      EventsOn('avd-booted', this._bootedListener)
      EventsOn('avd-shutdown', this._shutdownListener)
    }
  }
})
