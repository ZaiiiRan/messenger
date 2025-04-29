import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import IShortUser from '../models/IShortUser'
import shortUserStore from '../store/ShortUserStore'

async function fetchSocialUser(id: number | string): Promise<AxiosResponse<any, any>> {
    const response = await api.get(`/social/users/${id}`)
    updateStore(response)
    return response
}

async function addFriend(id: number | string): Promise<AxiosResponse<any, any>> {
    const response = await api.post(`/social/friends/management/${id}`)
    updateStore(response)
    return response
}

async function removeFriend(id: number | string): Promise<AxiosResponse<any, any>> {
    const response = await api.delete(`/social/friends/management/${id}`)
    updateStore(response)
    return response
}

async function blockUser(id: number | string): Promise<AxiosResponse<any, any>> {
    const response = await api.post(`/social/block/management/${id}`)
    updateStore(response)
    return response
}

async function unblockUser(id: number | string): Promise<AxiosResponse<any, any>> {
    const response = await api.delete(`/social/block/management/${id}`)
    updateStore(response)
    return response
}

function updateStore(response: AxiosResponse<any, any>) {
    const shortUser = response.data.user as IShortUser
    if (shortUserStore.get(shortUser.userId)) {
        shortUserStore.set(shortUser)
    }
}

export { fetchSocialUser, addFriend, removeFriend, blockUser, unblockUser }