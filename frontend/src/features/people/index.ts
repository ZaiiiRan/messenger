import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'peopleFeature', en)
i18n.addResourceBundle('ru', 'peopleFeature', ru)

import PeopleList from './ui/PeopleList'
export { PeopleList }