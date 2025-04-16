import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'messengingPage', en)
i18n.addResourceBundle('ru', 'messengingPage', ru)

import MessengingPage from './ui/MessengingPage'
export { MessengingPage }