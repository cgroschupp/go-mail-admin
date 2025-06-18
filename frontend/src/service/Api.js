import axios from 'axios'
import { useAuthStore } from '@/stores/auth.js'

export default () => {
    const store = useAuthStore()
    const instance = axios.create({
        baseURL: process.env.VUE_APP_API_URL,
    });
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
