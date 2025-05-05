import axios, { AxiosError, AxiosResponse } from 'axios'
import { useTranslation } from 'react-i18next'
import { useModal } from '../../modal'
import { useCallback, useEffect, useRef, useState } from 'react'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import { ShortUser } from '../../../entities/SocialUser'
import { Input } from '../../../shared/ui/Input'

interface UserSelectionProps {
    onSelect: (id: string | number) => void,
    fetchFunction: (search: string, limit: number, offset: number) => Promise<AxiosResponse<any, any>>,
    checkSelected: (userId: string | number) => boolean
}

const LIMIT = 10

const UserSelection: React.FC<UserSelectionProps> = ({ onSelect, fetchFunction, checkSelected }) => {
    const { t } = useTranslation('userSelection')
    const { openModal, setModalText, setModalTitle } = useModal()
    const [users, setUsers] = useState<any[]>([])
    const [isFetching, setFetching] = useState<boolean>(false)
    const [end, setEnd] = useState<boolean>(false)
    const [offset, setOffset] = useState<number>(0)
    const [search, setSearch] = useState<string>('')

    const loadUsers = async (newSearch = search, newOffset = offset, newEnd = end, newLimit = LIMIT) => {
        if (newEnd || isFetching) return
        setFetching(true)

        const source = axios.CancelToken.source()

        try {
            const response = await fetchFunction(newSearch, newLimit, newOffset)
            const newUsers = response.data.users
            if (newUsers.length < LIMIT) setEnd(true)
            setUsers((prevUsers) => [...prevUsers, ...newUsers])
            setOffset((prevOffset) => prevOffset + LIMIT)
        } catch (e: any) {
            if (e instanceof AxiosError && e.status === 404) {
                setEnd(true)
            } else {
                setModalTitle(t('Error'))

                const errorKey: ApiErrorsKey = e.response?.data?.error
                setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
                openModal()
            }
        } finally {
            setFetching(false)
        }

        return () => {
            source.cancel("Operation canceled due to new request")
        }
    }

    const lastSearchRef = useRef<string>()

    useEffect(() => {
        const trimmedSearch = search.trim()
        if (trimmedSearch === lastSearchRef.current) return
        lastSearchRef.current = trimmedSearch

        setUsers([])
        setOffset(0)
        setEnd(false)
        loadUsers(trimmedSearch, 0, false)
    }, [search])

    const observerRef = useRef<IntersectionObserver | null>(null)
    const lastUserRef = useCallback((node: HTMLDivElement | null) => {
        if (isFetching) return
        if (observerRef.current) observerRef.current.disconnect()
        observerRef.current = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && !end) {
                loadUsers(search.trim(), offset, end)
            }
        })
        if (node) observerRef.current.observe(node)
    }, [isFetching])


    return (
        <>
            <div className='h-12 w-full'>
                    <Input
                        placeholder={t('Username or Email')}
                        className='w-full px-2 py-1 2k:px-3 2k:py-2 4k:px-4 4k:py-35 rounded-lg 
                            md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                        value={search}
                        onChange={(e) => setSearch(e.target.value) }
                    />
            </div>

            <div className='h-96 scrollContainer flex flex-col overflow-y-scroll gap-3 w-full box-border'>
                {
                    users.map((user, index) => (
                        <div
                            key={user.userId}
                            ref={index === users.length - 1 ? lastUserRef : null}
                        >
                            <ShortUser 
                                user={user}
                                onClick={() => {onSelect(user.userId)}}
                                isSelectable
                                isSelected={checkSelected(user.userId)}
                            />
                        </div>
                    ))
                }
                {
                    end && !isFetching && users.length === 0 && (
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

export default UserSelection