import { observer } from 'mobx-react'
import Member from '../models/member'
import styles from './ChatMember.module.css'
import { IChat } from '../../../entities/Chat'
import { useTranslation } from 'react-i18next'
import { changeChatMemberRole, removeMembersFromChat } from '../../chatsFetching'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import { useModal } from '../../modal'
import IsFetchingStates from '../models/isFetchingStates'
import { Dispatch, SetStateAction, useState } from 'react'
import { SocialUserDialog } from '../../../entities/SocialUser'
import { PrivateMessageSenderDialog } from '../../messageSender'

interface MemberProps {
    chat: IChat
    member: Member,
    isButtonsDisabled: () => boolean,
    isFetching: IsFetchingStates,
    setIsFetching: Dispatch<SetStateAction<IsFetchingStates>>,
}

const getRolePriority = (role: string): number => {
    switch(role) {
        case 'owner': return 3
        case 'admin': return 2
        case 'member': return 1
        default: return 0
    }
}

const ChatMember: React.FC<MemberProps> = ({ chat, member, isButtonsDisabled, isFetching, setIsFetching }) => {
    const { t } = useTranslation('chatProperties')
    const { openModal } = useModal()
    const [showUserModal, setShowUserModal] = useState<boolean>(false)
    const [showSendMessageModal, setShowSendMessageModal] = useState<boolean>(false)

    const openUserModal = () => {
        if (chat.you.userId !== member.userId) 
            setShowUserModal(true)
    }

    const showErrorModal = (e: any) => {
        const errorKey: ApiErrorsKey = e.response?.data?.error
        const errMsg = t(apiErrors[errorKey]) || t('Internal server error')
        openModal(t('Error'), errMsg)
    }

    const checkAccessForRemove = () => {
        if (chat.you.userId === member.userId) return false
        const yourRolePriority = getRolePriority(chat.you.role)
        const memberRolePriority = getRolePriority(member.role)

        if (member.addedBy !== chat.you.userId) {
            if (yourRolePriority <= memberRolePriority) return false
        } else if (member.addedBy === chat.you.userId && memberRolePriority >= yourRolePriority) return false

        return true
    }

    const checkAccessForRoleChanging = (newRole: string) => {
        if (chat.you.userId === member.userId) return false
        const yourRolePriority = getRolePriority(chat.you.role)

        if (yourRolePriority === 1) return false
        const memberRolePriority = getRolePriority(member.role)
        const newRolePriority = getRolePriority(newRole)

        if (newRolePriority === 3) return false
        if (yourRolePriority <= newRolePriority || (memberRolePriority == yourRolePriority && newRolePriority < memberRolePriority) 
            || memberRolePriority > yourRolePriority) return false
        if (newRolePriority === memberRolePriority) return false
        return true
    }

    const removeMemberAction = () => {
        if (isButtonsDisabled()) return
        
        const removeMemberFunc = async () => {
            try {
                setIsFetching({ ...isFetching, removeMember: true })
                await removeMembersFromChat(chat.chat.id, [member.userId])
            } catch (e) {
                showErrorModal(e)
            } finally {
                setIsFetching({ ...isFetching, removeMember: false })
            }
        }
        
        openModal(t('Remove member'), `${t('Are you sure you want to remove')} ${member.user.username}?`, removeMemberFunc)
    }

    const changeMemberRoleAction = (role: string) => {
        if (isButtonsDisabled()) return

        const changeMemberRoleFunc = async () => {
            try {
                setIsFetching({ ...isFetching, changeRole: true })
                await changeChatMemberRole(chat.chat.id, member.userId, role)
            } catch (e) {
                showErrorModal(e)
            } finally {
                setIsFetching({ ...isFetching, changeRole: false })
            }
        }

        openModal(t('Change role'), `${t('Are you sure you want to change')} ${member.user.username}${t('\'s role to')} ${t(role)}?`, changeMemberRoleFunc)
    }
    
    return (
        <div 
            className={`${styles.Member} flex items-center px-3 py-2 
                2k:px-6 2k:py-3 4k:px-9 4k:py-4 rounded-3xl xl:gap-5 mobile:gap-4 2k:gap-8 4k:gap-12`}
        >
            <div 
                className={`md:h-2/3 mobile:h-3/5 rounded-3xl aspect-square ${chat.you.userId !== member.userId ? 'cursor-pointer' : ''}`}
                onClick={openUserModal}
            >
                <div className='flex items-center justify-center w-full h-full Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'>
                    <div className='flex items-center justify-center w-1/2 aspect-square'>
                        {
                            member.user.isActivated && !member.user.isDeleted && !member.user.isBanned ? (
                                <svg viewBox="0 0 16 19" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <defs/>
                                    <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                </svg>
                            ) : (
                                <svg viewBox="-3.5 0 19 19" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M11.383 13.644A1.03 1.03 0 0 1 9.928 15.1L6 11.172 2.072 15.1a1.03 1.03 0 1 1-1.455-1.456l3.928-3.928L.617 5.79a1.03 1.03 0 1 1 1.455-1.456L6 8.261l3.928-3.928a1.03 1.03 0 0 1 1.455 1.456L7.455 9.716z"></path></g></svg>
                            )
                        }
                    </div>
                </div>
            </div>

            <div 
                className={`${chat.you.userId !== member.userId ? styles.MemberText : ''} flex flex-col 4k:gap-2 max-w-[70%] ${chat.you.userId !== member.userId ? 'cursor-pointer' : ''}`}
                onClick={openUserModal}
            >
                <div 
                    className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                        md:text-xl sm:text-lg mobile:text-text-base text-ellipsis overflow-hidden'
                >
                    { member.user.firstname } { member.user.lastname }
                </div>

                <div
                    className='2xl:text-base xl:text-sm lg:text-sm 2k:text-lg 4k:text-xl
                        md:text-base sm:text-sm mobile:text-sm  text-ellipsis overflow-hidden'
                >
                    { member.user.username } { member.role === 'admin' || member.role === 'owner' ? `(${t(member.role)})` : ''} 
                </div>
            </div>
            
            <div className='ml-auto flex gap-6 h-full items-center'>
                {
                    member.role === 'admin' && checkAccessForRoleChanging('member') && (
                        <div 
                            className='Icon IconClickable cursor-pointer flex items-center justify-center h-1/4 aspect-square'
                            onClick={() => changeMemberRoleAction('member')}
                        >
                            <svg viewBox="0 0 1920 1920" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="m.153 526.146 92.168-92.299 867.767 867.636 867.636-867.636 92.429 92.299-960.065 959.935z" fillRule="evenodd"></path> </g></svg>
                        </div>
                    )
                }
                {
                    member.role === 'member' && checkAccessForRoleChanging('admin') && (
                        <div 
                            className='Icon IconClickable cursor-pointer flex items-center justify-center h-1/4 aspect-square'
                            onClick={() => changeMemberRoleAction('admin')}
                        >
                            <svg viewBox="0 0 1920 1920" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="m.153 1393.854 92.168 92.299 867.767-867.636 867.636 867.636 92.429-92.299L960.088 433.92z" fillRule="evenodd"></path> </g></svg>
                        </div>
                    )
                }
                {
                    checkAccessForRemove() && (
                        <div 
                            className='Icon IconClickable cursor-pointer flex items-center justify-center h-1/4 aspect-square'
                            onClick={removeMemberAction}
                        >
                            <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <path fillRule="evenodd" d="M13.41425,12.00025 L18.70725,6.70725 C19.09825,6.31625 19.09825,5.68425 18.70725,5.29325 C18.31625,4.90225 17.68425,4.90225 17.29325,5.29325 L12.00025,10.58625 L6.70725,5.29325 C6.31625,4.90225 5.68425,4.90225 5.29325,5.29325 C4.90225,5.68425 4.90225,6.31625 5.29325,6.70725 L10.58625,12.00025 L5.29325,17.29325 C4.90225,17.68425 4.90225,18.31625 5.29325,18.70725 C5.48825,18.90225 5.74425,19.00025 6.00025,19.00025 C6.25625,19.00025 6.51225,18.90225 6.70725,18.70725 L12.00025,13.41425 L17.29325,18.70725 C17.48825,18.90225 17.74425,19.00025 18.00025,19.00025 C18.25625,19.00025 18.51225,18.90225 18.70725,18.70725 C19.09825,18.31625 19.09825,17.68425 18.70725,17.29325 L13.41425,12.00025 Z"></path> </g></svg>
                        </div>
                    )
                }
                
            </div>

            <SocialUserDialog
                show={showUserModal}
                setShow={(show: boolean) => setShowUserModal(show)}
                id={member.userId}
                onMessageClick={() => setShowSendMessageModal(true) }
            />
            <PrivateMessageSenderDialog
                show={showSendMessageModal}
                setShow={(show: boolean) => setShowSendMessageModal(show)}
                recipient={member.user}
                zIndex={51}
            />
        </div>
    )
}

export default observer(ChatMember)