import { createRouter, createWebHistory } from 'vue-router'
import ActivityView from '../views/ActivityView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
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
      path: '/altitudegraph',
      name: 'altitudegraph',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AltitudeGraph.vue')
    }
  ]
})

export default router
