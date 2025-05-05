import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatCreating', en)
i18n.addResourceBundle('ru', 'chatCreating', ru)

import ChatCreatingModal from './ui/ChatCreatingModal'

export { ChatCreatingModal }