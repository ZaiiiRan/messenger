import { useEffect, useState } from 'react'
import socialUserAPI from '../api/SocialUserFetching'
import SocialUserInfoSkeleton from './SocialUserInfoSkeleton'
import SocialUserInfo from './SocialUserInfo'
import { useModal } from '../../../features/modal'
import { useTranslation } from 'react-i18next'
import { apiErrors } from '../../../shared/api'
import { AxiosError } from 'axios'
import { ApiErrorsKey } from '../../../shared/api'
import ISocialUser from '../models/ISocialUser'

interface SocialUserProps {
    id: number,
    onError: () => void,
    setUserManipulation: React.Dispatch<React.SetStateAction<boolean>>,
    onMessageClick: (event: React.MouseEvent<HTMLButtonElement>) => void
}

const SocialUser: React.FC<SocialUserProps> = ({ id, onError, setUserManipulation, onMessageClick }) => {
    const { t } = useTranslation('socialUser')
    const [data, setData] = useState<ISocialUser | null>(null)
    const [isFetching, setFetching] = useState<boolean>(true)
    const { openModal, setModalTitle, setModalText } = useModal()

    const load = async () => {
        try {
            const response = await socialUserAPI.fetch(id)
            setData(response.data)
        } catch (e: any) {
            setModalTitle(t('Error'))
            if (e instanceof AxiosError && e.status === 404) {
                setModalText(t('User not found'))
            }

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
            openModal()
            onError()
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
                        setUserManipulation={setUserManipulation} 
                        onMessageClick={onMessageClick} 
                    />
                )
            }
        </div>
    )
}

export default SocialUser
