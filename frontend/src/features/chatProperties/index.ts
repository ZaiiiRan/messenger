import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatProperties', en)
i18n.addResourceBundle('ru', 'chatProperties', ru)

import ChatProperties from "./ui/ChatProperties"

export { ChatProperties }
