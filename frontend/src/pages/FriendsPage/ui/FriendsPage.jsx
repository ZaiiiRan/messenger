import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { FriendsMenu } from '../../../features/friends'
import { UserList } from '../../../features/friends'
import { shortUsersFetching } from '../../../entities/ShortUser'


const FriendsPage = () => {
    const [selected, setSelected] = useState(null)

    const open = (optionGroup) => {
        setSelected(optionGroup)
    }

    const goBack = () => {
        setSelected(null)
    }

    return (
        <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1}}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
            className='w-full h-full flex relative lg:gap-10 xl:gap-12 2xl:gap-14 2k:gap-24 4k:gap-36'
        >
            <FriendsMenu open={open}/>

            <AnimatePresence mode='wait'>
                {
                    selected === 'search_friends' && (
                        <UserList key="search_friends" goBack={goBack} title={'Find friends'} fetchFunction={shortUsersFetching.fetchShortUser} minSearchLength={4} checkAfterUpdate={false} />
                    )
                }
                {
                    selected === 'your_friends' && (
                        <UserList key="your_friends" goBack={goBack} title={'Your friends'} fetchFunction={shortUsersFetching.fetchFriends} />
                    )
                }
                {
                    selected === 'incoming_requests' && (
                        <UserList key="incoming_requests" goBack={goBack} title={'Incoming requests'} fetchFunction={shortUsersFetching.fetchIncomingFriendRequests} />
                    )
                }
                {
                    selected === 'outgoing_requests' && (
                        <UserList key="outgoing_requests" goBack={goBack} title={'Outgoing requests'} fetchFunction={shortUsersFetching.fetchOutgoingFriendRequests} />
                    )
                }
                {
                    selected === 'black_list' && (
                        <UserList key="black_list" goBack={goBack} title={'Black list'} fetchFunction={shortUsersFetching.fetchBlackList} />
                    )
                }
            </AnimatePresence>
        </motion.div>
    )
}

export default FriendsPage
