import { makeAutoObservable } from 'mobx'

class ModalStore {
    isOpen: boolean = false
    title: string = ''
    text: string = ''

    constructor() {
        makeAutoObservable(this)
    }

    openModal() {
        this.isOpen = true
    }

    closeModal() {
        this.isOpen = false
    }

    setTitle(title: string) {
        this.title = title
    }

    setText(text: string) {
        this.text = text
    }

}

const modalStore = new ModalStore()
export default modalStore