import IWebSocketMessage from '../../../shared/api/models/IWebSocketMessage'
import { chatStore, IMessage } from '../../../entities/Chat/'
import { runInAction } from 'mobx'
import { fetchChat } from '../../../features/chatsFetching'
import { shortUserStore } from '../../../entities/SocialUser'
import { userStore } from '../../../entities/user'

function handleNewMessageNotification(wsMessage: IWebSocketMessage) {
    const message = wsMessage.content as IMessage

    if (!chatStore.has(message.chatId)) fetchUnknownChat(message.chatId)

    const chat = chatStore.get(message.chatId)
    if (chat) {
        if (!shortUserStore.has(message.memberId) && message.memberId !== userStore.user?.userId) fetchUnknownUser(message.memberId)

        const isNewer = !chat.lastMessage || new Date(message.sentAt).getTime() > new Date(chat.lastMessage.sentAt).getTime()
        if (isNewer) {
            runInAction(() => {
                chat.lastMessage = message
            })
        }

        if (chat.messages.some(item => message.id === item.id)) return
        runInAction(() => {
            chat.messages.unshift(message)
        })
    }
}

async function fetchUnknownChat(chatId: string | number) {
    await fetchChat(chatId)
}

async function fetchUnknownUser(userId: string | number) {
    await shortUserStore.get(userId)
}

export default handleNewMessageNotification