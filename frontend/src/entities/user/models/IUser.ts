interface IUser {
    user_id: number,
    username: string,
    email: string,
    is_activated: boolean,
    is_banned: boolean,
    is_deleted: boolean,
    lastname: string,
    firstname: string,
    birthdate?: string | Date | number | null,
    phone?: string | null
}

export default IUser