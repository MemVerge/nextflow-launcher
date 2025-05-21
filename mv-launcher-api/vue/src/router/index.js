import { createRouter, createWebHistory } from 'vue-router'
import CreateJob from '../components/CreateJob.vue'
import ListJobs from '../components/ListJobs.vue'
import AWSBatchSetup from '../components/AWSBatchSetup.vue'

const routes = [
    {
      path: '/',
    name: 'CreateJob',
      component: CreateJob
    },
    {
      path: '/jobs',
    name: 'ListJobs',
      component: ListJobs
  },
  {
    path: '/batch-setup',
    name: 'AWSBatchSetup',
    component: AWSBatchSetup
    }
  ]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
