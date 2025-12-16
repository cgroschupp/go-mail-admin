import { defineStore } from 'pinia'
import Client from "../service/Client";

export const useAuthStore = defineStore('auth', {
  persist: true,
  state: () => {
    return {
      isAuthenticated: false,
      loginFailed: false,
    }
  },
  actions: {
    async login(username, password) {
      await Client.login(username, password).then((res) => {
        if (res.status == 200) {
          this.isAuthenticated = true;
          this.loginFailed = false
          this.router.push({ name: 'Home' })
        } else {
          this.loginFailed = true
        }
      }).catch((res) => {
        this.loginFailed = true
      });
    },
    async logout() {
      if (this.isAuthenticated) {
        Client.logout().then(() => {
          this.isAuthenticated = false
        }).catch(() => {
          this.isAuthenticated = false
        });
      }
      if (this.router.name != "Login") {
        this.router.push({ name: 'Login' })
      }
    }
  },
})