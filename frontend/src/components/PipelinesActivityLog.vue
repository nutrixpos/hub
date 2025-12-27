<template>
  <div class="activity-log-wrapper">
    <div class="!flex !flex-col !h-full !bg-gray-50 !font-sans !text-gray-900">
      <div class="!flex !flex-1 !overflow-hidden">
        <!-- Left Panel: Run List -->
        <div class="!w-1/3 !min-w-[320px] !border-r !border-gray-200 !flex !flex-col !bg-white">
          <div class="!p-4 !border-b !border-gray-200 !bg-gray-50/50">
            <!-- Search -->
            <div class="!relative !mb-3">
              <component :is="getIcon('Search')" class="!absolute !left-3 !top-2.5 !text-gray-400" :size="16" />
              <input 
                type="text" 
                placeholder="Search runs..."
                class="!w-full !pl-9 !pr-4 !py-2 !bg-white !border !border-gray-200 !rounded-lg !text-sm focus:!outline-none focus:!ring-2 focus:!ring-[#001F3E]/20 focus:!border-[#001F3E] !transition-all !m-0"
                v-model="searchTerm"
              />
            </div>
            <!-- Filter Tabs -->
            <div class="!flex !gap-2">
              <button
                v-for="status in ['all', 'completed', 'failed']"
                :key="status"
                @click="filterStatus = status"
                class="!px-3 !py-1.5 !rounded-md !text-xs !font-medium !capitalize !transition-colors !border"
                :class="[
                  filterStatus === status 
                  ? '!bg-white !border-gray-300 !text-gray-900 !shadow-sm' 
                  : '!bg-transparent !border-transparent !text-gray-500 hover:!bg-gray-200'
                ]"
              >
                {{ status }}
              </button>
            </div>
          </div>

          <!-- List -->
          <div class="!flex-1 !overflow-y-auto custom-scrollbar">
            <div v-if="filteredRuns.length === 0" class="!p-8 !text-center !text-gray-400">
              <div class="!mb-2 !flex !justify-center !opacity-50">
                <component :is="getIcon('Search')" :size="24"/>
              </div>
              <p class="!text-sm !m-0">No runs found.</p>
            </div>
            <div v-else>
              <div 
                v-for="run in filteredRuns"
                :key="run.id"
                @click="selectedRunId = run.id"
                class="!p-4 !border-b !border-gray-100 !cursor-pointer !transition-all hover:!bg-gray-50 !relative !group"
                :class="[
                  selectedRunId === run.id ? '!bg-blue-50/50 !border-l-4 !border-l-[#001F3E]' : '!border-l-4 !border-l-transparent'
                ]"
              >
                <div class="!flex !justify-between !items-start !mb-1">
                  <h4 
                    class="!text-sm !font-bold !truncate !pr-2 !m-0"
                    :class="selectedRunId === run.id ? '!text-[#001F3E]' : '!text-gray-700'"
                  >
                    {{ run.workflowName }}
                  </h4>
                  <span class="!text-xs !text-gray-400 !font-mono !shrink-0">{{ getDuration(run.startTime, run.endTime) }}</span>
                </div>
                
                <div class="!flex !items-center !gap-2 !mb-2">
                  <span class="!text-[10px] !font-bold !px-1.5 !py-0.5 !rounded !flex !items-center !gap-1 !border"
                    :class="{
                      '!bg-green-50 !text-green-700 !border-green-100': run.status === 'completed',
                      '!bg-red-50 !text-red-700 !border-red-100': run.status === 'failed',
                      '!bg-blue-50 !text-blue-700 !border-blue-100': run.status !== 'completed' && run.status !== 'failed'
                    }"
                  >
                    <component :is="getIcon('CheckCircle')" v-if="run.status === 'completed'" :size="10" stroke-width="3" />
                    <component :is="getIcon('XCircle')" v-if="run.status === 'failed'" :size="10" stroke-width="3" />
                    <component :is="getIcon('Loader2')" v-if="run.status === 'running'" :size="10" class="!animate-spin" />
                    <span class="!ml-1">{{ run.status.toUpperCase() }}</span>
                  </span>
                  <span class="!text-[10px] !text-gray-400">
                    {{ formatDate(run.startTime) }}
                  </span>
                </div>
                
                <div class="!mt-1.5 !flex !items-center">
                  <span class="!text-[10px] !font-mono !text-gray-400 !bg-gray-100 !px-1.5 !py-0.5 !rounded !border !border-gray-200">
                    #{{ run.id.replace(/^run-/, '') }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Panel: Logs -->
        <div class="!flex-1 !bg-[#1e1e1e] !flex !flex-col !overflow-hidden">
          <template v-if="selectedRun">
            <div class="!h-14 !border-b !border-gray-700 !flex !items-center !justify-between !px-6 !bg-[#252526] !shrink-0">
              <div class="!flex !items-center !gap-3">
                <div class="!text-gray-400">
                  <component :is="getIcon('Terminal')" :size="18" />
                </div>
                <div class="!text-sm !font-mono !text-gray-300">
                  Console Output
                </div>
              </div>
              <div class="!flex !gap-4 !text-xs !text-gray-400 !font-mono">
                <span>Run #{{ selectedRun.id.replace(/^run-/, '') }}</span>
              </div>
            </div>
            
            <div class="!flex-1 !overflow-y-auto !p-6 !font-mono !text-sm custom-scrollbar">
              <div 
                v-for="(log, idx) in selectedRun.logs" 
                :key="idx" 
                class="!flex !gap-4 !mb-2 hover:!bg-white/5 !p-1 !rounded !-mx-1 !transition-colors !leading-relaxed !group"
              >
                <span class="!text-gray-600 !shrink-0 !select-none !w-28 !text-xs !pt-0.5">
                  {{ formatLogTime(log.timestamp) }}
                </span>
                <span 
                  class="!shrink-0 !font-bold !w-16 !text-xs !pt-0.5"
                  :class="{
                    '!text-red-400': log.level === 'error',
                    '!text-yellow-400': log.level === 'warning',
                    '!text-blue-400': log.level !== 'error' && log.level !== 'warning'
                  }"
                >
                  {{ log.level.toUpperCase() }}
                </span>
                <span 
                  class="!break-all !whitespace-pre-wrap !flex-1"
                  :class="log.level === 'error' ? '!text-red-100' : '!text-gray-300'"
                >
                  {{ log.message }}
                </span>
              </div>

              <div v-if="selectedRun.status === 'running'" class="!flex !items-center !gap-2 !text-gray-500 !mt-4 !pl-32 !animate-pulse">
                <div class="!w-2 !h-2 !bg-green-500 !rounded-full"></div>
                <span class="!text-xs">Live stream active...</span>
              </div>
              <div class="!h-12"></div>
            </div>
          </template>
          
          <div v-else class="!flex-1 !flex !flex-col !items-center !justify-center !text-gray-600">
            <component :is="getIcon('Terminal')" :size="48" class="!opacity-20 !mb-4"/>
            <p class="!m-0 !mt-4">Select a run to view logs</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, getCurrentInstance } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import axios from 'axios';


