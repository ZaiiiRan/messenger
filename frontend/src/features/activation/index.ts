import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'activationFeature', en)
i18n.addResourceBundle('ru', 'activationFeature', ru)


import ActivationAccount from "./ui/ActivationAccount"

export { ActivationAccount }