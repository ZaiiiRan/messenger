/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable react/prop-types */
import { MainWidget } from '../../../shared/ui/MainWidget'
import { useTranslation } from 'react-i18next'
import { Input } from '../../../shared/ui/Input'
import { useEffect, useState, useRef, useCallback } from 'react'
import { ShortUser, ShortUserSkeleton } from '../../../entities/ShortUser'
import { useModal } from '../../../features/modal'
import { apiErrors } from '../../../shared/api'
import axios from 'axios'
import { SocialUser } from '../../../entities/SocialUser'

const UserList = ({ key, goBack, title, fetchFunction, minSearchLength = 0, checkAfterUpdate = true }) => {
    const { t } = useTranslation('friendsFeature')
    const limit = 10
    const [offset, setOffset] = useState(0)
    const [data, setData] = useState([])
    const [isFetching, setFetching] = useState(false)
    const [search, setSearch] = useState('')
    const [end, setEnd] = useState(false)
    const { openModal, setModalTitle, setModalText } = useModal()
    const [selectedUser, setSelectedUser] = useState(null)
    const [userManipulation, setUserManipulation] = useState(false)

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
        } catch (e) {
            if (e.status === 404) {
                setEnd(true)
            } else {
                setModalTitle(t('Error'))
                setModalText(apiErrors[e.response?.data?.error] || t('Internal server error'))
                openModal()
            }
        } finally {
            setFetching(false)
        }

        return () => {
            source.cancel("Operation canceled due to new request")
        }
    }

    const lastSearchRef = useRef()

    useEffect(() => {
        const trimmedSearch = search.trim()
        if (trimmedSearch === lastSearchRef.current) return
        lastSearchRef.current = trimmedSearch

        setData([])
        setOffset(0)
        setEnd(false)
        loadUsers(trimmedSearch, 0, false)
    }, [search])

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

    const observerRef = useRef()
    const lastUserRef = useCallback((node) => {
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
            {
                !selectedUser ? (
                    <MainWidget key={key} title={ t(title) } goBack={ goBack }>
                        <div className='flex flex-col items-center'>
                            <Input 
                                className='px-3 py-2 2k:px-4 2k:py-3 4k:px-6 4k:py-5 rounded-lg 
                                    md:text-lg mobile:text-sm lg:text-sm 2xl:text-lg 2k:text-2xl 4k:text-4xl sm:w-2/3 mobile:w-full lg:w-full 2xl:w-2/3'
                                placeholder={ t('Username or email') }
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                            />
                        </div>
                        <div className='scrollContainer flex flex-col overflow-y-scroll gap-3 w-full box-border'>
                            {
                                data.map((user, index) => (
                                    <div
                                        key={user.user_id}
                                        ref={index === data.length - 1 ? lastUserRef : null}
                                    >
                                        <ShortUser 
                                            lastname={user.lastname}
                                            firstname={user.firstname}
                                            username={user.username}
                                            isActivated={user.is_activated}
                                            isBanned={user.is_banned}
                                            isDeleted={user.is_deleted}
                                            onClick={() => setSelectedUser(user)}
                                        />
                                    </div>
                                ))
                            }
                            {isFetching && (
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
                    </MainWidget>
                ) : (
                        <MainWidget title={ selectedUser.username } goBack={ () => setSelectedUser(null) }>
                            <SocialUser id={selectedUser.user_id} onError={() => setSelectedUser(null)} setUserManipulation={checkAfterUpdate ? setUserManipulation : () => {}} />
                        </MainWidget>
                    
                )
            }
        </>
    )
}

export default UserList