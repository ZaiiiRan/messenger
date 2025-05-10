import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { IChat } from '../../../entities/Chat'
import IsFetchingStates from '../models/isFetchingStates'
import { useTranslation } from 'react-i18next'
import { useModal } from '../../modal'
import { observer } from 'mobx-react'
import Member from '../models/member'
import { userStore } from '../../../entities/user'
import { IShortUser, shortUserStore } from '../../../entities/SocialUser'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import ChatMember from './ChatMember'
import { Input } from '../../../shared/ui/Input';

interface MembersManagementProps {
    chat: IChat,
    isButtonsDisabled: () => boolean,
    isMember: () => boolean,
    isAdmin: () => boolean,
    isOwner: () => boolean,
    isFetching: IsFetchingStates,
    setIsFetching: Dispatch<SetStateAction<IsFetchingStates>>
}

const MemberManagement: React.FC<MembersManagementProps> = ({ chat, isButtonsDisabled, isMember, isAdmin, isOwner, isFetching, setIsFetching }) => {
    const { t } = useTranslation('chatProperties')
    const { openModal } = useModal()
    const [members, setMembers] = useState<Member[]>([])
    const [search, setSearch] = useState<string>('')

    const showErrorModal = (e: any) => {
        const errorKey: ApiErrorsKey = e.response?.data?.error
        const errMsg = t(apiErrors[errorKey]) || t('Internal server error')
        openModal(t('Error'), errMsg)
    }

    useEffect(() => {
        let isMounted = true

        const getMembers = async () => {
            const you: Member = {
                ... chat.you,
                user: userStore.user as IShortUser
            }

            let members: Member[] = []
            const chatMembers = chat.members
            for (let i = 0; i < chatMembers.length; i++) {
                try {
                    const user = await shortUserStore.get(chatMembers[i].userId)
                    if (user) {
                        members.unshift({
                            ... chatMembers[i],
                            user: user as IShortUser
                        })
                    }
                } catch (e: any) {
                    showErrorModal(e)
                }
            }

            members = [you, ...members]
            setMembers(members)
        }

        if (isMounted) {
            getMembers()
        }

        return () => {
            isMounted = false
        }
    }, [chat.members, chat.you, chat.members.length, search])

    const filteredMembers = members.filter((member) => {
        const searchLower = search.trim().toLowerCase()
        const username = member.user.username.toLowerCase() || ''
        return username.includes(searchLower)
    })

    return(
        <>
            <div
                className='font-extrabold 
                    md:text-lg mobile:text-base 2k:text-2xl 4k:text-4xl'
            >
                { t('Members') }
            </div>

            <div className='h-12 w-full'>
                    <Input
                        placeholder={t('Username')}
                        className='w-full px-2 py-1 2k:px-3 2k:py-2 4k:px-4 4k:py-35 rounded-lg 
                            md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                        value={search}
                        onChange={(e) => setSearch(e.target.value) }
                    />
            </div>

            <div className='scrollContainer flex flex-col overflow-y-scroll gap-3 w-full box-border flex-grow'>
                {
                    filteredMembers.map((member) => (
                        <ChatMember
                            key={member.userId}
                            chat={chat}
                            member={member}
                            isButtonsDisabled={isButtonsDisabled}
                            isFetching={isFetching}
                            setIsFetching={setIsFetching}
                        />
                    ))
                }
                {
                    filteredMembers.length === 0 && (
                        <div 
                            className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                                md:text-xl sm:text-lg mobile:text-text-base text-center'
                        >
                            { t('We couldn\'t find anyone') }
                        </div>
                    )
                }
            </div>
        </>
    )
}

export default observer(MemberManagement)
