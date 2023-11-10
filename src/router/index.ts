import { createRouter, createWebHistory } from 'vue-router'
import ActivityView from '../views/ActivityView.vue'
import ActivityViewV0 from '../views/v0/ActivityView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ActivityView
    },
    {
      path: '/v0',
      name: 'v0',
      component: ActivityViewV0
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
      path: '/components/graph',
      name: 'graph',
      component: () => import('../views/ComponentGraphView.vue')
    }
  ]
})

export default router
