import { useCallback } from 'react'
import modalStore from '../store/modalStore'

const useModal = () => {
    const openModal = useCallback((title: string, text: string) => modalStore.openModal(title, text), [])
    return { openModal }
}

export default useModal