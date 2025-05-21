<template>
  <div class="batch-setup">
    <div class="header">
      <h2>AWS Batch Setup</h2>
      <p class="description">Configure your AWS Batch compute environment and job queue</p>
    </div>

    <form @submit.prevent="submitForm" class="setup-form">
      <div class="form-section">
        <h3>AWS Settings</h3>
        <div class="form-group">
          <label for="region">AWS Region</label>
          <input
            type="text"
            id="region"
            v-model="form.region"
            class="form-control"
            placeholder="e.g., us-west-2"
            required
          />
          <small class="help-text">Enter your AWS region (e.g., us-west-2, us-east-1, eu-west-1)</small>
        </div>
        <div class="form-group">
          <label for="uniquePrefix">Unique Prefix</label>
          <input
            type="text"
            id="uniquePrefix"
            v-model="form.uniquePrefix"
            class="form-control"
            placeholder="e.g., my-project"
            required
            pattern="[a-z0-9\-]+"
            title="Only lowercase letters, numbers, and hyphens are allowed"
          />
          <small class="help-text">This will be used as a prefix for all AWS resources. Use only lowercase letters, numbers, and hyphens. This field is required.</small>
        </div>
      </div>

      <div class="form-section">
        <h3>Compute Environment Settings</h3>
        <div class="form-group">
          <label for="computeEnvName">Compute Environment Name (Optional)</label>
          <input
            type="text"
            id="computeEnvName"
            v-model="form.computeEnvName"
            class="form-control"
            placeholder="Leave empty for auto-generated name"
          />
          <small class="help-text">If not specified, a name will be generated using the unique prefix.</small>
        </div>
        <div class="form-group">
          <label for="instanceTypes">Instance Types</label>
          <input
            type="text"
            id="instanceTypes"
            v-model="form.instanceTypes"
            class="form-control"
            placeholder="e.g., t3.micro,t3.small"
            required
          />
          <small class="help-text">Comma-separated list of EC2 instance types (e.g., t3.micro,t3.small)</small>
        </div>
        <div class="form-group">
          <label class="allocation-strategy-label">Allocation Strategy</label>
          <div class="allocation-strategy-options">
            <label class="radio-label">
              <input
                type="radio"
                v-model="form.allocationStrategy"
                value="BEST_FIT"
                class="radio-input"
              />
              <span class="radio-text">Best Fit</span>
            </label>
            <label class="radio-label">
              <input
                type="radio"
                v-model="form.allocationStrategy"
                value="BEST_FIT_PROGRESSIVE"
                class="radio-input"
              />
              <span class="radio-text">Best Fit Progressive</span>
            </label>
          </div>
          <small class="help-text">Best Fit: Selects the instance type that best matches the job requirements. Best Fit Progressive: Selects the instance type that best matches the job requirements and has the lowest cost.</small>
        </div>
        <div class="form-group">
          <label class="instance-type-label">Instance Type</label>
          <div class="instance-type-options">
            <label class="radio-label">
              <input
                type="radio"
                v-model="form.useSpot"
                :value="false"
                class="radio-input"
              />
              <span class="radio-text">On-Demand</span>
            </label>
            <label class="radio-label">
              <input
                type="radio"
                v-model="form.useSpot"
                :value="true"
                class="radio-input"
              />
              <span class="radio-text">Spot</span>
            </label>
          </div>
          <small class="help-text">On-Demand: Pay for compute capacity by the hour or second. Spot: Request unused EC2 instances at up to 90% discount.</small>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label for="minvCpus">Minimum vCPUs</label>
            <input
              type="number"
              id="minvCpus"
              v-model.number="form.minvCpus"
              class="form-control"
              min="0"
              required
            />
          </div>
          <div class="form-group">
            <label for="maxvCpus">Maximum vCPUs</label>
            <input
              type="number"
              id="maxvCpus"
              v-model.number="form.maxvCpus"
              class="form-control"
              min="0"
              required
            />
          </div>
          <div class="form-group">
            <label for="desiredvCpus">Desired vCPUs</label>
            <input
              type="number"
              id="desiredvCpus"
              v-model.number="form.desiredvCpus"
              class="form-control"
              min="0"
              required
            />
          </div>
        </div>
        <div class="form-group">
          <label for="subnetId">Subnet IDs</label>
          <input
            type="text"
            id="subnetId"
            v-model="form.subnetId"
            class="form-control"
            placeholder="e.g., subnet-12345678,subnet-87654321"
            required
            pattern="subnet-[a-z0-9]+(,subnet-[a-z0-9]+)*"
            title="Enter one or more subnet IDs separated by commas (e.g., subnet-12345678,subnet-87654321)"
          />
          <small class="help-text">Enter one or more subnet IDs separated by commas. Multiple subnets are recommended for high availability.</small>
        </div>
        <div class="form-group">
          <label for="securityGroupId">Security Group ID</label>
          <input
            type="text"
            id="securityGroupId"
            v-model="form.securityGroupId"
            class="form-control"
            placeholder="e.g., sg-12345678"
            required
            pattern="sg-[a-z0-9]+"
            title="Security Group ID must start with 'sg-' followed by alphanumeric characters"
          />
        </div>
      </div>

      <div class="form-section">
        <h3>Job Queue Settings</h3>
        <div class="form-group">
          <label for="jobQueueName">Job Queue Name (Optional)</label>
          <input
            type="text"
            id="jobQueueName"
            v-model="form.jobQueueName"
            class="form-control"
            placeholder="Leave empty for auto-generated name"
          />
          <small class="help-text">If not specified, a name will be generated using the unique prefix.</small>
        </div>
      </div>

      <div class="form-actions">
        <router-link to="/" class="cancel-button">Cancel</router-link>
        <button type="submit" class="submit-button" :disabled="isSubmitting">
          {{ isSubmitting ? 'Setting up AWS Batch...' : 'Setup AWS Batch' }}
        </button>
      </div>
    </form>

    <!-- Status Messages -->
    <div v-if="statusMessage" :class="['status-message', statusType]">
      {{ statusMessage }}
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'AWSBatchSetup',
  data() {
    return {
      form: {
        region: 'us-west-2',
        uniquePrefix: '',
        computeEnvName: '',
        jobQueueName: '',
        instanceTypes: 't3.micro',
        minvCpus: null,
        maxvCpus: null,
        desiredvCpus: null,
        subnetId: '',
        securityGroupId: '',
        useSpot: false,
        allocationStrategy: 'BEST_FIT',
        enableMultiQueue: false
      },
      isSubmitting: false,
      statusMessage: '',
      statusType: 'info'
    }
  },
  methods: {
    async submitForm() {
      // Validate required fields
      if (!this.form.region || !this.form.uniquePrefix || !this.form.instanceTypes || 
          !this.form.subnetId || !this.form.securityGroupId) {
        this.statusMessage = 'Please fill in all required fields'
        this.statusType = 'error'
        return
      }

      // Validate vCPU values
      if (this.form.minvCpus === null || this.form.maxvCpus === null || this.form.desiredvCpus === null) {
        this.statusMessage = 'Please specify all vCPU values'
        this.statusType = 'error'
        return
      }

      // Validate unique prefix
      if (!/^[a-z0-9\-]+$/.test(this.form.uniquePrefix)) {
        this.statusMessage = 'Please provide a valid unique prefix (lowercase letters, numbers, and hyphens only)'
        this.statusType = 'error'
        return
      }

      this.isSubmitting = true
      this.statusMessage = 'Starting AWS Batch setup...'
      this.statusType = 'info'

      try {
        // Convert comma-separated instance types to array
        const formData = {
          region: this.form.region,
          unique_prefix: this.form.uniquePrefix,
          compute_env_name: this.form.computeEnvName,
          job_queue_name: this.form.jobQueueName,
          instance_types: this.form.instanceTypes.split(',').map(t => t.trim()),
          min_vcpus: this.form.minvCpus || 0,
          max_vcpus: this.form.maxvCpus || 0,
          desired_vcpus: this.form.desiredvCpus || 0,
          subnet_id: this.form.subnetId,
          security_group_id: this.form.securityGroupId,
          use_spot: this.form.useSpot,
          allocation_strategy: this.form.allocationStrategy,
          enable_multi_queue: this.form.enableMultiQueue
        }

        console.log('Submitting form data:', formData)

        const response = await axios.post('/v1/batch/setup', formData)
        
        this.statusMessage = 'AWS Batch setup completed successfully!'
        this.statusType = 'success'
        
        // Emit success event with the created resources
        this.$emit('setup-complete', response.data.resources)
        
        // Reset form after successful submission
        this.resetForm()
      } catch (error) {
        console.error('Error setting up AWS Batch:', error)
        this.statusMessage = error.response?.data?.error || 'Failed to setup AWS Batch. Please try again.'
        this.statusType = 'error'
      } finally {
        this.isSubmitting = false
      }
    },
    resetForm() {
      this.form = {
        region: 'us-west-2',
        uniquePrefix: '',
        computeEnvName: '',
        jobQueueName: '',
        instanceTypes: 't3.micro',
        minvCpus: null,
        maxvCpus: null,
        desiredvCpus: null,
        subnetId: '',
        securityGroupId: '',
        useSpot: false,
        allocationStrategy: 'BEST_FIT',
        enableMultiQueue: false
      }
    }
  }
}
</script>

