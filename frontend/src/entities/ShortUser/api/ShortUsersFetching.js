import { api } from '../../../shared/api'

class ShortUsersFetching {
    async fetchShortUser(search, limit, offset) {
        const response = await api.post('/social/search-users', { offset, limit, search })
        return response
    }

    async fetchFriends(search, limit, offset) {
        const response = await api.post('/social/get-friends', { offset, limit, search })
        return response
    }

    async fetchIncomingFriendRequests(search, limit, offset) {
        const response = await api.post('/social/get-incoming-friend-requests', { offset, limit, search })
        return response
    }

    async fetchOutgoingFriendRequests(search, limit, offset) {
        const response = await api.post('/social/get-outgoing-friend-requests', { offset, limit, search })
        return response
    }

    async fetchBlackList(search, limit, offset) {
        const response = await api.post('/social/get-blocked-users', { offset, limit, search })
        return response
    }
}

const shortUsersFetching = new ShortUsersFetching()

export { shortUsersFetching }