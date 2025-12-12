<template>
<div class="workflow-card-wrapper">
    <div 
    class="!group !bg-white !rounded-xl !border !shadow-sm hover:!shadow-xl hover:!-translate-y-1 !transition-all !duration-200 !cursor-pointer !flex !flex-col !relative !overflow-hidden"
    :class="[
        workflow.enabled == true ? '!border-gray-200 hover:!border-[#001F3E]' : '!border-gray-200 !opacity-80 hover:!opacity-100'
    ]"
    @click="navigate"
    >
        <!-- Running Progress Bar -->
        <div v-if="isRunning" class="!absolute !top-0 !left-0 !w-full !h-1.5 !bg-gray-100 !z-[60] !overflow-hidden">
            <div class="!h-full !bg-[#001F3E] !z-50 animate-indeterminate"></div>
        </div>

        <div class="!p-6 !flex-1 !flex !flex-col">
            <div class="!flex !justify-between !items-start !mb-4">
                <!-- Trigger Icon -->
                <div class="!p-2.5 !rounded-lg !shadow-sm" :class="workflow.trigger ? '!bg-pink-50 !text-pink-600' : '!bg-gray-50 !text-gray-400'">
                    <component :is="getIcon(workflow.trigger ? 'Zap' : 'HelpCircle')" :size="20" />
                </div>
                
                <div class="!flex !items-center !gap-3">
                    <!-- Manual Run Button -->
                    <button 
                        v-if="workflow.trigger?.type === 'trigger-manual'"
                        @click.stop="toggleRun"
                        :disabled="isRunning"
                        class="!flex !items-center !justify-center !w-8 !h-8 !rounded-full !border !transition-all !shadow-sm !mr-1"
                        :class="[
                            isRunning 
                            ? '!bg-gray-50 !text-gray-400 !border-gray-200 !cursor-not-allowed' 
                            : '!bg-white !text-[#001F3E] !border-[#001F3E] hover:!bg-[#001F3E] hover:!text-white'
                        ]"
                        title="Run manually"
                    >
                        <component :is="getIcon('Loader2')" v-if="isRunning" :size="14" class="!animate-spin" />
                        <component :is="getIcon('Play')" v-else :size="14" class="!fill-current !ml-0.5" />
                    </button>

                    <!-- Status Badge -->
                    <div class="!px-2.5 !py-0.5 !rounded-full !text-xs !font-bold !flex !items-center !gap-1.5"
                        :class="workflow.enabled == true ? '!bg-green-100 !text-green-800' : '!bg-gray-100 !text-gray-800'"
                    >
                        <span v-if="workflow.enabled == true" class="!relative !flex !h-2.5 !w-2.5">
                            <span class="!animate-ping !absolute !inline-flex !h-full !w-full !rounded-full !bg-green-50 !opacity-75"></span>
                            <span class="!relative !inline-flex !rounded-full !h-2.5 !w-2.5 !bg-green-600"></span>
                        </span>
                        {{ workflow.enabled == true ? 'ENABLED' : 'DISABLED' }}
                    </div>
                    
                    <!-- Toggle Switch -->
                    <button
                        @click.stop="toggleEnable"
                        class="!relative !inline-flex !h-6 !w-11 !flex-shrink-0 !cursor-pointer !rounded-full !border-2 !border-transparent !transition-colors !duration-200 !ease-in-out focus:!outline-none"
                        :class="workflow.enabled == true ? '!bg-[#001F3E]' : '!bg-gray-200'"
                        role="switch"
                        :aria-checked="workflow.enabled == true"
                        :title="workflow.enabled == true ? 'Disable Workflow' : 'Enable Workflow'"
                    >
                        <span
                            aria-hidden="true"
                            class="!pointer-events-none !inline-block !h-5 !w-5 !transform !rounded-full !bg-white !shadow !ring-0 !transition !duration-200 !ease-in-out"
                            :class="workflow.enabled == true ? '!translate-x-5' : '!translate-x-0'"
                        />
                    </button>
                </div>
            </div>

            <!-- Name & Desc -->
            <h3 class="!text-lg !font-bold !text-gray-900 !mb-2 group-hover:!text-[#001F3E] !transition-colors">
                {{ workflow.name }}
            </h3>
            <p class="!text-sm !text-gray-500 !line-clamp-2 !leading-relaxed !mb-4">
                {{ workflow.description || 'No description provided.' }}
            </p>
            
            <!-- Workflow Steps -->
            <div class="!mt-auto !pt-4 !border-t !border-gray-100 !flex !flex-col !gap-3">
                <div class="!flex !items-center !gap-3">
                    <div class="!w-6 !h-6 !rounded !bg-gray-100 !flex !items-center !justify-center !text-gray-500 !shrink-0">
                        <component :is="getIcon('Zap')" :size="14" />
                    </div>
                    <span class="!text-sm !font-medium !text-gray-700 !truncate">
                        {{ workflow.trigger?.type }}
                    </span>
                </div>
                
                <div v-if="workflow.actions.length > 0" class="!ml-2.5 !w-0.5 !h-2 !bg-gray-200"></div>

                <div v-for="action in workflow.actions" :key="action.id" class="!flex !items-center !gap-3">
                    <div class="!w-6 !h-6 !rounded !bg-[#001F3E]/10 !flex !items-center !justify-center !text-[#001F3E] !shrink-0">
                        <component :is="getIcon('Activity')" :size="14" />
                    </div>
                    <span class="!text-sm !text-gray-600 !truncate">
                        {{ action.type }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Footer / Last Run -->
        <div class="!px-6 !py-4 !border-t !border-gray-100 !bg-gray-50/50 !rounded-b-xl !flex !justify-between !items-center">
            <div class="!flex !items-center !gap-2 !min-w-0">
                <template v-if="lastRun">
                    <div class="!w-2 !h-2 !rounded-full !shrink-0"
                        :class="{
                            '!bg-green-500': lastRun.status === 'completed',
                            '!bg-red-500': lastRun.status === 'failed',
                            '!bg-blue-500 !animate-pulse': lastRun.status !== 'completed' && lastRun.status !== 'failed'
                        }"
                    ></div>
                    <div class="!flex !flex-col">
                        <span class="!text-xs !font-semibold !text-gray-700 !capitalize !leading-none">
                            {{ lastRun.status }}
                        </span>
                        <div class="!flex !items-center !gap-1 !text-[10px] !text-gray-400 !leading-none !mt-1">
                            <span>{{ formatDate(lastRun.endTime || lastRun.startTime) }}</span>
                            <span class="!text-gray-300">•</span>
                            <span class="!font-mono !text-gray-500">#{{ lastRun.id.replace(/^run-/, '') }}</span>
                        </div>
                    </div>
                </template>
                <span v-else class="!text-xs !font-medium !text-gray-400">No runs yet</span>
            </div>
            <router-link 
                :to="`/console/workflows/put/${workflow.id}`"
                class="!text-gray-400 hover:!text-gray-600 !transition-colors !p-1.5 hover:!bg-gray-50 !rounded-md"
            >
                <component :is="getIcon('Pencil')" :size="16" />
        </router-link>
        </div>
    </div>
