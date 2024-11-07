import { useCallback } from 'react'
import modalStore from '../../../app/stores/modalStore/modalStore'

const useModal = () => {
    const openModal = useCallback(() => modalStore.openModal(), [])
    const closeModal = useCallback(() => modalStore.closeModal(), [])
    const setModalTitle = useCallback((title) => modalStore.setTitle(title), [])
    const setModalText = useCallback((text) => modalStore.setText(text), [])
    
    return { isOpen: modalStore.isOpen, openModal, closeModal, setModalTitle, setModalText }
}

export default useModal