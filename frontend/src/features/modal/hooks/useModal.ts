import { useCallback } from 'react'
import modalStore from '../store/modalStore'

const useModal = () => {
    const openModal = useCallback(() => modalStore.openModal(), [])
    const closeModal = useCallback(() => modalStore.closeModal(), [])
    const setModalTitle = useCallback((title: string) => modalStore.setTitle(title), [])
    const setModalText = useCallback((text: string) => modalStore.setText(text), [])
    
    return { isOpen: modalStore.isOpen, openModal, closeModal, setModalTitle, setModalText }
}

export default useModal