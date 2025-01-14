import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthTokenView from '../views/AuthTokenView.vue'
import TokenView from '@/views/TokenView.vue'
import LogoutView from '@/views/LogoutView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/auth/token/:token',
      name: 'auth_token',
      component: AuthTokenView,
    },
    {
      path: '/tokens',
      name: 'token',
      component: TokenView
    },
    {
      path: '/logout',
      name: 'logout',
      component: LogoutView
    }
  ],
})

export default router
