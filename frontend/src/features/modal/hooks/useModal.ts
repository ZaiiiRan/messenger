import { useCallback } from 'react'
import modalStore from '../store/modalStore'

const useModal = () => {
    const openModal = useCallback((title: string, text: string, actionFunction: (() => void) | undefined = undefined) => modalStore.openModal(title, text, actionFunction), [])
    return { openModal }
}

export default useModal