import axios, { AxiosError } from 'axios'
import { transformKeysToCamelCase } from '../../utils/transformKeysToCamelCase'
import { transformKeysToSnakeCase } from '../../utils/transformKeysToSnakeCase'

export const API_URL = import.meta.env.VITE_API_URL

const api = axios.create({
    withCredentials: true,
    baseURL: API_URL
})

let isRefreshing = false
let failedQueue: { resolve: (token: string) => void; reject: (error: any) => void }[] = []

const processQueue = (error: any, token: string | null = null) => {
    failedQueue.forEach((prom) => {
        if (error) {
            prom.reject(error)
        } else if (token) {
            prom.resolve(token)
        }
    })
    failedQueue = []
}

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
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

        if (isRefreshing) {
            return new Promise((resolve, reject) => {
                failedQueue.push({ resolve, reject })
            })
                .then((token) => {
                    originalRequest.headers.Authorization = `Bearer ${token}`
                    return api.request(originalRequest)
                })
                .catch((err) => Promise.reject(err))
        }

        isRefreshing = true

        try {
            const response = await axios.get(`${API_URL}/auth/refresh`, { withCredentials: true })
            const newToken = response.data.accessToken
            localStorage.setItem('token', newToken)

            originalRequest.headers.Authorization = `Bearer ${newToken}`

            processQueue(null, newToken)
            return api.request(originalRequest)
        } catch (e) {
            processQueue(e, null)
            if (e instanceof AxiosError && e.status === 401) {
                localStorage.removeItem('token')
                window.location.href = '/login'
            }
            throw e
        } finally {
            isRefreshing = false
        }
    }
    throw error
}))

export default api