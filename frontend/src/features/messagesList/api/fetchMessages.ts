import { runInAction } from 'mobx'
import { chatStore, IMessage } from '../../../entities/Chat'
import { shortUserStore } from '../../../entities/SocialUser'
import { userStore } from '../../../entities/user'
import { api } from '../../../shared/api'

async function fetchMessages(chatId: number | string, limit: number, offset: number): Promise<IMessage[]> {
    const response = await api.post(`/chats/${chatId}/messages-list`, { limit, offset })

    const messagesObjs = response.data.messages

    const messages: IMessage[] = messagesObjs.map((messageObj: any) => {
        const message = messageObj as IMessage

        if (message.memberId != userStore.user?.userId && !shortUserStore.has(message.memberId)) fetchUnknownUser(message.memberId)

        return message
    })

    const chat = chatStore.get(chatId)
    if (chat) {
        runInAction(() => {
            chat.messages = [...chat.messages, ...messages]
        })
    }

    return messages
}

async function fetchUnknownUser(userId: string | number) {
    await shortUserStore.get(userId)
}

export default fetchMessages