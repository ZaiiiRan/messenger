import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'

class SocialUserAPI {
    async fetch(id: number | string): Promise<AxiosResponse<any, any>> {
        return api.get(`/social/users/${id}`)
    }

    async addFriend(id: number | string): Promise<AxiosResponse<any, any>> {
        return api.post(`/social/friends/management/${id}`)
    }

    async removeFriend(id: number | string): Promise<AxiosResponse<any, any>> {
        return api.delete(`/social/friends/management/${id}`)
    }

    async blockUser(id: number | string): Promise<AxiosResponse<any, any>> {
        return api.post(`/social/block/management/${id}`)
    }

    async unblockUser(id: number | string): Promise<AxiosResponse<any, any>> {
        return api.delete(`/social/block/management/${id}`)
    }
}

const socialUserAPI = new SocialUserAPI()

export default socialUserAPI