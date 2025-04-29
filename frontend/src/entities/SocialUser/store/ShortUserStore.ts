import { makeAutoObservable } from 'mobx'
import IShortUser from '../models/IShortUser'

class ShortUserStore {
    store: Map<string | number, IShortUser>

    constructor() {
        makeAutoObservable(this)
        this.store = new Map<string | number, IShortUser>()
    }

    set(user: IShortUser) {
        this.store.set(user.userId, user)
    }

    get(id: string | number): IShortUser | undefined {
        return this.store.get(id)
    }

    delete(id: string | number): boolean {
        return this.store.delete(id)
    }
}

export default (new ShortUserStore())