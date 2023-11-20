import OpenActivityView from '@/views/OpenActivityView.vue'
import ActivityView from '@/views/legacy/ActivityView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: OpenActivityView
    },
    {
      path: '/legacy',
      name: 'legacy',
      component: ActivityView
    },
    {
      path: '/components/elevationgraphplot',
      name: 'elevationgraphplot',
      component: () => import('@/views/example/ElevationGraph.vue')
    },
    {
      path: '/components/elevationgraphview',
      name: 'elevationgraph',
      component: () => import('@/views/example/ElevationGraphView.vue')
    }
  ]
})

export default router