const { proxy } = getCurrentInstance();

// ==========================================
// STATE VARIABLES (Replaces Props)
// ==========================================

const workflows = ref([
  {
    id: '1',
    name: 'Restock to n8n',
  }
]);

const allRuns = ref({
  '1': [
      {
          id: 'run-demo-1',
          startTime: new Date(Date.now() - 3600000).toISOString(),
          endTime: new Date(Date.now() - 3595000).toISOString(),
          status: 'completed',
          logs: [
              { timestamp: new Date(Date.now() - 3600000).toISOString(), level: 'info', message: 'Workflow execution started.' },
              { timestamp: new Date(Date.now() - 3599000).toISOString(), level: 'info', message: 'Checking stock levels for configured items.' },
              { timestamp: new Date(Date.now() - 3598000).toISOString(), level: 'warning', message: 'Item WIDGET-X is low on stock (Count: 4).' },
              { timestamp: new Date(Date.now() - 3597000).toISOString(), level: 'info', message: 'Executing action: n8n Webhook.' },
              { timestamp: new Date(Date.now() - 3595000).toISOString(), level: 'info', message: 'Workflow completed successfully.' }
          ]
      }
  ],
});

// ==========================================
// LOGIC
// ==========================================

const getWorkflows = () => {
   axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/workflows`, {
       headers: {
           Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
       },
   })
   .then((response)=>{
    // workflows.value = []

    for (var i=0;i<response.data.data.length;i++){
      workflows.value.push({
        id: response.data.data[i].id,
        name: response.data.data[i].name,
      })

      if (!allRuns.value[response.data.data[i].id]){
        allRuns.value[response.data.data[i].id] = []
      }

      allRuns.value[response.data.data[i].id] = response.data.data[i].runs


      for (var j=0;j<response.data.data[i].runs.length;j++){
        allRuns.value[response.data.data[i].id].push({
          id: response.data.data[i].runs[j],
          startTime: response.data.data[i].runs[j].start_time,
          endTime: response.data.data[i].runs[j].start_time,
          status: response.data.data[i].runs[j].status,
          logs: response.data.data[i].runs[j].logs,
        })
      }
    }

   })
   .catch((err) => {
       console.log(err)
   });
}

getWorkflows();


const searchTerm = ref('');
const filterStatus = ref('all');
const selectedRunId = ref(null);

const getIcon = (name) => {
  return LucideIcons[name] || LucideIcons.HelpCircle;
};

// Flatten and sort runs
const flatRuns = computed(() => {
  const list = [];
  Object.entries(allRuns.value).forEach(([wfId, runs]) => {
    const wf = workflows.value.find(w => w.id === wfId);
    if (runs && Array.isArray(runs)) {
        runs.forEach(r => {
        list.push({
            ...r,
            workflowName: wf ? wf.name : 'Unknown Workflow',
            workflowId: wfId
        });
        });
    }
  });
  return list.sort((a, b) => new Date(b.startTime).getTime() - new Date(a.startTime).getTime());
});

// Filter runs based on search and status
const filteredRuns = computed(() => {
  return flatRuns.value.filter(run => {
    const matchesSearch = 
      run.id.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      run.workflowName.toLowerCase().includes(searchTerm.value.toLowerCase());
    const matchesStatus = filterStatus.value === 'all' || run.status === filterStatus.value;
    return matchesSearch && matchesStatus;
  });
});

// Get the currently selected run object
const selectedRun = computed(() => {
  let run = flatRuns.value.find(r => r.id === selectedRunId.value);
  if (!run && filteredRuns.value.length > 0) {
    return filteredRuns.value[0];
  }
  return run || null;
});

// Auto-select first run if none selected
watch(filteredRuns, (newRuns) => {
  if (newRuns.length > 0 && !selectedRunId.value) {
    selectedRunId.value = newRuns[0].id;
  }
}, { immediate: true });

const formatDate = (iso) => {
  if (!iso) return '';
  return new Date(iso).toLocaleString(undefined, {
    month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit'
  });
};

const formatLogTime = (iso) => {
  if (!iso) return '';
  return new Date(iso).toLocaleTimeString([], { 
    hour12: false, 
    hour: '2-digit', 
    minute:'2-digit', 
    second:'2-digit', 
    fractionalSecondDigits: 3 
  });
};

const getDuration = (start, end) => {
  if (!end) return '...';
  const diff = new Date(end).getTime() - new Date(start).getTime();
  return `${(diff / 1000).toFixed(2)}s`;
};
</script>

<style scoped>
.activity-log-wrapper {
  all: unset !important;
  display: block !important;
  width: 100% !important;
  height: 100% !important;
  box-sizing: border-box !important;
  font-family: 'Inter', sans-serif !important;
}

.activity-log-wrapper * {
  box-sizing: border-box !important;
}

/* Custom scrollbar copied from global styling for local containment */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px !important;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent !important;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e1 !important;
  border-radius: 20px !important;
}
</style>