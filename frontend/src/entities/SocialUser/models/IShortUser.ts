interface IShortUser {
    userId: number | string,
    username: string,
    firstname: string,
    lastname: string,
    isDeleted: boolean,
    isBanned: boolean,
    isActivated: boolean
}

export default IShortUser