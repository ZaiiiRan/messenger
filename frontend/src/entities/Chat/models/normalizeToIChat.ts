import { chatStore } from '..'
import IChat from './IChat'
import IChatInfo from './IChatInfo'
import IChatMember from './IChatMember'
import IMessage from './IMessage'
import normalizeToIChatMember from './normalizeToIChatMember'

function normalizeToIChat(object: any): IChat {
    const chatInfo = object.chat as IChatInfo
    const lastMessage = object.lastMessage as IMessage
    const members: IChatMember[] = []

    if (object.members) {
        object.members.forEach((value: any) => {
            const member: IChatMember = normalizeToIChatMember(value)
            members.push(member)
        })
    }

    const you: IChatMember = normalizeToIChatMember(object.you)

    let messages: IMessage[] = []
    const chatCandidate = chatStore.get(chatInfo.id)
    if (chatCandidate) {
        messages = chatCandidate.messages
    } else if (lastMessage) {
        messages.push(lastMessage)
    }

    const chat: IChat = {
        chat: chatInfo,
        lastMessage: lastMessage,
        members: members,
        you: you,
        messages: messages
    }

    return chat
}

export default normalizeToIChat