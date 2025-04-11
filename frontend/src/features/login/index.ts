import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'loginFeature', en)
i18n.addResourceBundle('ru', 'loginFeature', ru)

import Login from "./ui/Login"

export { Login }