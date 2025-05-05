import IChatInfo from './IChatInfo'
import IChatMember from './IChatMember'
import IMessage from './IMessage'

interface IChat {
    chat: IChatInfo,
    members: IChatMember[],
    you: IChatMember,
    lastMessage?: IMessage | null,
    messages: IMessage[],
}

export default IChat