import { api } from '../../../shared/api'
import { IChatInfo, IChatMember, normalizeToIChat } from '../../../entities/Chat'
import chatStore from '../../../entities/Chat/store/ChatStore'
import shortUserStore from '../../../entities/SocialUser/store/ShortUserStore'
import { IShortUser } from '../../../entities/SocialUser'
import { IChat } from '../../../entities/Chat'
import { runInAction } from 'mobx'

async function fetchGroupChats(limit: number, offset: number): Promise<IChat[]> {
    const response = await api.post('/chats/group-list', { limit, offset })

    const chats = response.data.chats.map((value: any) => saveChat(value))

    return chats
}

async function fetchPrivateChats(limit: number, offset: number): Promise<IChat[]> {
    const response = await api.post('/chats/private-list', { limit, offset })

    const chats = response.data.chats.map((value: any) => saveChat(value))

    return chats
}

async function fetchChat(id: string | number): Promise<IChat> {
    const response = await api.get(`/chats/${id}`)
    const chat = saveChat(response.data)
    return chat
}

async function fetchPrivateChat(memberId: string | number): Promise<IChat> {
    const response = await api.get(`/chats/private/${memberId}`)
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
    const you = response.data.you as IChatMember
    const chatInfo = response.data.chat as IChatInfo

    const chat = chatStore.get(id)
    if (chat) {
        runInAction(() => {
            chat.chat = chatInfo,
            chat.you = you
        })
    }
}

async function returnToChat(id: string | number): Promise<void> {
    const response = await api.patch(`/chats/${id}/return`)
    const you = response.data.you as IChatMember
    const chatInfo = response.data.chat as IChatInfo

    const chat = chatStore.get(id)
    if (chat) {
        runInAction(() => {
            chat.chat = chatInfo,
            chat.you = you
        })
    }
}

async function renameChat(id: string | number, newName: string): Promise<void> {
    const response = await api.patch(`/chats/${id}`, { name: newName })
    const chatInfo = response.data.chat as IChatInfo
    
    const chat = chatStore.get(id)
    if (chat) {
        runInAction(() => {
            chat.chat = chatInfo
        })
    }
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

export { fetchGroupChats, fetchPrivateChats, fetchChat, fetchPrivateChat, deleteChat, leaveFromChat, returnToChat, renameChat, saveChat }