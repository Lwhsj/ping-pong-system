import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '../layout/MainLayout.vue'
import LiveScore from '../views/LiveScore.vue'
import MatchHistory from '../views/MatchHistory.vue'
import MatchDetail from '../views/MatchDetail.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: MainLayout,
      redirect: '/matches',
      children: [
        {
          path: 'matches',
          name: 'matches',
          component: MatchHistory
        },
        {
          path: 'scoreboard',
          name: 'scoreboard',
          component: LiveScore
        },
        {
          path: 'match-detail',
          name: 'match-detail',
          component: MatchDetail
        }
      ]
    }
  ]
})

export default router
