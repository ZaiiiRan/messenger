import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import shortUserStore from '../store/ShortUserStore'
import IShortUser from '../models/IShortUser'

async function fetchShortUser(search: string, limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('/social/users/search', { offset, limit, search })
    updateStore(response)
    return response
}

async function fetchFriends(search: string, limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('social/friends/friend-list', { offset, limit, search })
    updateStore(response)
    return response
}

async function fetchIncomingFriendRequests(search: string, limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('social/friends/friend-requests/incoming', { offset, limit, search })
    updateStore(response)
    return response
}

async function fetchOutgoingFriendRequests(search: string, limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('/social/friends/friend-requests/outgoing', { offset, limit, search })
    updateStore(response)
    return response
}

async function fetchBlackList(search: string, limit: number, offset: number): Promise<AxiosResponse<any, any>> {
    const response = await api.post('social/block/block-list', { offset, limit, search })
    updateStore(response)
    return response
}

function updateStore(response: AxiosResponse<any, any>) {
    const users = response.data.users

    users.forEach((user: any) => {
        let shortUser = user as IShortUser
        let candidate = shortUserStore.get(shortUser.userId)
        if (candidate) {
            shortUserStore.set(shortUser)
        }
    })
}

export { fetchShortUser, fetchFriends, fetchIncomingFriendRequests, fetchOutgoingFriendRequests, fetchBlackList }