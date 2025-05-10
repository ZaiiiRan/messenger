import { makeAutoObservable } from 'mobx'
import ModalData from '../models/modalData'

class ModalStore {
    modals: ModalData[] = []

    constructor() {
        makeAutoObservable(this)
    }

    openModal(title: string, text: string, actionFunction: (() => void) | undefined = undefined) {
        if (actionFunction) {
            this.modals.push({ id: Date.now().toString(), title, text, actionFunction })
            return
        } 
        
        if (this.modals.length >= 10) return
        const id = Date.now().toString()
        this.modals.push({ id, title, text })
    }

    closeModal(id: string) {
        this.modals = this.modals.filter(modal => modal.id !== id)
    }
}

const modalStore = new ModalStore()
export default modalStore