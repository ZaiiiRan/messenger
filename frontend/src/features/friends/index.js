import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'friendsFeature', en)
i18n.addResourceBundle('ru', 'friendsFeature', ru)

import FriendsMenu from "./ui/FriendsMenu"
import FindFriends from './ui/FindFriends'
import Friends from './ui/Friends'
import IncomingFriendRequests from './ui/IncomingFriendRequests'
import OutgoingFriendRequests from './ui/OutgoingFriendRequests'
import BlackList from './ui/BlackList'
export { FriendsMenu, FindFriends, Friends, IncomingFriendRequests, OutgoingFriendRequests, BlackList }