interface ISocialUserData {
    userId: number,
    username: string,
    firstname: string,
    lastname: string,
    isDeleted: boolean,
    isBanned: boolean,
    isActivated: boolean,
    email: string,
    phone?: string | null,
    birthdate?: string | Date | number | null,
}

export default ISocialUserData