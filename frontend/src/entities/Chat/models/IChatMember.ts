interface IChatMember {
    userId: number | string,
    role: string,
    isRemoved: boolean,
    isLeft: boolean,
    addedBy: number
}

export default IChatMember