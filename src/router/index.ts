// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
