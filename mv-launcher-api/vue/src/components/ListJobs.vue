<template>
  <div class="job-list">
    <div class="header-section">
      <h2>Jobs</h2>
      <div class="queue-selector">
        <label for="queue">Select Queue:</label>
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

    <div v-if="error" class="error-message">{{ error }}</div>
    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="!selectedQueue" class="no-queue-selected">Please select a queue to view jobs</div>
    <div v-else-if="filteredAndSortedJobs.length === 0" class="no-jobs">No jobs found in the selected queue</div>
    <div v-else class="table-outer-container">
      <div class="table-container">
        <table class="jobs-table">
          <thead>
            <tr>
              <th @click="sortBy('name')" class="sortable">
                Name
                <span v-if="sortKey === 'name'" class="sort-indicator">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th @click="sortBy('status')" class="sortable">
                Status
                <span v-if="sortKey === 'status'" class="sort-indicator">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th @click="sortBy('created_at')" class="sortable">
                Created
                <span v-if="sortKey === 'created_at'" class="sort-indicator">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th @click="sortBy('duration')" class="sortable">
                Duration
                <span v-if="sortKey === 'duration'" class="sort-indicator">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="job in filteredAndSortedJobs" :key="job.id">
              <td>{{ job.name }}</td>
              <td>
                <span :class="['status-badge', job.status ? job.status.toLowerCase() : '']">
                  {{ job.status }}
                </span>
              </td>
              <td>{{ formatDate(job.created_at) }}</td>
              <td>{{ formatDuration(job.duration) }}</td>
              <td>
                <button @click="showJobDetails(job)" class="view-btn">View Details</button>
                <button v-if="job.status === 'SUCCEEDED'" @click="downloadNextflowLog(job)" class="download-btn">Download Log</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Job Details Modal -->
    <div v-if="selectedJob" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ selectedJob.name }}</h3>
          <button class="close-button" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="job-info-section">
            <h4>Job Information</h4>
            <div class="info-grid">
              <div class="info-item">
                <strong>Status:</strong>
                <span :class="['status-badge', selectedJob.status ? selectedJob.status.toLowerCase() : '']">{{ selectedJob.status }}</span>
              </div>
              <div class="info-item">
                <strong>Job ID:</strong>
                <span>{{ selectedJob.id }}</span>
              </div>
              <div class="info-item">
                <strong>Batch Job ID:</strong>
                <span>{{ selectedJob.batch_job_id }}</span>
              </div>
              <div class="info-item">
                <strong>Job Definition:</strong>
                <span>{{ selectedJob.job_definition }}</span>
              </div>
              <div class="info-item">
                <strong>Job Queue:</strong>
                <span>{{ selectedJob.job_queue }}</span>
              </div>
              <div class="info-item">
                <strong>Attempts:</strong>
                <span>{{ selectedJob.attempts }}</span>
              </div>
            </div>
          </div>

          <div class="job-timing-section">
            <h4>Timing Information</h4>
            <div class="info-grid">
              <div class="info-item">
                <strong>Created:</strong>
                <span>{{ formatDate(selectedJob.created_at) }}</span>
              </div>
              <div class="info-item">
                <strong>Started:</strong>
                <span>{{ formatDate(selectedJob.started_at) }}</span>
              </div>
              <div class="info-item">
                <strong>Stopped:</strong>
                <span>{{ formatDate(selectedJob.stopped_at) }}</span>
              </div>
              <div class="info-item">
                <strong>Duration:</strong>
                <span>{{ formatDuration(selectedJob.duration) }}</span>
              </div>
            </div>
          </div>

          <div class="job-resources-section">
            <h4>Resource Information</h4>
            <div class="info-grid">
              <div class="info-item">
                <strong>Memory:</strong>
                <span>{{ selectedJob.memory }} MB</span>
              </div>
              <div class="info-item">
                <strong>vCPUs:</strong>
                <span>{{ selectedJob.vcpus }}</span>
              </div>
            </div>
          </div>

          <div class="job-status-section">
            <h4>Status Information</h4>
            <div class="info-grid">
              <div class="info-item">
                <strong>Exit Code:</strong>
                <span>{{ selectedJob.exit_code }}</span>
              </div>
              <div class="info-item">
                <strong>Status Reason:</strong>
                <span>{{ selectedJob.status_reason }}</span>
              </div>
              <div class="info-item">
                <strong>Container Reason:</strong>
                <span>{{ selectedJob.container_reason }}</span>
              </div>
            </div>
          </div>

          <div class="job-logs-section">
            <h4>Container Logs</h4>
            <div class="logs-container">
              <div v-if="loadingLogs" class="loading-logs">Loading logs...</div>
              <div v-else-if="logs.length === 0" class="no-logs">No logs available</div>
              <div v-else class="logs-content">
                <div v-for="(log, index) in logs" :key="index" class="log-entry">
                  <span class="log-timestamp">{{ formatDate(log.timestamp) }}</span>
                  <span class="log-message">{{ log.message }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

export default {
  name: 'ListJobs',
  setup() {
    const jobs = ref([])
    const loading = ref(true)
    const selectedQueue = ref('')
    const queues = ref([])
    const error = ref(null)
    const autoRefresh = ref(false)
    const refreshInterval = ref(null)
    const sortKey = ref('created_at')
    const sortOrder = ref('desc')
    const searchQuery = ref('')
    const selectedJob = ref(null)
    const logs = ref([])
    const loadingLogs = ref(false)

    const fetchQueues = async () => {
      loading.value = true
      error.value = null
      try {
        const response = await axios.get('/v1/batch/queues')
        if (response.data && response.data.queues) {
          queues.value = response.data.queues
          if (queues.value.length > 0) {
            selectedQueue.value = queues.value[0].name
            await fetchJobs()
          } else {
            error.value = 'No queues available'
          }
        } else {
          error.value = 'Invalid response from server'
        }
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load queues'
      } finally {
        loading.value = false
      }
    }

    const fetchJobs = async () => {
      if (!selectedQueue.value) {
        jobs.value = []
        return
      }
      loading.value = true
      error.value = null
      try {
        const response = await axios.get(`/v1/jobs?queue=${encodeURIComponent(selectedQueue.value)}`)
        jobs.value = response.data
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load jobs'
        jobs.value = []
      } finally {
        loading.value = false
      }
    }

    const showJobDetails = async (job) => {
      selectedJob.value = job
      loadingLogs.value = true
      logs.value = []
      try {
        const response = await axios.get(`/v1/jobs/${job.id}/logs`)
        logs.value = response.data.cloudwatch_logs || []
      } catch (error) {
        logs.value = []
      } finally {
        loadingLogs.value = false
      }
    }

    const closeModal = () => {
      selectedJob.value = null
      logs.value = []
    }

    const downloadNextflowLog = async (job) => {
      try {
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
        alert('Failed to download log')
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleString()
    }

    const formatDuration = (seconds) => {
      if (!seconds) return 'N/A'
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      const secs = seconds % 60
      return `${hours}h ${minutes}m ${secs}s`
    }

    const sortBy = (key) => {
      if (sortKey.value === key) {
        sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
      } else {
        sortKey.value = key
        sortOrder.value = 'desc'
      }
    }

    const filteredAndSortedJobs = computed(() => {
      let filtered = jobs.value
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(job => job.name && job.name.toLowerCase().includes(query))
      }
      return filtered.sort((a, b) => {
        let aVal = a[sortKey.value]
        let bVal = b[sortKey.value]
        if (sortKey.value === 'created_at') {
          aVal = new Date(aVal).getTime()
          bVal = new Date(bVal).getTime()
        }
        if (typeof aVal === 'string') {
          return sortOrder.value === 'asc' ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal)
        }
        return sortOrder.value === 'asc' ? aVal - bVal : bVal - aVal
      })
    })

    const toggleAutoRefresh = () => {
      if (autoRefresh.value) {
        startAutoRefresh()
      } else {
        stopAutoRefresh()
      }
    }

    const startAutoRefresh = () => {
      stopAutoRefresh()
      refreshInterval.value = setInterval(fetchJobs, 60000)
    }

    const stopAutoRefresh = () => {
      if (refreshInterval.value) {
        clearInterval(refreshInterval.value)
        refreshInterval.value = null
      }
    }

    onMounted(async () => {
      await fetchQueues()
    })

    return {
      jobs,
      loading,
      selectedQueue,
      queues,
      error,
      autoRefresh,
      sortKey,
      sortOrder,
      searchQuery,
      filteredAndSortedJobs,
      selectedJob,
      logs,
      loadingLogs,
      formatDate,
      formatDuration,
      showJobDetails,
      closeModal,
      downloadNextflowLog,
      fetchJobs,
      sortBy,
      toggleAutoRefresh
    }
  }
}
</script>

