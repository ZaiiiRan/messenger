import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'

export default class Auth {
    static async login(login: string, password: string): Promise<AxiosResponse<any, any>> {
        return api.post('/auth/login', { login, password })
    }

    static async register(username: string, email: string, password: string, firstname: string, 
    lastname: string, phone: string | undefined | null, birthdate: string | undefined | null): Promise<AxiosResponse<any, any>> {
        return api.post('/auth/register', { username, email, password, firstname, lastname, phone, birthdate })
    }

    static async logout(): Promise<AxiosResponse<any, any>> {
        return api.delete('/auth/logout')
    }

    static async refresh(): Promise<AxiosResponse<any, any>> {
        return api.get('/auth/refresh')
    }
}
