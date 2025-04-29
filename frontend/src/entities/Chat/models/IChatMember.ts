interface IChatMember {
    userId: number,
    role: string,
    isRemoved: boolean,
    isLeft: boolean,
    addedBy: number
}

export default IChatMember