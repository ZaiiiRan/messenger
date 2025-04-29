import { makeAutoObservable } from 'mobx'
import IChat from '../models/IChat'

class ChatStore {
    store: Map<string | number, IChat>

    constructor() {
        makeAutoObservable(this)
        this.store = new Map<string | number, IChat>()
    }

    set(chat: IChat) {
        this.store.set(chat.chat.id, chat)
    }

    get(chatId: string | number) {
        return this.store.get(chatId)
    }

    delete(chatId: string | number) {
        this.store.delete(chatId)
    }

    getAll() {
        return this.store.values().toArray()
    }
}

export default (new ChatStore())