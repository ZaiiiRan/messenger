import { AxiosResponse } from 'axios'
import { api } from '../../../shared/api'

export default class Activation {
    static async activate(code: string): Promise<AxiosResponse<any, any>> {
        return api.patch('/auth/activate', { code })
    }

    static async resend(): Promise<AxiosResponse<any, any>> {
        return api.get('/auth/resend')
    }
}
