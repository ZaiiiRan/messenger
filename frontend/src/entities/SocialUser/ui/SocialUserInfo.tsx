import styles from './SocialUser.module.css'
import { useTranslation } from 'react-i18next'
import { Button } from '../../../shared/ui/Button'
import { addFriend, removeFriend, blockUser, unblockUser } from '../api/SocialUserFetching'
import { useModal } from '../../../features/modal'
import { apiErrors } from '../../../shared/api'
import { useState } from 'react'
import { Loader } from '../../../shared/ui/Loader'
import { AxiosError } from 'axios'
import { ApiErrorsKey } from '../../../shared/api'
import ISocialUser from '../models/ISocialUser'

interface SocialUserInfoProps {
    data: ISocialUser,
    onUpdate: (data: any) => void,
    setUserManipulation: React.Dispatch<React.SetStateAction<boolean>>,
    onMessageClick?: () => void
}

const SocialUserInfo: React.FC<SocialUserInfoProps> = ({ data, onUpdate, setUserManipulation, onMessageClick  }) => {
    const { t } = useTranslation('socialUser')
    const { openModal, setModalTitle, setModalText } = useModal()
    const [isLoading, setLoading] = useState(false)

    const addFriendAction = async (action: string) => {
        try {
            setLoading(true)
            const response = await addFriend(data.user.userId)
            setModalTitle(t('Success'))
            if (action === 'request') {
                setModalText(t('Friend request sent'))
            } else if (action === 'add') {
                if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been added as a friend'))
                else setModalText(`${data.user.username} ${t('has been added as a friend')}`)
            }
            
            openModal()
            onUpdate(response.data.user)
            setUserManipulation(true)
        } catch (e: any) {
            setModalTitle(t('Error'))
            if (e instanceof AxiosError && e.status === 404) {
                setModalText(t('User not found'))
            }

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    const removeFriendAction = async (action: string) => {
        try {
            setLoading(true)
            const response = await removeFriend(data.user.userId)
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
            onUpdate(response.data.user)
            setUserManipulation(true)
        } catch (e: any) {
            setModalTitle(t('Error'))
            if (e instanceof AxiosError && e.status === 404) {
                setModalText(t('User not found'))
            }

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    const blockUserAction = async () => {
        try {
            setLoading(true)
            const response = await blockUser(data.user.userId)
            
            setModalTitle(t('Success'))
            if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been blocked'))
            else setModalText(`${data.user.username} ${t('has been blocked')}`)
            openModal()
            onUpdate(response.data.user)
            setUserManipulation(true)
        } catch (e: any) {
            setModalTitle(t('Error'))
            if (e instanceof AxiosError && e.status === 404) {
                setModalText(t('User not found'))
            }

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
            openModal()
        } finally {
            setLoading(false)
        }
    }

    const unblockUserAction = async () => {
        try {
            setLoading(true)
            const response = await unblockUser(data.user.userId)
            setModalTitle(t('Success'))
            if (data.user.username.length > 15) setModalText(t('The user') + ' ' + t('has been unblocked'))
            else setModalText(`${data.user.username} ${t('has been unblocked')}`)
            openModal()
            onUpdate(response.data.user)
            setUserManipulation(true)
        } catch (e: any) {
            setModalTitle(t('Error'))
            if (e instanceof AxiosError && e.status === 404) {
                setModalText(t('User not found'))
            }

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
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
                    {
                        data.friendStatus !== 'blocked by target' && data.friendStatus !== 'blocked' && (
                            <div>Email: { data.user.email }</div>
                        )
                    }
                    
                    { 
                        data.friendStatus !== 'blocked by target' && data.friendStatus !== 'blocked' && data.user.phone && (
                            <div>{ t('Phone number') }: { data.user.phone }</div>
                        )    
                    }
                    { 
                        data.friendStatus !== 'blocked by target' && data.friendStatus !== 'blocked' && data.user.birthdate && (
                            <div>{ t('Birthdate') }: { new Date(data.user.birthdate).toLocaleDateString() }</div>
                        )    
                    }
                    { 
                        (data.friendStatus === 'blocked by target' || data.friendStatus === 'blocked') && (
                            <div>{t('Viewing information about this person is restricted')}</div>
                        )    
                    }
                </div>
            </div>
        

            <div className='flex justify-between'>
                <div className='flex flex-col gap-7 mt-2 2k:mt-4 4k:mt-6 2k:gap-10 4k:gap-14'>
                    {
                        !data.friendStatus && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={() => addFriendAction('request')}
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
                        data.friendStatus !== 'blocked' && data.friendStatus !== 'blocked by target' && data.friendStatus === 'incoming request' && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={() => addFriendAction('add')}
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
                        data.friendStatus !== 'blocked' && data.friendStatus !== 'blocked by target' && data.friendStatus === 'incoming request' && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={() => removeFriendAction('decline')}
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
                        data.friendStatus !== 'blocked' && data.friendStatus !== 'blocked by target' && data.friendStatus === 'outgoing request' && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={() => removeFriendAction('cancel')}
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
                        data.friendStatus !== 'blocked' && data.friendStatus !== 'blocked by target' && data.friendStatus === 'accepted' && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={() => removeFriendAction('remove')}
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
                        data.friendStatus !== 'blocked' && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={blockUserAction}
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
                        data.friendStatus === 'blocked' && (
                            <Button
                                className='h-14 flex items-center justify-center 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                    rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                onClick={unblockUserAction}
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
                
                {
                    onMessageClick && (
                        <div className='flex flex-col gap-7 mt-2 2k:mt-4 4k:mt-6 2k:gap-10 4k:gap-14'>
                            {   
                                data.friendStatus !== 'blocked' && data.friendStatus !== 'blocked by target' && (
                                    <Button
                                        className='h-14 aspect-square flex items-center justify-center 2k:h-20 4k:h-32
                                            rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                        onClick={onMessageClick}
                                        disabled={isLoading}
                                    >
                                        <div className='h-1/2 aspect-square flex items-center justify-center'>
                                            <svg viewBox="0 0 45.2986 43.9343" fill="none" xmlns="http://www.w3.org/2000/svg" >
                                                <defs/>
                                                <path id="Vector" d="M38.53 6.46C34.83 2.84 29.95 0.59 24.72 0.1C19.5 -0.4 14.26 0.89 9.92 3.75C5.57 6.6 2.38 10.84 0.91 15.73C-0.57 20.62 -0.24 25.86 1.83 30.54C2.05 30.98 2.12 31.47 2.03 31.94L0.05 41.21C-0.03 41.56 -0.01 41.93 0.09 42.28C0.2 42.62 0.39 42.94 0.66 43.2C0.87 43.41 1.13 43.57 1.41 43.68C1.7 43.79 2 43.84 2.3 43.83L2.75 43.83L12.4 41.95C12.89 41.89 13.39 41.96 13.84 42.15C18.66 44.16 24.06 44.48 29.09 43.05C34.13 41.61 38.49 38.52 41.43 34.29C44.37 30.07 45.7 24.98 45.19 19.91C44.68 14.84 42.36 10.09 38.64 6.5L38.53 6.46ZM13.57 24.13C13.13 24.13 12.69 24 12.32 23.76C11.95 23.52 11.66 23.18 11.49 22.78C11.32 22.38 11.28 21.94 11.36 21.51C11.45 21.09 11.66 20.7 11.98 20.39C12.3 20.08 12.7 19.88 13.13 19.79C13.57 19.71 14.02 19.75 14.44 19.92C14.85 20.08 15.2 20.36 15.45 20.72C15.7 21.08 15.83 21.51 15.83 21.94C15.83 22.52 15.59 23.08 15.17 23.49C14.75 23.9 14.17 24.13 13.57 24.13ZM22.59 24.13C22.14 24.13 21.71 24 21.34 23.76C20.97 23.52 20.68 23.18 20.51 22.78C20.34 22.38 20.29 21.94 20.38 21.51C20.47 21.09 20.68 20.7 21 20.39C21.31 20.08 21.71 19.88 22.15 19.79C22.59 19.71 23.04 19.75 23.45 19.92C23.86 20.08 24.22 20.36 24.46 20.72C24.71 21.08 24.84 21.51 24.84 21.94C24.84 22.52 24.61 23.08 24.18 23.49C23.76 23.9 23.19 24.13 22.59 24.13ZM31.61 24.13C31.16 24.13 30.73 24 30.35 23.76C29.98 23.52 29.69 23.18 29.52 22.78C29.35 22.38 29.31 21.94 29.4 21.51C29.48 21.09 29.7 20.7 30.01 20.39C30.33 20.08 30.73 19.88 31.17 19.79C31.6 19.71 32.06 19.75 32.47 19.92C32.88 20.08 33.23 20.36 33.48 20.72C33.73 21.08 33.86 21.51 33.86 21.94C33.86 22.52 33.62 23.08 33.2 23.49C32.78 23.9 32.2 24.13 31.61 24.13Z" fill="#FFFFFF" fillOpacity="1.000000" fillRule="nonzero"/>
                                            </svg>
                                        </div>
                                    </Button>
                                )
                            }
                        </div>
                    )
                }
                
            </div>
        </>
    )
}

export default SocialUserInfo
