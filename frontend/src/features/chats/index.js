import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatsFeature', en)
i18n.addResourceBundle('ru', 'chatsFeature', ru)

import ChatWidget from './ui/ChatWidget'
import SendMessageModal from './ui/SendMessageModal'

export { ChatWidget, SendMessageModal }