<style scoped>
.job-list {
  padding: 2rem;
  background-color: #f8f9fa;
  min-height: 100vh;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.queue-selector {
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 300px;
}

.queue-selector label {
  font-weight: 500;
  color: #4a5568;
  white-space: nowrap;
}

.queue-select {
  padding: 0.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background-color: white;
  min-width: 250px;
  font-size: 1rem;
  color: #2d3748;
  cursor: pointer;
}

.queue-select:disabled {
  background-color: #f7fafc;
  cursor: not-allowed;
}

.queue-select:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 2px rgba(0, 102, 204, 0.1);
}

.search-section {
  flex: 1;
  margin: 0 20px;
}

.search-input {
  width: 100%;
  max-width: 300px;
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.search-input:focus {
  border-color: #0066cc;
  box-shadow: 0 0 0 2px rgba(0, 102, 204, 0.1);
}

.refresh-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.btn-refresh {
  padding: 0.5rem 1rem;
  background-color: #0066cc;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.btn-refresh:hover {
  background-color: #0052a3;
}

.btn-refresh:disabled {
  background-color: #cbd5e0;
  cursor: not-allowed;
}

.auto-refresh {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.auto-refresh label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: #4a5568;
  cursor: pointer;
}

.error-message {
  background-color: #fff5f5;
  color: #c53030;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  text-align: center;
}

.loading, .no-queue-selected, .no-jobs {
  text-align: center;
  padding: 2rem;
  color: #4a5568;
  font-size: 1.1rem;
}

.table-outer-container {
  background-color: white;
  border-radius: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.table-container {
  overflow-x: auto;
}

.jobs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.95rem;
}

.jobs-table th,
.jobs-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.jobs-table th {
  background-color: #f8fafc;
  font-weight: 600;
  color: #4a5568;
  white-space: nowrap;
}

.jobs-table th.sortable {
  cursor: pointer;
  user-select: none;
}

.jobs-table th.sortable:hover {
  background-color: #edf2f7;
}

.sort-indicator {
  margin-left: 0.5rem;
  color: #0066cc;
}

.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  text-transform: capitalize;
}

