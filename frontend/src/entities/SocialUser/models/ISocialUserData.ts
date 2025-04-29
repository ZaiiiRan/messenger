import IShortUser from "./IShortUser"

interface ISocialUserData extends IShortUser {
    email: string,
    phone?: string | null,
    birthdate?: string | Date | number | null,
}

export default ISocialUserData