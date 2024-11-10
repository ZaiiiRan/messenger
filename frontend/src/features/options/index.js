import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'optionsFeature', en)
i18n.addResourceBundle('ru', 'optionsFeature', ru)

import OptionsList from "./ui/OptionsList"
import AppearanceOptions from './ui/AppearanceOptions'
export { OptionsList, AppearanceOptions }