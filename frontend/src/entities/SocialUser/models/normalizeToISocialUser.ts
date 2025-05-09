import ISocialUser from './ISocialUser'
import ISocialUserData from './ISocialUserData'

function normalizeToIShortUser(data: any): ISocialUser {
    const socialUserData = data.user as ISocialUserData
    return {
        user: socialUserData,
        friendStatus: data.friendStatus
    }
}

export default normalizeToIShortUser