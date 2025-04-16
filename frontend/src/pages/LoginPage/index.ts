import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'loginPage', en)
i18n.addResourceBundle('ru', 'loginPage', ru)

import LoginPage from "./ui/LoginPage"

export { LoginPage }