import { AxiosError } from 'axios'
import { fetchPrivateChat, saveChat } from '../../chatsFetching/api/chatsFetching'
import sendMessage from './sendMessage'
import { api } from '../../../shared/api'
import { chatStore, IChat } from '../../../entities/Chat'

async function sendPrivateMessage(userId: string | number, message: string) {
    const chat = await getOrCreatePrivateChat(userId)
    sendMessage(chat.chat.id, message)
}

async function getOrCreatePrivateChat(userId: string | number): Promise<IChat> {
    const chat = chatStore.getPrivateChat(userId)
    if (chat) {
        return chat
    }

    try {
        return await fetchPrivateChat(userId)
    } catch (e: any) {
        if (e instanceof AxiosError && e.status === 403) {
            return await createPrivateChat(userId)
        }
        throw e
    }
}

async function createPrivateChat(userId: string | number): Promise<IChat> {
    const response = await api.post('/chats', { members: [userId], isGroup: false })

    const chat = saveChat(response.data)

    return chat
}

export default sendPrivateMessage