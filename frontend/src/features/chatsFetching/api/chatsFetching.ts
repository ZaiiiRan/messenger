import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import { normalizeToIChat } from '../../../entities/Chat'
import chatStore from '../../../entities/Chat/store/ChatStore'
import shortUserStore from '../../../entities/SocialUser/store/ShortUserStore'
import { IShortUser } from '../../../entities/SocialUser'
import { IChat } from '../../../entities/Chat'

async function fetchGroupChats(limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('/chats/group-list', { limit, offset })

    response.data.chats.forEach((value: any) => saveChat(value))

    return response
}

async function fetchPrivateChats(limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('/chats/private-list', { limit, offset })

    response.data.chats.forEach((value: any) => saveChat(value))

    return response
}

async function fetchChat(id: string | number): Promise<IChat> {
    const response = await api.get(`/chats/${id}`)
    const chat = saveChat(response.data)
    return chat
}

function saveChat(data: any) {
    const chat = normalizeToIChat(data)
    chatStore.set(chat)
    saveUsers(data)
    return chat
}

function saveUsers(chat: any) {
    chat.members.forEach((value: any) => {
        const user = value.user as IShortUser
        shortUserStore.set(user)
    })
}

export { fetchGroupChats, fetchPrivateChats, fetchChat, saveChat }