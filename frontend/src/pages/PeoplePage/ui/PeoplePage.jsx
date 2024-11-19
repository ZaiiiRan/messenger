import { motion, AnimatePresence } from 'framer-motion'
import { useEffect, useState } from 'react'
import { shortUsersFetching } from '../../../entities/ShortUser'
import { PeopleMenu, PeopleListWidget, UserWidget } from '../../../widgets/people'

const PeoplePage = () => {
    const [selected, setSelected] = useState(null)
    const [selectedUser, setSelectedUser] = useState(null)
    const [userManipulation, setUserManipulation] = useState(false)

    const open = (optionGroup) => {
        setSelected(optionGroup)
    }

    const goBack = () => {
        setSelected(null)
    }

    useEffect(() => {
        setUserManipulation(false)
        setSelectedUser(null)
    }, [selected])

    return (
        <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1}}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
            className='w-full h-full flex relative lg:gap-10 xl:gap-12 2xl:gap-14 2k:gap-24 4k:gap-36'
        >
            <PeopleMenu open={open}/>
            <AnimatePresence mode='wait'>
            {
                selected &&(
                    <motion.div 
                    initial={{ opacity: 0, x: -500 }}
                    animate={{opacity: 1, x: 0 }}
                    exit={{ opacity: 0, x: -500 }}
                    transition={{ duration: 0.3 }}
                    key={selected} 
                    className='relative lg:sticky lg:w-[55%] lg:h-full lg:z-10 mobile:w-full_screen mobile:absolute mobile:top-0 mobile:left-0 mobile:h-full mobile:z-20 portrait:w-full_screen portrait:absolute portrait:top-0 portrait:left-0 portrait:h-full portrait:z-20'>

                {
                    selected === 'search_friends' && (
                        <PeopleListWidget
                            initialAnimation={null}
                            animation={null}
                            exitAnimation={null}
                            key="search_friends"
                            goBack={goBack}
                            title={'Find friends'}
                            fetchFunction={shortUsersFetching.fetchShortUser}
                            minSearchLength={4}
                            setSelectedUser={setSelectedUser}
                            className='w-full h-full'
                        />
                    )
                }
                {
                    selected === 'your_friends' && (
                        <PeopleListWidget
                            initialAnimation={null}
                            animation={null}
                            exitAnimation={null}
                            key="your_friends"
                            goBack={goBack}
                            title={'Your friends'}
                            fetchFunction={shortUsersFetching.fetchFriends}
                            setSelectedUser={setSelectedUser}
                            className='w-full h-full'
                            selectedUser={selectedUser}
                            userManipulation={userManipulation}
                            setUserManipulation={setUserManipulation}
                        />
                    )
                }
                {
                    selected === 'incoming_requests' && (
                        <PeopleListWidget
                            initialAnimation={null}
                            animation={null}
                            exitAnimation={null}
                            key="incoming_requests"
                            goBack={goBack}
                            title={'Incoming requests'}
                            fetchFunction={shortUsersFetching.fetchIncomingFriendRequests}
                            setSelectedUser={setSelectedUser}
                            className='w-full h-full'
                            selectedUser={selectedUser}
                            userManipulation={userManipulation}
                            setUserManipulation={setUserManipulation}
                        />
                    )
                }
                {
                    selected === 'outgoing_requests' && (
                        <PeopleListWidget
                            initialAnimation={null}
                            animation={null}
                            exitAnimation={null}
                            key="outgoing_requests"
                            goBack={goBack}
                            title={'Outgoing requests'}
                            fetchFunction={shortUsersFetching.fetchOutgoingFriendRequests}
                            setSelectedUser={setSelectedUser}
                            className='w-full h-full'
                            selectedUser={selectedUser}
                            userManipulation={userManipulation}
                            setUserManipulation={setUserManipulation}
                        />
                    )
                }
                {
                    selected === 'black_list' && (
                        <PeopleListWidget
                            initialAnimation={null}
                            animation={null}
                            exitAnimation={null}
                            key="black_list"
                            goBack={goBack}
                            title={'Black list'}
                            fetchFunction={shortUsersFetching.fetchBlackList}
                            setSelectedUser={setSelectedUser}
                            className='w-full h-full'
                            selectedUser={selectedUser}
                            userManipulation={userManipulation}
                            setUserManipulation={setUserManipulation}
                        />
                    )
                }

            <AnimatePresence mode='wait'>
                {
                    selectedUser && (
                        <UserWidget 
                            selectedUser={selectedUser}
                            setUserManipulation={setUserManipulation}
                            onLoadError={() => setSelectedUser(null)}
                            goBack={() => setSelectedUser(null)}
                            className=' w-full h-full absolute top-0 left-0 z-30'
                        />
                    )
                }
            </AnimatePresence>
            </motion.div>
                )

            }
            
            </AnimatePresence>
        </motion.div>
    )
}

export default PeoplePage
