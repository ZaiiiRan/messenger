interface IShortUser {
    userId: number | string,
    username: string,
    firstname: string,
    lastname: string,
    isDeleted: false,
    isBanned: false,
    isActivated: false
}

export default IShortUser