import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'registerPage', en)
i18n.addResourceBundle('ru', 'registerPage', ru)

import RegisterPage from "./ui/RegisterPage"

export { RegisterPage }