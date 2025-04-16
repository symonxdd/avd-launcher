export function getStateClass(state) {
  switch (state) {
    case 'Running':
      return 'running'
    case 'Powered off':
      return 'poweredOff'
    case 'Loading...':
      return 'loading'
    case 'Shutting down...':
      return 'shuttingDown'
    default:
      return ''
  }
}
