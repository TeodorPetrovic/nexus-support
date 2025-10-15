import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Page1 from '@/views/Page1.vue'
import Page2 from '@/views/Page2.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/discovery'
  },
  {
    path: '/discovery',
    name: 'Discovery',
    component: Page1
  },
  {
    path: '/screenshots',
    name: 'Screenshots',
    component: Page2
  },
  {
    path: '/stream/:ip/:hostname',
    name: 'Stream',
    component: () => import('@/views/StreamView.vue'),
    props: true
  },
  // Keep old routes for compatibility
  {
    path: '/page1',
    redirect: '/discovery'
  },
  {
    path: '/page2',
    redirect: '/screenshots'
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
