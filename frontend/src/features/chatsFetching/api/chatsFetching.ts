import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import normalizeToIChat from '../models/normalizeToIChat'
import chatStore from '../../../entities/Chat/store/ChatStore'
import shortUserStore from '../../../entities/SocialUser/store/ShortUserStore'
import { IShortUser } from '../../../entities/SocialUser'

async function fetchGroupChats(limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('/chats/group-list', { limit, offset })

    response.data.chats.forEach((value: any) => {
        const chat = normalizeToIChat(value)
        chatStore.set(chat)

        saveUsers(value)
    })

    return response
}

async function fetchPrivateChats(limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('/chats/private-list', { limit, offset })

    response.data.chats.forEach((value: any) => {
        const chat = normalizeToIChat(value)
        chatStore.set(chat)

        saveUsers(value)
    })

    return response
}

function saveUsers(chat: any) {
    chat.members.forEach((value: any) => {
        const user = value.user as IShortUser
        shortUserStore.set(user)
    })
}

export { fetchGroupChats, fetchPrivateChats }