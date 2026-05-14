import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Domains from "../views/Domains.vue";
import Alias from "../views/Alias.vue";
import AliasEdit from "../views/AliasEdit.vue";
import Accounts from "../views/Accounts.vue";
import AccountEdit from "../views/AccountEdit.vue";
import TLSPolicy from "../views/TLSPolicy.vue";
import TLSPolicyEdit from "../views/TLSPolicyEdit.vue";
import Login from "../views/Login.vue"
import Logout from "../views/Logout.vue";
import { useAuthStore } from '../stores/auth.js'

const routes = [
  {
    path: '/domains',
    name: 'Domains',
    component: Domains
  },
  {
    path: '/alias',
    name: 'Alias',
    component: Alias
  },
  {
    path: '/alias/:id',
    name: 'AliasEdit',
    component: AliasEdit
  },
  {
    path: '/account',
    name: 'Accounts',
    component: Accounts
  },
  {
    path: '/account/new',
    name: 'AccountNew',
    component: AccountEdit
  },
  {
    path: '/account/:id',
    name: 'AccountEdit',
    component: AccountEdit
  },
  {
    path: '/tls',
    name: 'TLS',
    component: TLSPolicy,
  },
  {
    path: '/tls/new',
    name: 'TLSNew',
    component: TLSPolicyEdit
  },
  {
    path: '/tls/:id',
    name: 'TLSEdit',
    component: TLSPolicyEdit
  },
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/logout',
    name: 'Logout',
    component: Logout
  }
]

let router = createRouter({
  history: createWebHistory(),
  routes: routes,
})

router.beforeEach((to, from, next) => {
  const store = useAuthStore()
  const isAuthenticated = store.isAuthenticated
  if (to.name !== 'Login' && !isAuthenticated) next({ name: 'Login' })
  else next()
})

export default router
