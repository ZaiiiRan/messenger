import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'optionsFeature', en)
i18n.addResourceBundle('ru', 'optionsFeature', ru)

import LanguageChanging from './ui/LaguageChanging'
import ThemeChanging from './ui/ThemeChanging'
export { LanguageChanging, ThemeChanging }