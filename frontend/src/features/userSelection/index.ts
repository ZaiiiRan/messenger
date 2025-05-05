import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'userSelection', en)
i18n.addResourceBundle('ru', 'userSelection', ru)

import UserSelection from "./ui/UserSelection"

export { UserSelection }