.status-badge.running {
  background-color: #e6fffa;
  color: #2c7a7b;
}

.status-badge.succeeded {
  background-color: #f0fff4;
  color: #2f855a;
}

.status-badge.failed {
  background-color: #fff5f5;
  color: #c53030;
}

.view-btn {
  padding: 0.5rem 1rem;
  background-color: #0066cc;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.view-btn:hover {
  background-color: #0052a3;
}

.download-btn {
  padding: 0.5rem 1rem;
  background-color: #2d3748;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  cursor: pointer;
  margin-left: 0.5rem;
  transition: background-color 0.2s ease;
}

.download-btn:hover {
  background-color: #1a202c;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 20px;
  position: relative;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  padding: 5px;
}

.close-button:hover {
  color: #333;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
}

.logs-container {
  background: #f5f5f5;
  border-radius: 4px;
  padding: 10px;
  max-height: 300px;
  overflow-y: auto;
  margin-top: 10px;
}

.log-entry {
  display: flex;
  gap: 10px;
  padding: 5px 0;
  border-bottom: 1px solid #eee;
  font-family: monospace;
  font-size: 0.9em;
}

.log-timestamp {
  color: #666;
  min-width: 200px;
}

.log-message {
  flex: 1;
  word-break: break-word;
}

.loading-logs, .no-logs {
  text-align: center;
  padding: 20px;
  color: #666;
}
</style>
  