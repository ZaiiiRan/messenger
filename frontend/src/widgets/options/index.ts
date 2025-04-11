import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'optionsWidget', en)
i18n.addResourceBundle('ru', 'optionsWidget', ru)

import OptionsList from './ui/OptionsList'
import AppearanceOptions from './ui/AppearanceOptions'
export { OptionsList, AppearanceOptions }