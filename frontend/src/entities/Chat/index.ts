import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatEntity', en)
i18n.addResourceBundle('ru', 'chatEntity', ru)

import useChatStore from './hook/useChatStore'
import IChat from './models/IChat'
import IChatInfo from './models/IChatInfo'
import IChatMember from './models/IChatMember'
import IMessage from './models/IMessage'
import ChatCard from './ui/ChatCard'
import ChatCardSkeleton from './ui/ChatCardSkeleton'

export { useChatStore, ChatCard, ChatCardSkeleton }
export type { IChat, IChatInfo, IChatMember, IMessage }