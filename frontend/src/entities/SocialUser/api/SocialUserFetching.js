import { api } from '../../../shared/api'

class SocialUserAPI {
    async fetch(id) {
        return api.post('/social/get-user', { user_id: id })
    }

    async addFriend(id) {
        return api.post('/social/add-friend', { user_id: id })
    }

    async removeFriend(id) {
        return api.post('/social/remove-friend', { user_id: id })
    }

    async blockUser(id) {
        return api.post('/social/block-user', { user_id: id })
    }

    async unblockUser(id) {
        return api.post('/social/unblock-user', { user_id: id })
    }
}

const socialUserAPI = new SocialUserAPI()

export default socialUserAPI