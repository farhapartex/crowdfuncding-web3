import { initializeFaro, getWebInstrumentations } from '@grafana/faro-web-sdk'

const FARO_URL = import.meta.env.VITE_FARO_URL

export function initFrontendObservability() {
  if (!FARO_URL) return

  initializeFaro({
    url: FARO_URL,
    app: {
      name: 'crowdfunding-frontend',
      environment: import.meta.env.MODE,
    },
    instrumentations: getWebInstrumentations(),
  })
}
