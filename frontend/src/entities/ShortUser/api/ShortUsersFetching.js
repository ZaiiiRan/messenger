import { api } from '../../../shared/api'

class ShortUsersFetching {
    async fetchShortUser(search, limit, offset) {
        const response = await api.post('/social/users/search', { offset, limit, search })
        return response
    }

    async fetchFriends(search, limit, offset) {
        const response = await api.post('social/friends/friend-list', { offset, limit, search })
        return response
    }

    async fetchIncomingFriendRequests(search, limit, offset) {
        const response = await api.post('social/friends/friend-requests/incoming', { offset, limit, search })
        return response
    }

    async fetchOutgoingFriendRequests(search, limit, offset) {
        const response = await api.post('/social/friends/friend-requests/outgoing', { offset, limit, search })
        return response
    }

    async fetchBlackList(search, limit, offset) {
        const response = await api.post('social/block/block-list', { offset, limit, search })
        return response
    }
}

const shortUsersFetching = new ShortUsersFetching()

export { shortUsersFetching }