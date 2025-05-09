import { makeAutoObservable } from 'mobx'
import IChat from '../models/IChat'

class ChatStore {
    store: Map<string | number, IChat>

    constructor() {
        makeAutoObservable(this)
        this.store = new Map<string | number, IChat>()
    }

    set(chat: IChat) {
        makeAutoObservable(chat, {}, { deep: true })
        this.store.set(chat.chat.id, chat)
    }

    get(chatId: string | number) {
        return this.store.get(chatId)
    }

    delete(chatId: string | number) {
        this.store.delete(chatId)
    }

    has(chatId: string | number): boolean {
        return this.store.has(chatId)
    }

    getPrivateChat(userId: string | number): IChat | undefined {
        return Array.from(this.store.values()).find((chat) => !chat.chat.isGroupChat && chat.members.some(member => member.userId == userId))
    }

    clear() {
        this.store.clear()
    }

    getAll(): IChat[] {
        return this.sortChats(Array.from(this.store.values()))
    }

    getGroupChats(): IChat[] {
        return this.sortChats(
            Array.from(this.store.values()).filter((chat) => chat.chat.isGroupChat)
        )
    }

    getPrivateChats(): IChat[] {
        return this.sortChats(
            Array.from(this.store.values()).filter((chat) => !chat.chat.isGroupChat)
        )
    }

    getGroupChatCount(): number {
        return Array.from(this.store.values()).filter((chat) => chat.chat.isGroupChat).length
    }

    getPrivateChatCount(): number {
        return Array.from(this.store.values()).filter((chat) => !chat.chat.isGroupChat).length
    }

    private sortChats(chats: IChat[]): IChat[] {
        return chats.sort((a, b) => {
            if (!a.lastMessage && !b.lastMessage) return 0
            if (!a.lastMessage) return 1
            if (!b.lastMessage) return -1
            return new Date(b.lastMessage.sentAt).getTime() - new Date(a.lastMessage.sentAt).getTime()
        })
    }
}

export default (new ChatStore())