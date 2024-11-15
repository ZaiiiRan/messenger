/* eslint-disable react/prop-types */
import styles from './SocialUser.module.css'
import { useTranslation } from 'react-i18next'
import { Button } from '../../../shared/ui/Button'
import socialUserAPI from '../api/SocialUserFetching'
import { useModal } from '../../../features/modal'
import { apiErrors } from '../../../shared/api'
import { useState } from 'react'
import { Loader } from '../../../shared/ui/Loader'

const SocialUserInfo = ({ data, onUpdate, setUserManipulation }) => {
    const { t } = useTranslation('socialUser')
    const { openModal, setModalTitle, setModalText } = useModal()
    const [isLoading, setLoading] = useState(false)

    const addFriend = async (action) => {
        try {
            setLoading(true)
            await socialUserAPI.addFriend(data.user.id)
            setModalTitle(t('Success'))
            if (action === 'request') {
                setModalText(t('Friend request sent'))
            } else if (action === 'add') {
                if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been added as a friend'))
                else setModalText(`${data.user.username} ${t('has been added as a friend')}`)
            }
            
            openModal()
            await onUpdate()
            setUserManipulation(true)
        } catch (e) {
            setModalTitle(t('Error'))
            if (e.status === 404) {
                setModalText(t('User not found'))
            }
            setModalText(t(apiErrors[e.response?.data?.error]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    const removeFriend = async (action) => {
        try {
            setLoading(true)
            await socialUserAPI.removeFriend(data.user.id)
            setModalTitle(t('Success'))
            if (action === 'decline') {
                if (data.user.username.length > 15) setModalText(t('Friend request from') + ' ' + t('user') + ' ' + t('was rejected'))
                else setModalText(`${t('Friend request from')} ${data.user.username} ${t('was rejected')}`)
            } else if (action === 'cancel') {
                setModalText(t('Friend request canceled'))
            } else if (action === 'remove') {
                if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been removed from friends'))
                else setModalText(`${data.user.username} ${t('has been removed from friends')}`)
            }
            openModal()
            await onUpdate()
            setUserManipulation(true)
        } catch (e) {
            setModalTitle(t('Error'))
            if (e.status === 404) {
                setModalText(t('User not found'))
            }
            setModalText(t(apiErrors[e.response?.data?.error]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    const blockUser = async () => {
        try {
            setLoading(true)
            await socialUserAPI.blockUser(data.user.id)
            
            setModalTitle(t('Success'))
            if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been blocked'))
            else setModalText(`${data.user.username} ${t('has been blocked')}`)
            openModal()
            await onUpdate()
            setUserManipulation(true)
        } catch (e) {
            setModalTitle(t('Error'))
            if (e.status === 404) {
                setModalText(t('User not found'))
            }
            setModalText(t(apiErrors[e.response?.data?.error]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    const unblockUser = async () => {
        try {
            setLoading(true)
            await socialUserAPI.unblockUser(data.user.id)
            setModalTitle(t('Success'))
            if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been unblocked'))
            else setModalText(`${data.user.username} ${t('has been unblocked')}`)
            openModal()
            await onUpdate()
            setUserManipulation(true)
        } catch (e) {
            setModalTitle(t('Error'))
            if (e.status === 404) {
                setModalText(t('User not found'))
            }
            setModalText(t(apiErrors[e.response?.data?.error]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    return (
        <>
            <div className='flex w-full h-auto gap-9 4k:gap-14'>
                <div className={`${styles.Avatar} aspect-square 2xl:w-1/5 lg:w-1/6 2k:w-1/6 md:w-1/6 sm:w-1/5 mobile:w-1/3 lg:h-full rounded-3xl`}>
                    <div 
                        className='flex items-center justify-center w-full aspect-square
                            Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'
                    >
                        <div className='flex items-center justify-center w-1/2 aspect-square'>
                            <svg width="100%" height="100%" viewBox="0 0 16 19" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <defs/>
                                <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                            </svg>
                        </div>
                    </div>
                </div>

                <div
                    className='lg:text-base 2k:text-2xl 4k:text-4xl
                        md:text-xl sm:text-lg mobile:text-sm flex flex-col gap-2 2k:gap-3 4k:gap-5 sm:gap-3'
                >
                    <div
                        className='font-bold lg:text-2xl 2k:text-3xl 4k:text-5xl
                            md:text-3xl sm:text-2xl mobile:text-xl'
                    >
                        { data.user.firstname } { data.user.lastname }
                    </div>
                    <div>{ t('Username') }: { data.user.username }</div>
                    <div>Email: { data.user.email }</div>
                    { 
                        data.user.phone && (
                            <div>{ t('Phone number') }: { data.user.phone }</div>
                        )    
                    }
                    { 
                        data.user.birthdate && (
                            <div>{ t('Birthdate') }: { new Date(data.user.birthdate).toLocaleDateString() }</div>
                        )    
                    }
                </div>
            </div>
        

            <div className='flex flex-col gap-7 mt-2 2k:mt-4 4k:mt-6 2k:gap-10 4k:gap-14'>
                {
                    !data.friend_status && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={() => addFriend('request')}
                            disabled={isLoading}
                        >
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Add as a friend')
                                )  
                            }
                        </Button>
                    )
                }
                {
                    data.friend_status !== 'blocked' && data.friend_status === 'incoming request' && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={() => addFriend('add')}
                            disabled={isLoading}
                        >
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Accept friend request')
                                )  
                            }
                        </Button>
                    )
                }
                {
                    data.friend_status !== 'blocked' && data.friend_status === 'incoming request' && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={() => removeFriend('decline')}
                            disabled={isLoading}
                        >
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Decline friend request')
                                )  
                            }
                        </Button>
                    )
                }
                {
                    data.friend_status !== 'blocked' && data.friend_status === 'outgoing request' && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={() => removeFriend('cancel')}
                            disabled={isLoading}
                        >
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Cancel friend request')
                                )  
                            }
                        </Button>
                    )
                }
                {
                    data.friend_status !== 'blocked' && data.friend_status === 'accepted' && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={() => removeFriend('remove')}
                            disabled={isLoading}
                        >
                            
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Unfriend')
                                )  
                            }
                        </Button>
                    )
                }
                {
                    data.friend_status !== 'blocked' && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={blockUser}
                            disabled={isLoading}
                        >
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Block')
                                )
                            }
                            
                        </Button>
                    )
                }
                {
                    data.friend_status === 'blocked' && (
                        <Button
                            className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            onClick={unblockUser}
                            disabled={isLoading}
                        >
                            {
                                isLoading ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : (
                                    t('Unblock')
                                )
                            }
                        </Button>
                    )
                }
            </div>
        </>
    )
}

export default SocialUserInfo
