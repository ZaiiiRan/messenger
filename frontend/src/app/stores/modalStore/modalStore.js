import { makeAutoObservable } from 'mobx'

class ModalStore {
    isOpen = false
    title = ''
    text = ''

    constructor() {
        makeAutoObservable(this)
    }

    openModal() {
        this.isOpen = true
    }

    closeModal() {
        this.isOpen = false
    }

    setTitle(title) {
        this.title = title
    }

    setText(text) {
        this.text = text
    }

}

const modalStore = new ModalStore()
export default modalStore