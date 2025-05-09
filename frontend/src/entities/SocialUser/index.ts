import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'socialUser', en)
i18n.addResourceBundle('ru', 'socialUser', ru)

import SocialUser from './ui/SocialUser'
import ISocialUser from './models/ISocialUser'
import ISocialUserData from './models/ISocialUserData'
import IShortUser from './models/IShortUser'
import { fetchShortUser, fetchFriends, fetchIncomingFriendRequests, fetchOutgoingFriendRequests, fetchBlackList, fetchFriendsAreNotChatting } from './api/ShortUsersFetching'
import { fetchSocialUserForStore, fetchSocialUsersForStore } from './api/SocialUserFetching'
import useShortUserStore from './hook/useShortUserStore'
import ShortUserSkeleton from './ui/ShortUserSkeleton'
import ShortUser from './ui/ShortUser'
import shortUserStore from './store/ShortUserStore'
import SocialUserDialog from './ui/SocialUserDialog'

export { SocialUser, ShortUser, ShortUserSkeleton, useShortUserStore, 
    fetchShortUser, fetchFriends, fetchIncomingFriendRequests, fetchOutgoingFriendRequests, fetchBlackList, fetchSocialUserForStore, fetchSocialUsersForStore,
    shortUserStore, SocialUserDialog, fetchFriendsAreNotChatting
}
export type { ISocialUser, ISocialUserData, IShortUser }