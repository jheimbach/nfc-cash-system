import Vue from 'vue'
import VueRouter, { RouteConfig, RouterOptions } from 'vue-router'
import Home from '../views/Home.vue'
import Accounts from '@/views/Accounts.vue'
import AccountDetail from '@/views/Account.vue'
import GroupDetail from '@/views/Group.vue'
import AccountCreate from '@/views/Account/Create.vue'
import AccountCreateSingle from '@/views/Account/Create/Single.vue'
import AccountCreateMultiple from '@/views/Account/Create/Multiple.vue'
import AccountCreateUpload from '@/views/Account/Create/Upload.vue'
import GroupCreate from '@/views/Group/Create.vue'
import GroupCreateSingle from '@/views/Group/Create/Single.vue'
import GroupCreateMultiple from '@/views/Group/Create/Multiple.vue'
import GroupCreateUpload from '@/views/Group/Create/Upload.vue'
import Calculator from '@/views/Calculator.vue'
import App from '@/views/App.vue'
import Login from '@/views/Login.vue'

Vue.use(VueRouter)

const routes: RouteConfig[] = [
  {
    path: '',
    component: App,
    meta: {
      needsAuth: true
    },
    children: [
      {
        path: '/',
        name: 'home',
        component: Home,
        meta: {
          needsAuth: true
        }
      },
      {
        path: '/accounts',
        name: 'accounts',
        component: Accounts,
        meta: {
          needsAuth: true
        }
      },
      {
        path: '/calculator',
        name: 'calculator',
        component: Calculator,
        meta: {
          needsAuth: true
        }
      },
      {
        path: '/accounts/create',
        component: AccountCreate,
        meta: {
          needsAuth: true
        },
        children: [{
          path: '',
          name: 'account_create',
          component: AccountCreateSingle,
          meta: {
            needsAuth: true
          }
        }, {
          path: 'multiple',
          name: 'account_create_multiple',
          component: AccountCreateMultiple,
          meta: {
            needsAuth: true
          }
        }, {
          path: 'upload',
          name: 'account_create_upload',
          component: AccountCreateUpload,
          meta: {
            needsAuth: true
          }
        }]
      },
      {
        path: '/account/:id',
        name: 'account',
        component: AccountDetail,
        meta: {
          needsAuth: true
        }
      },
      {
        path: '/groups',
        name: 'groups',
        component: () => import(/* webpackChunkName: "groups" */ '../views/Groups.vue'),
        meta: {
          needsAuth: true
        }
      },
      {
        path: '/groups/create',
        component: GroupCreate,
        meta: {
          needsAuth: true
        },
        children: [{
          path: '',
          name: 'group_create',
          component: GroupCreateSingle,
          meta: {
            needsAuth: true
          }
        }, {
          path: 'multiple',
          name: 'group_create_multiple',
          component: GroupCreateMultiple,
          meta: {
            needsAuth: true
          }
        }, {
          path: 'upload',
          name: 'group_create_upload',
          component: GroupCreateUpload,
          meta: {
            needsAuth: true
          }
        }]
      },
      {
        path: '/group/:id',
        name: 'group',
        component: GroupDetail,
        meta: {
          needsAuth: true
        }
      },
      {
        path: '/transactions',
        name: 'transactions',
        component: () => import(/* webpackChunkName: "transactions" */ '../views/Transactions.vue'),
        meta: {
          needsAuth: true
        }
      }
    ]
  },
  {
    path: '/login',
    component: Login,
    name: 'login'
  }
]

const routerOptions: RouterOptions = {
  mode: 'history',
  base: process.env.BASE_URL,
  routes: routes
}

const router = new VueRouter(routerOptions)

router.beforeEach((to, from, next) => {
  if (to.meta.needsAuth === true) {
    // @ts-ignore
    const jwt = JSON.parse(localStorage.getItem('jwt'))
    if (jwt === null) {
      next({
        name: 'login'
      })
    } else {
      if (parseInt(jwt.expires_in) < Math.round((new Date()).getTime() / 1000)) {
        next({
          name: 'login',
          query: {
            to: JSON.stringify({ name: to.name, params: to.params })
          }
        })
      }
      next()
    }
  }
  next()
})

export default router
