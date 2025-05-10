import IChatMember from './IChatMember'

function normalizeToIChatMember(object: any): IChatMember {
    return {
        userId: object.user.userId,
        role: object.role,
        isRemoved: object.isRemoved,
        isLeft: object.isLeft,
        addedBy: object.addedBy
    }
}

export default normalizeToIChatMember