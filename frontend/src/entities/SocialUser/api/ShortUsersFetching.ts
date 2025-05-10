import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import shortUserStore from '../store/ShortUserStore'
import IShortUser from '../models/IShortUser'

async function fetchShortUser(search: string, limit: number, offset: number): Promise<IShortUser[]> {
    const response = await api.post('/social/users/search', { offset, limit, search })
    const users = normalizeShortUserList(response)
    return users
}

async function fetchFriends(search: string, limit: number, offset: number): Promise<IShortUser[]> {
    const response = await api.post('/social/friends/friend-list', { offset, limit, search })
    const users = normalizeShortUserList(response)
    return users
}

async function fetchIncomingFriendRequests(search: string, limit: number, offset: number): Promise<IShortUser[]> {
    const response = await api.post('/social/friends/friend-requests/incoming', { offset, limit, search })
    const users = normalizeShortUserList(response)
    return users
}

async function fetchOutgoingFriendRequests(search: string, limit: number, offset: number): Promise<IShortUser[]> {
    const response = await api.post('/social/friends/friend-requests/outgoing', { offset, limit, search })
    const users = normalizeShortUserList(response)
    return users
}

async function fetchBlackList(search: string, limit: number, offset: number): Promise<IShortUser[]> {
    const response = await api.post('/social/block/block-list', { offset, limit, search })
    const users = normalizeShortUserList(response)
    return users
}

async function fetchFriendsAreNotChatting(chatId: number | string, search: string, limit: number, offset: number): Promise<IShortUser[]> {
    const response = await api.post(`/chats/${chatId}/members/friends-are-not-chatting`, { offset, limit, search })
    const users = normalizeShortUserList(response)
    return users
}

function normalizeShortUserList(response: AxiosResponse<any, any>): IShortUser[] {
    const users = response.data.users

    const result = users.map((user: any) => {
        const shortUser = user as IShortUser
        updateStore(user)
        return shortUser
    })

    return result
}

async function updateStore(user: IShortUser) {
    let candidate = await shortUserStore.get(user.userId)
    if (candidate) {
        shortUserStore.set(user)
    }
}

export { fetchShortUser, fetchFriends, fetchIncomingFriendRequests, fetchOutgoingFriendRequests, fetchBlackList, fetchFriendsAreNotChatting }