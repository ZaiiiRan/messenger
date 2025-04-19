import axios, { AxiosError } from 'axios'
import { transformKeysToCamelCase } from '../../utils/transformKeysToCamelCase'
import { transformKeysToSnakeCase } from '../../utils/transformKeysToSnakeCase'

export const API_URL = import.meta.env.VITE_API_URL

const api = axios.create({
    withCredentials: true,
    baseURL: API_URL
})

api.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
    if (config.data) {
        config.data = transformKeysToSnakeCase(config.data)
    }
    return config
})

api.interceptors.response.use((response) => {
    response.data = transformKeysToCamelCase(response.data)
    return response
}, (async (error) => {
    const originalRequest = error.config
    if (error.response.status === 401 && error.config && !error.config._isRetry) {
        originalRequest._isRetry = true
        try {
            const response = await axios.get(`${API_URL}/auth/refresh`, { withCredentials: true })
            localStorage.setItem('token', response.data.accessToken)
            return api.request(originalRequest)
        } catch (e) {
            if (e instanceof AxiosError && e.status === 401) {
                localStorage.removeItem('token')
            }
            console.log(e)
        }
    }
    throw error
}))

export default api