import { api } from '../../../shared/api'
import { IChatInfo, IChatMember, normalizeToIChat, normalizeToIChatMember } from '../../../entities/Chat'
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

async function addMembersToChat(id: string | number, members: (string | number)[]): Promise<void> {
    const response = await api.post(`/chats/${id}/members`, { members })
    const chatInfo = response.data.chat as IChatInfo
    const newMembers = response.data.newMembers

    const chat = chatStore.get(id)
    if (chat) {
        const chatMembers: IChatMember[] = newMembers.map((value: any) => normalizeToIChatMember(value))
        runInAction(() => {
            chat.chat = chatInfo
            const uniqueMembers = chatMembers.filter(
                (newMember) => !chat.members.some((existingMember) => existingMember.userId === newMember.userId)
            )
            chat.members.push(...uniqueMembers)
        })
    }
}

async function removeMembersFromChat(id: string | number, members: (string | number)[]): Promise<void> {
    const response = await api.patch(`/chats/${id}/members`, { members })
    const chatInfo = response.data.chat as IChatInfo
    const removedMembers = response.data.removedMembers

    const chat = chatStore.get(id)
    if (chat) {
        runInAction(() => {
            chat.chat = chatInfo
            chat.members = chat.members.filter((member) => !removedMembers.some((removedMember: any) => removedMember.user.userId === member.userId))
        })
    }
}

async function changeChatMemberRole(id: string | number, memberId: string | number, role: string): Promise<void> {
    const response = await api.patch(`/chats/${id}/members/${memberId}/role`, { role })
    const chatInfo = response.data.chat as IChatInfo
    const updatedMember = normalizeToIChatMember(response.data.member)

    const chat = chatStore.get(id)
    if (chat) {
        runInAction(() => {
            chat.chat = chatInfo
            chat.members = chat.members.map((member) => {
                if (member.userId === updatedMember.userId) {
                    return updatedMember
                }
                return member
            })
        })
    }
}

function saveChat(data: any) {
    const chat = normalizeToIChat(data)
    chatStore.set(chat)
    saveUsers(data.members)
    return chat
}

function saveUsers(users: any) {
    if (users) {
        users.forEach((value: any) => {
            const user = value.user as IShortUser
            shortUserStore.set(user)
        })
    }
}

export { fetchGroupChats, fetchPrivateChats, fetchChat, fetchPrivateChat, deleteChat, leaveFromChat, returnToChat, renameChat, saveChat, addMembersToChat,
    removeMembersFromChat, changeChatMemberRole }
