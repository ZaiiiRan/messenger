import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import { IChatInfo, normalizeToIChat } from '../../../entities/Chat'
import chatStore from '../../../entities/Chat/store/ChatStore'
import shortUserStore from '../../../entities/SocialUser/store/ShortUserStore'
import { IShortUser } from '../../../entities/SocialUser'
import { IChat } from '../../../entities/Chat'
import { runInAction } from 'mobx'

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

async function deleteChat(id: string | number): Promise<void> {
    const response = await api.delete(`/chats/${id}`)
    const chatInfo = response.data.deletedChat as IChatInfo
    chatStore.delete(chatInfo.id)
}

async function leaveFromChat(id: string | number): Promise<void> {
    const response = await api.patch(`/chats/${id}/leave`)
    console.log(response.data)
}

function saveChat(data: any) {
    const chat = normalizeToIChat(data)
    chatStore.set(chat)
    saveUsers(data)
    return chat
}

function saveUsers(chat: any) {
    if (chat.members) {
        chat.members.forEach((value: any) => {
            const user = value.user as IShortUser
            shortUserStore.set(user)
        })
    }
}

export { fetchGroupChats, fetchPrivateChats, fetchChat, deleteChat, leaveFromChat, saveChat }