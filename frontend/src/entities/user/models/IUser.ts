interface IUser {
    userId: number | string,
    username: string,
    email: string,
    isActivated: boolean,
    isBanned: boolean,
    isDeleted: boolean,
    lastname: string,
    firstname: string,
    birthdate?: string | Date | number | null,
    phone?: string | null
}

export default IUser