interface IMessage {
    id: string | number,
    memberId: string | number,
    chatId: string | number,
    content: string,
    sentAt: Date,
    lastEdited: string | number | Date | null | undefined
}

export default IMessage