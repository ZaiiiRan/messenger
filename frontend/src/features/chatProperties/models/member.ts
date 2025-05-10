import { IChatMember } from '../../../entities/Chat'
import { IShortUser } from '../../../entities/SocialUser'

interface Member extends IChatMember {
    user: IShortUser
}

export default Member