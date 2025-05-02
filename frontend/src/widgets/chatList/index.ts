import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatListWidget', en)
i18n.addResourceBundle('ru', 'chatListWidget', ru)

import ChatList from "./ui/ChatList"

export { ChatList }