</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted,defineProps } from 'vue';
import * as LucideIcons from 'lucide-vue-next';

const emit = defineEmits(['edit', 'run', 'toggle', 'delete']);



// ==========================================
// MOCK DATA (Replaces Props)
// ==========================================

const NODE_DEFINITIONS = [
  {
    id: 'trigger_low_stock',
    label: 'Low Stock Alert'
  },
  {
    id: 'trigger-manual',
    label: 'Manual Trigger'
  },
  {
    id: 'action-n8n-webhook',
    label: 'n8n Webhook'
  }
];

const props = defineProps({
  workflow: {
    type: Object,
    required: true
  }
});

const isRunning = ref(false);

const lastRun = ref({
  id: 'run-12345',
  status: 'completed',
  startTime: new Date(Date.now() - 3600000).toISOString(),
  endTime: new Date(Date.now() - 3595000).toISOString()
});

let autoRunTimer = null;

// ==========================================
// LOCAL LOGIC
// ==========================================

const toggleEnable = () => {
    emit('edit', {
        ...props.workflow,
        enabled: !props.workflow.enabled
    });
};

const toggleRun = () => {
    if (isRunning.value) return;
    isRunning.value = true;
    
    // Duration: Minimum 10 seconds (10000ms) plus up to 3s random variance
    const duration = 10000 + Math.random() * 3000;
    
    setTimeout(() => {
        isRunning.value = false;
        lastRun.value = {
            id: `run-${Math.floor(Math.random() * 10000)}`,
            status: Math.random() > 0.3 ? 'completed' : 'failed',
            startTime: new Date(Date.now() - duration).toISOString(),
            endTime: new Date().toISOString()
        };
    }, duration);
};

const startSimulationLoop = () => {
    // Schedule next run between 15s and 30s from now (Intermittent)
    const nextInterval = 15000 + (Math.random() * 15000); 
    
    autoRunTimer = setTimeout(() => {
        if (workflow.value.status === 'active') {
            toggleRun();
        }
        startSimulationLoop();
    }, nextInterval);
};

// onMounted(() => {
//     // Trigger first random run attempt after mount
//     setTimeout(() => {
//         if (props.workflow.value.status === 'active' && Math.random() > 0.5) {
//              toggleRun();
//         }
//         startSimulationLoop();
//     }, 2000);
// });

// onUnmounted(() => {
//     if (autoRunTimer) clearTimeout(autoRunTimer);
// });

const getIcon = (name) => {
  return LucideIcons[name] || LucideIcons.HelpCircle;
};

const formatDate = (dateStr) => {
    if (!dateStr) return '';
    return new Date(dateStr).toLocaleString(undefined, { 
        month: 'short', 
        day: 'numeric', 
        hour: '2-digit', 
        minute: '2-digit' 
    });
};
</script>

<style scoped>
.workflow-card-wrapper {
  all: unset !important;
  display: block !important;
  width: 100% !important;
  box-sizing: border-box !important;
  font-family: 'Inter', sans-serif !important;
}

.workflow-card-wrapper * {
  box-sizing: border-box !important;
}

@keyframes indeterminate {
  0% { left: -35%; right: 100%; }
  60% { left: 100%; right: -90%; }
  100% { left: 100%; right: -90%; }
}

.animate-indeterminate {
  position: absolute !important;
  top: 0 !important;
  bottom: 0 !important;
  background-color: #001F3E !important;
  animation: indeterminate 2s cubic-bezier(0.65, 0.815, 0.735, 0.395) infinite !important;
}
</style>