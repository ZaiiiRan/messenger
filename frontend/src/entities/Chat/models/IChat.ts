import IChatInfo from './IChatInfo'
import IChatMember from './IChatMember'

interface IChat {
    chat: IChatInfo,
    members: IChatMember[],
    you: IChatMember,
    lastMessage?: any
}

export default IChat