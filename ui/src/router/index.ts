import Vue from 'vue'
import VueRouter, { RouteConfig, RouterOptions } from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes: RouteConfig[] = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/accounts',
    name: 'accounts',
    component: () => import(/* webpackChunkName: "accounts" */ '../views/Accounts.vue'),
    children: [
      {
        path: '/:id',
        component: () => import(/* webpackChunkName: "account_detail" */ '../views/Account.vue'),
        children: [
          {
            path: '/edit'
          }
        ]
      }
    ]
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
