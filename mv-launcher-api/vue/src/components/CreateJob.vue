<template>
  <div class="two-step-container">
    <aside class="sidebar">
      <div :class="['sidebar-step', { active: step === 1 }]" @click="goToStep(1)">
        <span class="step-number">1</span>
        AWS Settings
      </div>
      <div :class="['sidebar-step', { active: step === 2, disabled: !canGoToStep2 }]" @click="canGoToStep2 && goToStep(2)">
        <span class="step-number">2</span>
        Job Settings
      </div>
    </aside>
    <div class="step-content">
      <h1 class="title">Memory Machine Batch Nextflow Launcher</h1>
      <h2 class="form-title">Submit Job</h2>
      <div v-if="step === 1" class="step-form">
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
            <label for="accessKeyId">AWS Access Key ID</label>
            <input
              type="text"
              id="accessKeyId"
              v-model="form.accessKeyId"
              class="form-control"
              placeholder="Enter your AWS Access Key ID"
              required
            />
          </div>
        <div class="form-group">
            <label for="secretAccessKey">AWS Secret Access Key</label>
            <input
              type="password"
              id="secretAccessKey"
              v-model="form.secretAccessKey"
              class="form-control"
              placeholder="Enter your AWS Secret Access Key"
              required
            />
        </div>
        <div class="form-group">
            <p class="batch-setup-link">
              Need to set up AWS Batch? <router-link to="/batch-setup" class="link-button">Configure AWS Batch</router-link>
            </p>
          </div>
        </div>

        <button class="next-btn" :disabled="!canGoToStep2" @click="goToStep(2)">Next</button>
      </div>
      <div v-else-if="step === 2" class="step-form">
        <div class="form-group">
          <label class="form-label">Job Name</label>
          <input v-model="name" type="text" class="form-control" placeholder="Enter job name">
        </div>
        <div class="form-group">
          <label class="form-label">nf-core Pipeline</label>
          <select v-model="pipeline" class="form-select">
            <option value="none"></option>
            <option value="nf-core/airrflow">nf-core/airrflow</option>
            <option value="nf-core/demultiplex">nf-core/demultiplex</option>
            <option value="nf-core/hlatyping">nf-core/hlatyping</option>
            <option value="nf-core/fastquorum">nf-core/fastquorum</option>
            <option value="nf-core/oncoanalyser">nf-core/oncoanalyser</option>
            <option value="nf-core/rnavar">nf-core/rnavar</option>
            <option value="nf-core/ampliseq">nf-core/ampliseq</option>
            <option value="nf-core/mag">nf-core/mag</option>
            <option value="nf-core/variantbenchmarking">nf-core/variantbenchmarking</option>
            <option value="nf-core/genomeassembler">nf-core/genomeassembler</option>
            <option value="nf-core/scnanoseq">nf-core/scnanoseq</option>
            <option value="nf-core/eager">nf-core/eager</option>
            <option value="nf-core/metatdenovo">nf-core/metatdenovo</option>
            <option value="nf-core/taxprofiler">nf-core/taxprofiler</option>
            <option value="nf-core/scrnaseq">nf-core/scrnaseq</option>
            <option value="nf-core/pacvar">nf-core/pacvar</option>
            <option value="nf-core/molkart">nf-core/molkart</option>
            <option value="nf-core/funcscan">nf-core/funcscan</option>
            <option value="nf-core/raredisease">nf-core/raredisease</option>
            <option value="nf-core/multiplesequencealign">nf-core/multiplesequencealign</option>
            <option value="nf-core/phyloplace">nf-core/phyloplace</option>
            <option value="nf-core/sarek">nf-core/sarek</option>
            <option value="nf-core/pangenome">nf-core/pangenome</option>
            <option value="nf-core/proteinfamilies">nf-core/proteinfamilies</option>
            <option value="nf-core/pairgenomealign">nf-core/pairgenomealign</option>
            <option value="nf-core/fastqrepair">nf-core/fastqrepair</option>
            <option value="nf-core/riboseq">nf-core/riboseq</option>
            <option value="nf-core/drugresponseeval">nf-core/drugresponseeval</option>
            <option value="nf-core/denovotranscript">nf-core/denovotranscript</option>
            <option value="nf-core/nanostring">nf-core/nanostring</option>
            <option value="nf-core/pixelator">nf-core/pixelator</option>
            <option value="nf-core/rangeland">nf-core/rangeland</option>
            <option value="nf-core/references">nf-core/references</option>
            <option value="nf-core/rnaseq">nf-core/rnaseq</option>
            <option value="nf-core/methylseq">nf-core/methylseq</option>
            <option value="nf-core/phaseimpute">nf-core/phaseimpute</option>
            <option value="nf-core/metapep">nf-core/metapep</option>
            <option value="nf-core/bacass">nf-core/bacass</option>
            <option value="nf-core/detaxizer">nf-core/detaxizer</option>
            <option value="nf-core/crisprseq">nf-core/crisprseq</option>
            <option value="nf-core/demo">nf-core/demo</option>
            <option value="nf-core/smrnaseq">nf-core/smrnaseq</option>
            <option value="nf-core/chipseq">nf-core/chipseq</option>
            <option value="nf-core/isoseq">nf-core/isoseq</option>
            <option value="nf-core/proteinfold">nf-core/proteinfold</option>
            <option value="nf-core/reportho">nf-core/reportho</option>
            <option value="nf-core/mhcquant">nf-core/mhcquant</option>
            <option value="nf-core/callingcards">nf-core/callingcards</option>
            <option value="nf-core/epitopeprediction">nf-core/epitopeprediction</option>
            <option value="nf-core/rnasplice">nf-core/rnasplice</option>
            <option value="nf-core/differentialabundance">nf-core/differentialabundance</option>
            <option value="nf-core/bamtofastq">nf-core/bamtofastq</option>
            <option value="nf-core/readsimulator">nf-core/readsimulator</option>
            <option value="nf-core/metaboigniter">nf-core/metaboigniter</option>
            <option value="nf-core/rnafusion">nf-core/rnafusion</option>
            <option value="nf-core/nascent">nf-core/nascent</option>
            <option value="nf-core/fetchngs">nf-core/fetchngs</option>
            <option value="nf-core/circdna">nf-core/circdna</option>
            <option value="nf-core/cutandrun">nf-core/cutandrun</option>
            <option value="nf-core/atacseq">nf-core/atacseq</option>
            <option value="nf-core/viralintegration">nf-core/viralintegration</option>
            <option value="nf-core/marsseq">nf-core/marsseq</option>
            <option value="nf-core/hic">nf-core/hic</option>
            <option value="nf-core/hgtseq">nf-core/hgtseq</option>
            <option value="nf-core/viralrecon">nf-core/viralrecon</option>
            <option value="nf-core/nanoseq">nf-core/nanoseq</option>
            <option value="nf-core/coproid">nf-core/coproid</option>
            <option value="nf-core/mnaseseq">nf-core/mnaseseq</option>
            <option value="nf-core/hicar">nf-core/hicar</option>
            <option value="nf-core/bactmap">nf-core/bactmap</option>
            <option value="nf-core/diaproteomics">nf-core/diaproteomics</option>
            <option value="nf-core/clipseq">nf-core/clipseq</option>
            <option value="nf-core/pgdb">nf-core/pgdb</option>
            <option value="nf-core/dualrnaseq">nf-core/dualrnaseq</option>
            <option value="nf-core/cageseq">nf-core/cageseq</option>
            <option value="nf-core/kmermaid">nf-core/kmermaid</option>
            <option value="nf-core/proteomicslfq">nf-core/proteomicslfq</option>
            <option value="nf-core/imcyto">nf-core/imcyto</option>
            <option value="nf-core/slamseq">nf-core/slamseq</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Profile</label>
          <input v-model="profile" type="text" class="form-control" placeholder="Enter profile name (e.g., test, test_full)">
        </div>
        <div class="form-group">
          <label class="form-label">Input Directory</label>
          <div class="input-group">
            <select v-model="input_bucket" class="form-select" :disabled="isTestProfile">
              <option v-for="bucket in inputBucket" :key="bucket" :value="bucket">{{ bucket }}</option>
            </select>
            <input v-model="input_dir" type="text" class="form-control" placeholder="Enter path to input directory" :disabled="isTestProfile">
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">Working Directory</label>
          <div class="input-group">
            <select v-model="work_bucket" class="form-select">
              <option v-for="bucket in workBucket" :key="bucket" :value="bucket">{{ bucket }}</option>
            </select>
            <input v-model="work_dir" type="text" class="form-control" placeholder="Enter path to working directory">
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">Result Directory</label>
          <div class="input-group">
            <select v-model="result_bucket" class="form-select">
              <option v-for="bucket in resultBucket" :key="bucket" :value="bucket">{{ bucket }}</option>
            </select>
            <input v-model="result_dir" type="text" class="form-control" placeholder="Enter path to result directory">
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">Log Bucket</label>
          <select v-model="log_bucket" class="form-select">
            <option v-for="bucket in logBucket" :key="bucket" :value="bucket">{{ bucket }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Head Node Job Queue</label>
          <select v-model="head_node_queue" class="form-select">
            <option v-for="queue in jobQueues" :key="queue" :value="queue">{{ queue }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Task Job Queue</label>
          <select v-model="task_queue" class="form-select">
            <option v-for="queue in jobQueues" :key="queue" :value="queue">{{ queue }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Memory (GB)</label>
          <input v-model="memory" type="text" class="form-control" placeholder="Enter memory in GB (e.g., 20G)">
        </div>
        <div class="form-group">
          <label class="form-label">Max Retries</label>
          <input v-model="max_retries" type="number" class="form-control" placeholder="Enter max retries (default: 5)">
        </div>
        <div class="form-group">
          <label class="form-label">Additional Configuration</label>
          <textarea v-model="additional_config" class="form-control" rows="4" placeholder="Enter additional Nextflow configuration (optional)"></textarea>
        </div>
        <div class="form-actions">
          <button class="back-btn" @click="goToStep(1)">Back</button>
          <button @click="submitForm" class="submit-btn">Submit Job</button>
        </div>
        <div class="alert" :class="{'alert-success': responseMessage.startsWith('Success'), 'alert-danger': responseMessage.startsWith('Error')}" v-if="responseMessage">
          {{ responseMessage }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import axios from 'axios'
import { onMounted } from 'vue'

export default {
  name: 'PutForm',
  setup() {
    const step = ref(1)
    const name = ref('')
    const profile = ref('')
    const pipeline = ref('')
    const input_bucket = ref('')
    const input_dir = ref('')
    const work_bucket = ref('')
    const work_dir = ref('')
    const result_bucket = ref('')
    const result_dir = ref('')
    const log_bucket = ref('')
    const head_node_queue = ref('')
    const task_queue = ref('')
    const memory = ref('')
    const max_retries = ref(5)
    const additional_config = ref('')
    const responseMessage = ref('')
    const workBucket = ref([])
    const resultBucket = ref([])
    const logBucket = ref([])
    const inputBucket = ref([])
    const jobQueues = ref([])
    const form = ref({
      region: '',
      accessKeyId: '',
      secretAccessKey: '',
      sessionToken: ''
    })

    const isTestProfile = computed(() => {
      return profile.value === 'test' || profile.value === 'test_full'
    })

    const canGoToStep2 = computed(() => form.value.accessKeyId && form.value.secretAccessKey)
    const goToStep = (n) => { if (n === 1 || (n === 2 && canGoToStep2.value)) step.value = n }

    onMounted(async () => {
      try {
        const response = await axios.get('/v1/buckets')
        workBucket.value = response.data.map(bucket => bucket.name)
        resultBucket.value = response.data.map(bucket => bucket.name)
        logBucket.value = response.data.map(bucket => bucket.name)
        inputBucket.value = response.data.map(bucket => bucket.name)
      } catch (error) {}
      try {
        const response = await axios.get('/v1/batch/queues')
        jobQueues.value = response.data.queues.map(queue => queue.name)
      } catch (error) {}
    })

    const submitForm = async () => {
      // Validate input directory requirement
      if (!isTestProfile.value && (!input_bucket.value || !input_dir.value)) {
        responseMessage.value = 'Error: Input directory is required when not using test or test_full profile'
        return
      }

      const body = {
        'name': name.value,
        'pipeline': pipeline.value,
        'work_dir': work_bucket.value + '/' + work_dir.value,
        'result_dir': result_bucket.value + '/' + result_dir.value,
        'head_node_queue': head_node_queue.value,
        'task_queue': task_queue.value,
        'memory': memory.value,
        'max_retries': max_retries.value,
        'additional_config': additional_config.value,
        'aws_access_key': form.value.accessKeyId,
        'aws_secret_key': form.value.secretAccessKey,
        'log_bucket': log_bucket.value
      }
      if (profile.value) {
        body.profile = profile.value
      }
      if (!isTestProfile.value && input_bucket.value && input_dir.value) {
        body.input_dir = input_bucket.value + '/' + input_dir.value
      }
      try {
        const response = await axios.post('/v1/jobs', body)
        responseMessage.value = 'Success: ' + response.data.message
      } catch (error) {
        responseMessage.value = 'Error: ' + error.message
      }
    }

    return {
      step,
      goToStep,
      canGoToStep2,
      name,
      pipeline,
      work_dir,
      profile,
      input_bucket,
      input_dir,
      work_bucket,
      result_dir,
      result_bucket,
      log_bucket,
      head_node_queue,
      task_queue,
      memory,
      max_retries,
      additional_config,
      submitForm,
      workBucket,
      resultBucket,
      logBucket,
      inputBucket,
      responseMessage,
      jobQueues,
      form,
      isTestProfile
    }
  }
}
</script>

<style scoped>
.two-step-container {
  display: flex;
  min-height: 80vh;
  background: #f8f9fa;
}

.sidebar {
  width: 220px;
  background: #fff;
  border-radius: 16px 0 0 16px;
  box-shadow: 2px 0 8px rgba(0,0,0,0.04);
  padding: 2rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
  align-items: flex-start;
  margin: 2rem 0 2rem 2rem;
  min-height: 500px;
}

.sidebar-step {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 1.1rem;
  font-weight: 500;
  color: #4a5568;
  cursor: pointer;
  padding: 0.75rem 1.25rem;
  border-radius: 8px;
  transition: background 0.2s, color 0.2s;
}

.sidebar-step.active {
  background: #e6f0ff;
  color: #0066cc;
  font-weight: 700;
}

.sidebar-step.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.step-number {
  background: #0066cc;
  color: #fff;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  font-weight: 700;
}

.step-content {
  flex: 1;
  padding: 3rem 4rem;
  background: #fff;
  border-radius: 0 16px 16px 0;
  margin: 2rem 2rem 2rem 0;
  min-width: 0;
  box-shadow: 0 4px 6px rgba(0,0,0,0.07);
}

.title {
  color: #0066cc;
  text-align: left;
  margin-bottom: 1rem;
  font-size: 2rem;
  font-weight: bold;
}

.form-title {
  color: #2c3e50;
  text-align: left;
  margin-bottom: 2rem;
  font-size: 1.5rem;
  font-weight: 500;
}

.step-form {
  max-width: 600px;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2c3e50;
  font-weight: 500;
  font-size: 0.95rem;
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

.form-select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  z-index: 1;
  position: relative;
}

.form-select option {
  background-color: white;
  color: black;
  padding: 8px;
}

.input-group {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.input-group .form-select {
  flex: 1;
}

.input-group .form-control {
  flex: 2;
}

.next-btn, .back-btn {
  padding: 0.75rem 2rem;
  background-color: #0066cc;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-right: 1rem;
}

.next-btn:disabled {
  background-color: #b3d1f7;
  cursor: not-allowed;
}

.back-btn {
  background-color: #aaa;
}

.back-btn:hover {
  background-color: #888;
}

.next-btn:hover:not(:disabled) {
  background-color: #0052a3;
  transform: translateY(-1px);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.submit-btn {
  padding: 0.75rem 2rem;
  background-color: #0066cc;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.submit-btn:hover {
  background-color: #0052a3;
  transform: translateY(-1px);
}

.alert {
  margin-top: 1.5rem;
  padding: 1rem;
  border-radius: 8px;
  font-size: 0.95rem;
}

.alert-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-danger {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.form-section {
  margin-top: 2rem;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.section-title {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: #333;
}

.setup-btn {
  background-color: #4CAF50;
  color: white;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 1rem;
}

.setup-btn:hover {
  background-color: #45a049;
}

.form-checkbox {
  margin-right: 0.5rem;
}

.batch-setup-link {
  margin-top: 1.5rem;
  text-align: center;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.link-button {
  color: #0066cc;
  text-decoration: none;
  font-weight: 500;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.link-button:hover {
  background-color: #e6f0ff;
  text-decoration: underline;
}

@media (max-width: 900px) {
  .two-step-container {
    flex-direction: column;
  }
  .sidebar {
    flex-direction: row;
    width: 100%;
    margin: 0;
    border-radius: 16px 16px 0 0;
    min-height: unset;
    justify-content: center;
    gap: 1rem;
    padding: 1rem 0;
  }
  .step-content {
    margin: 0;
    border-radius: 0 0 16px 16px;
    padding: 2rem 1rem;
  }
}
</style>