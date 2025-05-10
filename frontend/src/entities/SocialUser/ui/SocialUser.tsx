import { useEffect, useState } from 'react'
import { fetchSocialUser } from '../api/SocialUserFetching'
import SocialUserInfoSkeleton from './SocialUserInfoSkeleton'
import SocialUserInfo from './SocialUserInfo'
import { useModal } from '../../../features/modal'
import { useTranslation } from 'react-i18next'
import { apiErrors } from '../../../shared/api'
import { AxiosError } from 'axios'
import { ApiErrorsKey } from '../../../shared/api'
import ISocialUser from '../models/ISocialUser'

interface SocialUserProps {
    id: number | string,
    onError?: () => void,
    setUserManipulation?: React.Dispatch<React.SetStateAction<boolean>>,
    onMessageClick?: () => void
}

const SocialUser: React.FC<SocialUserProps> = ({ id, onError, setUserManipulation, onMessageClick }) => {
    const { t } = useTranslation('socialUser')
    const [data, setData] = useState<ISocialUser | null>(null)
    const [isFetching, setFetching] = useState<boolean>(true)
    const { openModal } = useModal()

    const load = async () => {
        try {
            const response = await fetchSocialUser(id)
            setData(response.data)
        } catch (e: any) {
            if (e instanceof AxiosError && e.status === 404) {
                openModal(t('Error'), t('User not found'))
                onError && onError()
                return
            }

            const errorKey: ApiErrorsKey = e.response?.data?.error
            const errorMsg = t(apiErrors[errorKey]) || t('Internal server error')
            openModal(t('Error'), errorMsg)
            if (onError) onError()
        } finally {
            setFetching(false)
        }
    }

    useEffect(() => {
        setFetching(true)
        load()
    }, [id])

    return (
        <div
            className='flex flex-col gap-8 2k:gap-14 4k:gap-24'
        >
            {
                isFetching ? (
                    <SocialUserInfoSkeleton />
                ) : ( data &&
                    <SocialUserInfo 
                        data={data} 
                        onUpdate={setData} 
                        setUserManipulation={setUserManipulation ? setUserManipulation : () => {}} 
                        onMessageClick={onMessageClick} 
                    />
                )
            }
        </div>
    )
}

export default SocialUser
