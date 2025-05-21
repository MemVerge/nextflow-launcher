<template>
    <div class="job-list-container">
        <h1 class="title">Job List</h1>
        <div class="job-list-wrapper">
            <div class="header">
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
            </div>
            <div v-if="error" class="error-message">{{ error }}</div>
            <div v-if="loading" class="loading">Loading...</div>
            <div v-else-if="!selectedQueue" class="no-queue-selected">Please select a queue to view jobs</div>
            <div v-else-if="jobs.length === 0" class="no-jobs">No jobs found in the selected queue</div>
            <div v-else class="jobs-grid">
                <div v-for="job in jobs" :key="job.id" class="job-card">
                    <h3 class="job-name">{{ job.name }}</h3>
                    <div class="job-details">
                        <p><strong>Pipeline:</strong> {{ job.pipeline }}</p>
                        <p><strong>Status:</strong> <span :class="job.status">{{ job.status }}</span></p>
                        <p><strong>Created:</strong> {{ formatDate(job.created_at) }}</p>
                    </div>
                    <button @click="viewJobDetails(job.id)" class="view-btn">View Details</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'

export default {
    name: 'JobList',
    setup() {
        const jobs = ref([])
        const loading = ref(true)
        const selectedQueue = ref('')
        const queues = ref([])
        const error = ref(null)

        const fetchQueues = async () => {
            loading.value = true
            error.value = null
            try {
                console.log('Fetching queues...')
                const response = await axios.get('/v1/batch/queues')
                console.log('Queues response:', response.data)
                
                if (response.data && response.data.queues) {
                    queues.value = response.data.queues
                    console.log('Available queues:', queues.value)
                    
                    // Select first queue if available
                    if (queues.value.length > 0) {
                        selectedQueue.value = queues.value[0].name
                        console.log('Selected queue:', selectedQueue.value)
                        await fetchJobs()
                    } else {
                        console.log('No queues available')
                        error.value = 'No queues available'
                    }
                } else {
                    console.error('Invalid queues response:', response.data)
                    error.value = 'Invalid response from server'
                }
            } catch (err) {
                console.error('Error fetching queues:', err)
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
                console.log('Fetching jobs for queue:', selectedQueue.value)
                const response = await axios.get(`/v1/jobs?queue=${encodeURIComponent(selectedQueue.value)}`)
                console.log('Jobs response:', response.data)
                jobs.value = response.data
            } catch (err) {
                console.error('Error fetching jobs:', err)
                error.value = err.response?.data?.error || 'Failed to load jobs'
                jobs.value = []
            } finally {
                loading.value = false
            }
        }

        const formatDate = (dateString) => {
            return new Date(dateString).toLocaleString()
        }

        const viewJobDetails = (jobId) => {
            console.log('Viewing job details for:', jobId)
        }

        onMounted(async () => {
            console.log('Component mounted, fetching queues...')
            await fetchQueues()
        })

        return {
            jobs,
            loading,
            selectedQueue,
            queues,
            error,
            formatDate,
            viewJobDetails,
            fetchJobs,
            fetchQueues
        }
    }
}
</script>

<style scoped>
.job-list-container {
    width: 100%;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
    background-color: #f8f9fa;
}

.title {
    color: #0066cc;
    text-align: center;
    margin-bottom: 2rem;
    font-size: 2.5rem;
    font-weight: bold;
    width: 100%;
}

.job-list-wrapper {
    width: 100%;
    max-width: 1200px;
    background-color: white;
    padding: 2rem;
    border-radius: 16px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header {
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

.error-message {
    background-color: #fff5f5;
    color: #c53030;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
    text-align: center;
}

.no-queue-selected {
    text-align: center;
    padding: 2rem;
    color: #4a5568;
    font-size: 1.1rem;
    background-color: #f7fafc;
    border-radius: 8px;
    border: 1px dashed #cbd5e0;
}

.loading, .no-jobs {
    text-align: center;
    padding: 2rem;
    color: #4a5568;
    font-size: 1.1rem;
}

.jobs-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-top: 1rem;
}

.job-card {
    background-color: #f8f9fa;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.job-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.job-name {
    color: #2c3e50;
    font-size: 1.2rem;
    margin-bottom: 1rem;
}

.job-details {
    margin-bottom: 1rem;
}

.job-details p {
    margin: 0.5rem 0;
    color: #4a5568;
}

.status {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.9rem;
    font-weight: 500;
}

.status.running {
    background-color: #e6fffa;
    color: #2c7a7b;
}

.status.completed {
    background-color: #f0fff4;
    color: #2f855a;
}

.status.failed {
    background-color: #fff5f5;
    color: #c53030;
}

.view-btn {
    width: 100%;
    padding: 0.75rem;
    background-color: #0066cc;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.view-btn:hover {
    background-color: #0052a3;
}
</style> 