<template>
  <div class="job-list">
    <div class="header-section">
      <div class="header-top">
        <h2>Jobs</h2>
        <div class="search-section">
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="Search jobs by name..." 
            class="search-input"
          >
        </div>
        <div class="refresh-controls">
          <button @click="fetchJobs" class="btn-refresh" :disabled="loading">
            <span v-if="loading">Refreshing...</span>
            <span v-else>Refresh</span>
          </button>
          <div class="auto-refresh">
            <label>
              <input type="checkbox" v-model="autoRefresh" @change="toggleAutoRefresh">
              Auto-refresh (every 60s)
            </label>
          </div>
        </div>
      </div>
      <div class="header-bottom">
        <div class="queue-selector">
          <label for="queue">Select Head Node Queue:</label>
          <select 
            id="queue" 
            v-model="selectedQueue" 
            @change="fetchJobs" 
            class="queue-select"
            :disabled="loading"
          >
            <option value="">Select a queue</option>
            <option v-for="queue in queues" :key="queue.name" :value="queue.name">
              {{ queue.name }}
            </option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.header-bottom {
  display: flex;
  align-items: center;
}

.queue-selector {
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 400px;
}

.search-section {
  flex: 1;
  max-width: 400px;
}

.search-input {
  width: 100%;
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}
</style> 

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const showJobDetails = async (job) => {
  selectedJob.value = job
  loadingLogs.value = true
  logs.value = []
  try {
    // Debug: print the job object and batch_job_id
    console.log('showJobDetails called with job:', job)
    console.log('Using batch_job_id for logs:', job.batch_job_id)
    
    if (!job.batch_job_id) {
      console.error('No batch_job_id found in job object')
      logs.value = []
      return
    }
    
    const response = await axios.get(`/v1/jobs/${job.batch_job_id}/logs`)
    logs.value = response.data.cloudwatch_logs || []
  } catch (error) {
    console.error('Error fetching logs:', error)
    logs.value = []
  } finally {
    loadingLogs.value = false
  }
}

const downloadNextflowLog = async (job) => {
  try {
    // Log the job object to debug
    console.log('Job object for download:', job)
    
    // Use the job's ID for getting the log URL
    if (!job.id) {
      console.error('No job ID found in job object')
      alert('No job ID available for this job')
      return
    }
    
    const response = await axios.get(`/v1/jobs/${job.id}/log-url`)
    if (response.data.url) {
      const link = document.createElement('a')
      link.href = response.data.url
      link.download = `nextflow-${job.name}-${job.id}.log`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    } else {
      alert('No download URL available for this job')
    }
  } catch (error) {
    console.error('Error downloading log:', error)
    alert('Failed to download log')
  }
}
</script> 