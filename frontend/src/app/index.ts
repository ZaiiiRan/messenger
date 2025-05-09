import { i18n }  from '..//shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'app', en)
i18n.addResourceBundle('ru', 'app', ru)

import App from './App'
export { App }