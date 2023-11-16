import { createRouter, createWebHistory } from 'vue-router'
import ActivityView from '../views/ActivityView.vue'
import OpenActivityView from '../views/OpenActivityView.vue'

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
      path: '/activity',
      name: 'activity',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/ActivityView.vue')
    },
    {
      path: '/elevationgraph',
      name: 'elevationgraph',
      component: () => import('../views/ElevationGraph.vue')
    },
    {
      path: '/components/elevationgraph',
      name: 'graph',
      component: () => import('../views/example/ElevationGraphView.vue')
    }
  ]
})

export default router
