import { IChat, IChatInfo, IChatMember, IMessage } from "../../../entities/Chat"

function normalizeToIChat(object: any): IChat {
    const chatInfo = object.chat as IChatInfo
    const lastMessage = object.lastMessage as IMessage
    const members: IChatMember[] = []

    object.members.forEach((value: any) => {
        const member: IChatMember = {
            userId: value.user.userId,
            role: value.role,
            isRemoved: value.isRemoved,
            isLeft: value.isLeft,
            addedBy: value.addedBy
        } 
        members.push(member)
    })

    const you: IChatMember = {
        userId: object.you.user.userId,
        role: object.you.role,
        isRemoved: object.you.isRemoved,
        isLeft: object.you.isLeft,
        addedBy: object.you.addedBy
    } 

    const chat: IChat = {
        chat: chatInfo,
        lastMessage: lastMessage,
        members: members,
        you: you
    }

    return chat
}

export default normalizeToIChat