import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'startPage', en)
i18n.addResourceBundle('ru', 'startPage', ru)

import StartPage from "./ui/StartPage"

export { StartPage } 