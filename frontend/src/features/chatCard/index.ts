import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatCard', en)
i18n.addResourceBundle('ru', 'chatCard', ru)

import ChatCard from './ui/ChatCard'
import ChatCardSkeleton from './ui/ChatCardSkeleton'

export { ChatCard, ChatCardSkeleton }