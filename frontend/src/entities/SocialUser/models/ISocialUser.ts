import ISocialUserData from './ISocialUserData'

interface ISocialUser {
    user: ISocialUserData,
    friendStatus: string | null
}

export default ISocialUser