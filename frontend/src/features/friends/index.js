import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'friendsFeature', en)
i18n.addResourceBundle('ru', 'friendsFeature', ru)

import UserList from './ui/UserList'
import FriendsMenu from './ui/FriendsMenu'
export { UserList, FriendsMenu }