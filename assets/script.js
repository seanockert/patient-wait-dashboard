import { createApp } from './petite-vue.es.js'

function Dashboard(props) {
  return {
    busyLevel: 5,
    currentTime: 0,
    hoursMax: 8,
    levels: [
      { id: 0, label: 'critical', delay: 40 },
      { id: 1, label: 'high', delay: 30 },
      { id: 2, label: 'medium', delay: 20 },
      { id: 3, label: 'low', delay: 10 }
    ],
    patients: [],
    showAdmin: false,
    socket: null,
    storageKey: 'ed_dashboard',
    waitTime: 0,
    waitTimeAccuracy: 40,
    mounted() {
      this.initDashboard()
      this.initWebsocket()
    },
    initDashboard() {
      // Get patients list from local storage
      const storedData = localStorage.getItem(this.storageKey)
      this.patients = storedData ? JSON.parse(storedData) : []

      this.currentTime = new Date()
      this.updateWaitTime()

      // Update patients list when local storage changes
      window.addEventListener('storage', this.handleStorageChange)

      // Start clock
      setInterval(() => {
        this.currentTime = new Date()
      }, 1000)
    },
    initWebsocket() {
      this.socket = new WebSocket(
        `${window.location.protocol === 'https:' ? 'wss://' : 'ws://'}${
          window.location.host
        }`
      )

      this.socket.addEventListener('open', (event) => {
        console.log('WebSocket connection opened:', event)
      })

      this.socket.addEventListener('message', (event) => {
        const data = JSON.parse(event.data)
        if (data) {
          this.patients = data.patients
          this.updateState()
        }
      })
    },
    formatTime(time, showSeconds = false) {
      const format = {
        hour: 'numeric',
        minute: 'numeric',
        hour12: true
      }

      return new Intl.DateTimeFormat(
        'en-AU',
        showSeconds ? { ...format, second: 'numeric' } : format
      ).format(time)
    },
    formatDate(date) {
      const format = {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: 'numeric',
        minute: 'numeric',
        hour12: true
      }

      return new Intl.DateTimeFormat('en-AU', format).format(date)
    },
    getWaitTime(time1, time2, inHours = false) {
      const timeDifference = Math.abs(time1 - time2)

      if (inHours) {
        // Round to closest hour up to hoursMax
        const hoursDifference = Math.round(timeDifference / (1000 * 60 * 60))
        return hoursDifference > this.hoursMax
          ? `${this.hoursMax}+`
          : hoursDifference
      } else {
        // Return hours and minutes
        const hours = Math.floor(timeDifference / (1000 * 60 * 60))
        const minutes = Math.floor(
          (timeDifference % (1000 * 60 * 60)) / (1000 * 60)
        )

        return `${hours} hr ${minutes} min`
      }
    },
    updatePatient(level = 2, remove = false) {
      if (remove) {
        const patientIndex = this.patients
          .reverse()
          .findIndex((patient) => patient.level === level)

        if (patientIndex !== -1) {
          this.patients.splice(patientIndex, 1)
        }
      } else {
        this.patients.push({
          date: new Date().getTime(),
          level
        })
      }

      this.socket.send(JSON.stringify({ patients: this.patients }))
      this.updateState()
    },
    updateState() {
      localStorage.setItem(this.storageKey, JSON.stringify(this.patients))
      this.updateWaitTime()
    },
    updateWaitTime() {
      // Sum up the patients' delays
      const currentTime = new Date()
      const totalWaitTime = this.patients.reduce((acc, patient) => {
        const level = this.levels.find((level) => level.id === patient.level)
        if (level) {
          acc += level.delay
        }

        return acc
      }, 0)

      const busyLevel = (totalWaitTime / 60 / this.hoursMax) * 100
      this.busyLevel = busyLevel > 100 ? 100 : busyLevel

      this.waitTime = currentTime.setMinutes(
        currentTime.getMinutes() + totalWaitTime
      )
    },
    getPatientsByLevel(level) {
      return this.patients.filter((patient) => patient.level === level).length
    },
    handleStorageChange(event) {
      if (event.key === this.storageKey) {
        this.patients = JSON.parse(event.newValue)
      }
    }
  }
}

createApp({
  Dashboard
}).mount()
