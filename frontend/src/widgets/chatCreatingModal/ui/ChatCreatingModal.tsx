import { useTranslation } from 'react-i18next'
import { useModal } from '../../../features/modal'
import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { Button } from '../../../shared/ui/Button'
import styles from './ChatCreatingModal.module.css'
import { Input } from '../../../shared/ui/Input'
import { fetchFriends } from '../../../entities/SocialUser'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import createChat from '../api/createChat'
import { Loader } from '../../../shared/ui/Loader'
import { validateChatName, validateMembers } from '../../../entities/Chat'
import { UserSelection } from '../../../features/userSelection'
import { Dialog } from '../../../shared/ui/Dialog'

interface ChatCreatingModalProps {
    show: boolean,
    setShow: (show: boolean) => void,
    open: (chatId: string | number) => void
}

const ChatCreatingModal: React.FC<ChatCreatingModalProps> = ({ show, setShow, open }) => {
    const { t } = useTranslation('chatCreating')
    const { openModal } = useModal()
    const [chatName, setChatName] = useState<string>('')
    const [chatNameError, setChatNameError] = useState<boolean>(false)
    const [selectedUsers, setSelectedUsers] = useState<(number | string)[]>([])
    const [isCreating, setIsCreating] = useState<boolean>(false)

    const validate = () => {
        const showError = (message: string | undefined) => {
            if (message) {
                openModal(t('Error'), t(message))
            }
            
        }
    
        let error = validateChatName(chatName)
        if (error.error) {
            setChatNameError(true)
            showError(error.message)
            return false
        }
        setChatNameError(false)
    
        error = validateMembers(selectedUsers)
        if (error.error) {
            showError(error.message)
            return false
        }
    
        return true
    }

    const submit = async () => {
        if (!validate()) return
        setIsCreating(true)

        try {
            const chat = await createChat(chatName.trim(), selectedUsers)
            setShow(false)
            open(chat.chat.id)
        } catch (e: any) {
            const errorKey: ApiErrorsKey = e.response?.data?.error
            const errMsg = t(apiErrors[errorKey]) || t('An unexpected error occurred (maybe one of the users removed you from friends)')
            openModal(t('Error'), errMsg)
        } finally {
            setIsCreating(false)
        }
    }

    const selectUser = (id: number | string) => {
        if (!selectedUsers.includes(id)) {
            setSelectedUsers((prevSelectedUsers) => [...prevSelectedUsers, id])
        } else {
            setSelectedUsers((prevSelectedUsers) => prevSelectedUsers.filter((selectedId) => selectedId !== id))
        }
    }

    useEffect(() => {
        setChatName('')
        setSelectedUsers([])
    }, [show])

    return (
        <Dialog
            show={show}
            setShow={setShow}
            title={t('New chat')}
        >
            <div className='flex gap-6 2k:gap-10 4k:gap-12'>
                {/* Avatar */}
                <div className='md:h-24 mobile:h-16 2k:h-32 4k:h-40 rounded-3xl aspect-square'>
                    <div className='flex items-center justify-center w-full h-full Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'>
                        <div className='flex items-center justify-center w-1/2 aspect-square'>
                            <svg viewBox="0 0 26.6666 24" xmlns="http://www.w3.org/2000/svg">
                                <defs/>
                                <path id="coolicon" d="M2.66 6.66C2.66 2.98 5.65 0 9.33 0C13.01 0 16 2.98 16 6.66C16 10.34 13.01 13.33 9.33 13.33C5.65 13.33 2.66 10.34 2.66 6.66ZM9.33 2.66C7.12 2.66 5.33 4.45 5.33 6.66C5.33 8.87 7.12 10.66 9.33 10.66C11.54 10.66 13.33 8.87 13.33 6.66C13.33 4.45 11.54 2.66 9.33 2.66ZM18.66 6.66C19.08 6.66 19.5 6.76 19.87 6.95C20.25 7.14 20.57 7.42 20.82 7.76C21.07 8.1 21.23 8.5 21.3 8.91C21.36 9.33 21.33 9.75 21.2 10.16C21.07 10.56 20.84 10.92 20.54 11.22C20.25 11.51 19.88 11.74 19.48 11.87C19.08 12 18.66 12.03 18.24 11.96C17.82 11.89 17.43 11.73 17.09 11.48L15.52 13.64L15.52 13.64C16.2 14.13 16.99 14.46 17.82 14.59C17.91 14.61 18 14.62 18.09 14.63C18.84 14.71 19.59 14.63 20.3 14.4C21.1 14.14 21.83 13.7 22.43 13.1C23.02 12.51 23.47 11.78 23.73 10.98C23.99 10.18 24.06 9.33 23.93 8.5C23.8 7.67 23.47 6.88 22.98 6.2C22.54 5.59 21.98 5.08 21.33 4.71C21.25 4.66 21.17 4.62 21.08 4.58C20.33 4.19 19.5 4 18.66 4L18.66 6.66ZM16 24L18.66 24C18.66 18.84 14.48 14.66 9.33 14.66C4.17 14.66 0 18.84 0 24L2.66 24C2.66 20.31 5.65 17.33 9.33 17.33C13.01 17.33 16 20.31 16 24ZM23.99 24C23.99 23.3 23.85 22.6 23.59 21.96C23.32 21.31 22.93 20.72 22.43 20.23C21.94 19.73 21.35 19.34 20.7 19.07C20.05 18.8 19.36 18.66 18.66 18.66L18.66 16C19.57 16 20.47 16.15 21.33 16.45C21.46 16.5 21.59 16.55 21.72 16.6C22.69 17.01 23.58 17.6 24.32 18.34C25.06 19.08 25.65 19.96 26.05 20.93C26.11 21.06 26.16 21.2 26.2 21.33C26.51 22.18 26.66 23.09 26.66 24L23.99 24Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                            </svg>  
                        </div>
                    </div>
                </div>
                <div className='h-12 w-full'>
                    <Input
                        placeholder={t('Chat name')}
                        className='w-full px-2 py-1 2k:px-3 2k:py-2 4k:px-4 4k:py-35 rounded-lg 
                            md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                        value={chatName}
                        onChange={(e) => setChatName(e.target.value) }
                        error={chatNameError}
                    />
                </div>
                
            </div>

            <UserSelection onSelect={selectUser} fetchFunction={fetchFriends} checkSelected={(id) => selectedUsers.includes(id)} />

            <div className='self-end'>
                <Button
                    className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-56 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                        rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                    onClick={submit}
                    disabled={isCreating}
                >
                    {isCreating ? (
                        <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                    ) : t('Create') }
                </Button>
            </div>
        </Dialog>
    )
}

export default ChatCreatingModal
