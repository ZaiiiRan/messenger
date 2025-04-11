import { createRoot } from 'react-dom/client'
import './index.css'
import App from './app/App'

createRoot(document.getElementById('root') as HTMLElement).render(
  <>
    <App />
  </>,
)

const getPreferredColorScheme = () => {
  const darkQuery = "(prefers-color-scheme: dark)" as const
  const darkMQL = window.matchMedia ? window.matchMedia(darkQuery) : ({} as MediaQueryList)
  if (darkMQL.media === darkQuery && darkMQL.matches) {
    return "dark"
  }
  return "light"
}
document.documentElement.setAttribute("data-color-scheme", getPreferredColorScheme())