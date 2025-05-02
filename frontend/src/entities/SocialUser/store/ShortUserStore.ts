import { makeAutoObservable } from 'mobx'
import IShortUser from '../models/IShortUser'
import { fetchSocialUserForStore } from '../api/SocialUserFetching'

class ShortUserStore {
    store: Map<string | number, IShortUser>

    constructor() {
        makeAutoObservable(this)
        this.store = new Map<string | number, IShortUser>()
    }

    set(user: IShortUser) {
        makeAutoObservable(user)
        this.store.set(user.userId, user)
    }

    async get(id: string | number): Promise<IShortUser | undefined> {
        if (this.has(id)) {
            return this.store.get(id)
        }
        return fetchSocialUserForStore(id)
    }

    delete(id: string | number): boolean {
        return this.store.delete(id)
    }

    has(id: string | number): boolean {
        return this.store.has(id)
    }
}

export default (new ShortUserStore())