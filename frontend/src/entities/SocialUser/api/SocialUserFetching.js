import { api } from '../../../shared/api'

class SocialUserAPI {
    async fetch(id) {
        return api.get(`/social/users/${id}`)
    }

    async addFriend(id) {
        return api.post(`/social/friends/management/${id}`)
    }

    async removeFriend(id) {
        return api.delete(`/social/friends/management/${id}`)
    }

    async blockUser(id) {
        return api.post(`/social/block/management/${id}`)
    }

    async unblockUser(id) {
        return api.delete(`/social/block/management/${id}`)
    }
}

const socialUserAPI = new SocialUserAPI()

export default socialUserAPI