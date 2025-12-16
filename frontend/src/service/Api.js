import axios from 'axios'
import { useAuthStore } from '@/stores/auth.js'

let csrfToken = null

export function setCsrfToken(token) {
    csrfToken = token
}

export default () => {
    const store = useAuthStore()
    const instance = axios.create({
        baseURL: process.env.VUE_APP_API_URL,
        withCredentials: true,
    });

    // CSRF-Header vor jedem Request setzen
    instance.interceptors.request.use((config) => {
        if (csrfToken && !['get', 'head', 'options'].includes(config.method)) {
            config.headers['X-CSRF-Token'] = csrfToken
        }
        return config
    })

    instance.interceptors.response.use(function (response) {
        return response
    }, function (error) {
        if (error.response.status === 401) {
            store.logout();
        }

        return Promise.reject(error);
    })
    return instance
}
