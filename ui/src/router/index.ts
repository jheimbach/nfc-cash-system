import Vue from 'vue'
import VueRouter, { RouterOptions } from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/accounts',
    name: 'accounts',
    component: () => import(/* webpackChunkName: "accounts" */ '../views/Accounts.vue')
  },
  {
    path: '/groups',
    name: 'groups',
    component: () => import(/* webpackChunkName: "groups" */ '../views/Groups.vue')
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: () => import(/* webpackChunkName: "transactions" */ '../views/Transactions.vue')
  },
  {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const routerOptions: RouterOptions = {
  mode: 'history',
  base: process.env.BASE_URL,
  routes: routes
}

const router = (linkClass: string) => {
  routerOptions.linkActiveClass = linkClass
  return new VueRouter(routerOptions)
}

export default router
