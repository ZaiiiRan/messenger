import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'
import IShortUser from '../models/IShortUser'
import shortUserStore from '../store/ShortUserStore'

async function fetchSocialUser(id: number | string): Promise<AxiosResponse<any, any>> {
    const response = await api.get(`/social/users/${id}`)
    updateStore(response)
    return response
}

async function fetchSocialUsersForStore(ids: Array<number | string>): Promise<(AxiosResponse<any, any> | undefined)[] > {
    const uniqueIds = [...new Set(ids.map(String))]

    const requests = uniqueIds.map((id) => {
        if (!shortUserStore.has(id)) {
            return api.get(`/social/users/${id}`).catch(() => {
                return undefined
            })
        }
        return Promise.resolve(undefined)
    })
    const responses = await Promise.all(requests)

    responses.forEach(response => {
        if (response) {
            const shortUser = response.data.user as IShortUser
            shortUserStore.set(shortUser)
        }
    })

    return responses
}

async function fetchSocialUserForStore(id: number | string): Promise<IShortUser | undefined> {
    try {
        const response = await api.get(`/social/users/${id}`)
        const user = response.data.user
        shortUserStore.set(user)
        return user
    } catch (error) {
        return undefined
    }
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
    if (shortUserStore.has(shortUser.userId)) {
        shortUserStore.set(shortUser)
    }
}

export { fetchSocialUser, addFriend, removeFriend, blockUser, unblockUser, fetchSocialUserForStore, fetchSocialUsersForStore }