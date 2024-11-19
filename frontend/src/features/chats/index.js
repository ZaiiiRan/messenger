import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'chatsFeature', en)
i18n.addResourceBundle('ru', 'chatsFeature', ru)

import PeopleChatsList from './ui/PeopleChatsList'
import GroupChatsList from './ui/GroupChatsList'
import ChatWidget from './ui/ChatWidget'
import SendMessageModal from './ui/SendMessageModal'

export { PeopleChatsList, GroupChatsList, ChatWidget, SendMessageModal }