import { api } from '../../../shared/api'

export default class Activation {
    static async activate(code) {
        return api.post('/auth/activate', { code })
    }

    static async resend() {
        return api.get('/auth/resend')
    }
}
