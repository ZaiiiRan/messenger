import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'messageList', en)
i18n.addResourceBundle('ru', 'messageList', ru)

import MessageList from './ui/MessageList'

export { MessageList }