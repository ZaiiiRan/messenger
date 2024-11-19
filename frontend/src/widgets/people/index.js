import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'peopleWidget', en)
i18n.addResourceBundle('ru', 'peopleWidget', ru)

import PeopleMenu from './ui/PeopleMenu'
import PeopleListWidget from './ui/PeopleListWidget'
import UserWidget from './ui/UserWidget'
export { PeopleMenu, PeopleListWidget, UserWidget }