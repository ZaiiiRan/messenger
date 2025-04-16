import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'activationPage', en)
i18n.addResourceBundle('ru', 'activationPage', ru)

import ActivationPage from "./ui/ActivationPage"

export { ActivationPage }