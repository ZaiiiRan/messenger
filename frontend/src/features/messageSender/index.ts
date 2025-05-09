import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'messageSender', en)
i18n.addResourceBundle('ru', 'messageSender', ru)

import MessageSender from './ui/MessageSender'
import PrivateMessageSenderDialog from './ui/PrivateMessageSenderDialog'

export { MessageSender, PrivateMessageSenderDialog }