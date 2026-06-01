import { createRouter, createWebHashHistory } from 'vue-router'
import MatchStart from '../views/MatchStart.vue'
import MatchControl from '../views/MatchControl.vue'

const routes = [
  {
    path: '/',
    name: 'MatchStart',
    component: MatchStart
  },
  {
    path: '/control/:id',
    name: 'MatchControl',
    component: MatchControl
  }
]

const router = createRouter({
  // Use hash history for Electron to avoid issues with file:// protocol or 404s on reload
  history: createWebHashHistory(),
  routes
})

export default router
