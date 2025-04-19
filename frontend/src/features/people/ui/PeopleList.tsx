import { useTranslation } from 'react-i18next'
import { useEffect, useState, useRef, useCallback, SetStateAction, Dispatch } from 'react'
import { ShortUser, ShortUserSkeleton, IShortUser } from '../../../entities/ShortUser'
import { useModal } from '../../../features/modal'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import axios, { AxiosError, AxiosResponse } from 'axios'

interface PeopleListProps {
    search: string,
    fetchFunction: (search: string, limit: number, offset: number) => Promise<AxiosResponse<any, any>>,
    setSelectedUser: Dispatch<SetStateAction<IShortUser | null>>,
    minSearchLength: number, 
    userManipulation: boolean,
    setUserManipulation: Dispatch<SetStateAction<boolean>>,
    selectedUser: IShortUser | null
}

const PeopleList: React.FC<PeopleListProps> = ({ search, fetchFunction, setSelectedUser, minSearchLength = 0, userManipulation, setUserManipulation, selectedUser }) => {
    const { t } = useTranslation('peopleFeature')
    const limit: number = 10
    const [offset, setOffset] = useState<number>(0)
    const [data, setData] = useState<Array<any>>([])
    const [isFetching, setFetching] = useState<boolean>(false)
    const [end, setEnd] = useState<boolean>(false)
    const { openModal, setModalTitle, setModalText } = useModal()

    const loadUsers = async (newSearch = search, newOffset = offset, newEnd = end, newLimit = limit) => {
        if (newEnd || newSearch.length < minSearchLength || isFetching) return
        setFetching(true)

        const source = axios.CancelToken.source()

        try {
            const response = await fetchFunction(newSearch, newLimit, newOffset)
            const newUsers = response.data.users
            if (newUsers.length < limit) setEnd(true)
            setData((prevUsers) => [...prevUsers, ...newUsers])
            setOffset((prevOffset) => prevOffset + limit)
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

        setData([])
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

    useEffect(() => {
        if (!selectedUser && userManipulation) {
            const trimmedSearch = search.trim()
            const newLimit = data.length
            setData([])
            setOffset(0)
            setEnd(false)
            loadUsers(trimmedSearch, 0, false, newLimit)
            setUserManipulation(false)
        }
    }, [userManipulation, selectedUser])

    return (
        <>
            <div className='scrollContainer flex flex-col overflow-y-scroll gap-3 w-full box-border'>
                {
                    data.map((user, index) => (
                        <div
                            key={user.userId}
                            ref={index === data.length - 1 ? lastUserRef : null}
                        >
                            <ShortUser 
                                user={user}
                                onClick={() => setSelectedUser(user)}
                            />
                        </div>
                    ))
                }
                { isFetching && (
                    <>
                        {Array.from({ length: 5 }).map((_, index) => (
                            <ShortUserSkeleton key={index} />
                        ))}
                    </>
                )}
            </div>
            {
                end && !isFetching && data.length === 0 && (
                    <div 
                        className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                            md:text-xl sm:text-lg mobile:text-text-base text-center'
                    >
                        { t('We couldn\'t find anyone') }
                    </div>
                )
            }
        </>
    )
}

export default PeopleList
