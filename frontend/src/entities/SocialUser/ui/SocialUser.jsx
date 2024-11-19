/* eslint-disable react/prop-types */
/* eslint-disable react-hooks/exhaustive-deps */
import { useEffect, useState } from 'react'
import socialUserAPI from '../api/SocialUserFetching'
import SocialUserInfoSkeleton from './SocialUserInfoSkeleton'
import SocialUserInfo from './SocialUserInfo'
import { useModal } from '../../../features/modal'
import { useTranslation } from 'react-i18next'
import { apiErrors } from '../../../shared/api'

const SocialUser = ({ id, onError, setUserManipulation }) => {
    const { t } = useTranslation('socialUser')
    const [data, setData] = useState()
    const [isFetching, setFetching] = useState(true)
    const { openModal, setModalTitle, setModalText } = useModal()

    const load = async () => {
        try {
            const response = await socialUserAPI.fetch(id)
            setData(response.data)
        } catch (e) {
            setModalTitle(t('Error'))
            if (e.status === 404) {
                setModalText(t('User not found'))
            }
            setModalText(t(apiErrors[e.response?.data?.error]) || t('Internal server error'))
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
                ) : (
                    <SocialUserInfo data={data} onUpdate={setData} setUserManipulation={setUserManipulation} />
                )
            }
        </div>
    )
}

export default SocialUser
