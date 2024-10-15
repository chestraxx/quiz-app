import { createRouter, createWebHistory } from 'vue-router'
import JoinView from '../views/JoinView.vue'
import QuizView from '../views/QuizView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'join',
      component: JoinView
    },
    {
      path: '/quiz',
      name: 'quiz',
      component: QuizView
    }
  ]
})

export default router