<style scoped>
.batch-setup {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  text-align: center;
  margin-bottom: 2rem;
}

.header h2 {
  color: #2c3e50;
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.description {
  color: #666;
  font-size: 1.1rem;
}

.setup-form {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #eee;
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.form-section h3 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
  font-size: 1.3rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: #4a5568;
  font-weight: 500;
}

.help-text {
  display: block;
  margin-top: 0.5rem;
  color: #718096;
  font-size: 0.875rem;
}

.form-control, .form-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.2s ease;
  background-color: #f8f9fa;
}

.form-control:focus, .form-select:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
  background-color: white;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.submit-button, .cancel-button {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  text-align: center;
}

.submit-button {
  background: #0066cc;
  color: white;
  border: none;
  min-width: 150px;
}

.submit-button:hover:not(:disabled) {
  background: #0052a3;
  transform: translateY(-1px);
}

.submit-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.cancel-button {
  background: #f8f9fa;
  color: #4a5568;
  border: 1px solid #e2e8f0;
}

.cancel-button:hover {
  background: #e2e8f0;
}

.status-message {
  margin-top: 1.5rem;
  padding: 1rem;
  border-radius: 8px;
  text-align: center;
  font-weight: 500;
}

.status-message.info {
  background: #e3f2fd;
  color: #1976d2;
}

.status-message.success {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-message.error {
  background: #ffebee;
  color: #c62828;
}

@media (max-width: 768px) {
  .batch-setup {
    padding: 1rem;
  }

  .setup-form {
    padding: 1.5rem;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .form-actions {
    flex-direction: column;
  }

  .submit-button, .cancel-button {
    width: 100%;
  }
}

.allocation-strategy-label,
.instance-type-label {
  display: block;
  margin-bottom: 0.75rem;
  color: #4a5568;
  font-weight: 500;
}

.allocation-strategy-options,
.instance-type-options {
  display: flex;
  gap: 2rem;
  margin-bottom: 0.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.radio-input {
  width: 1.25rem;
  height: 1.25rem;
  margin: 0;
  cursor: pointer;
}

.radio-text {
  font-size: 1rem;
  color: #4a5568;
}

.help-text {
  display: block;
  margin-top: 0.5rem;
  color: #718096;
  font-size: 0.875rem;
  line-height: 1.4;
}
</style> 