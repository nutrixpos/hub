<template>
  <div class="workflow-editor-wrapper">
    <div class="!flex !h-full !bg-gray-50 !relative !overflow-hidden !font-sans !text-gray-900">
      
      <!-- Main Canvas Area -->
      <div 
        class="!flex-1 !flex !flex-col !h-full !transition-all !duration-300"
        :class="{ '!blur-sm !pointer-events-none !select-none': selectedNode }"
      >
        <!-- Editor Header (Name/Desc/Save) - No Navigation -->
        <div class="!h-16 !bg-white !border-b !border-gray-200 !flex !items-center !justify-between !px-6 !shrink-0 !z-10 !shadow-sm">
          <div class="!flex !items-center !gap-4">
            <div>
              <input 
                class="!text-lg !font-bold !text-gray-900 !border-none focus:!ring-0 !p-0 !bg-transparent !placeholder-gray-400 !w-full !outline-none"
                placeholder="Workflow Name"
                v-model="workflow.name"
              />
              <input 
                class="!block !text-sm !text-gray-500 !border-none focus:!ring-0 !p-0 !bg-transparent !w-64 !placeholder-gray-300 !outline-none"
                placeholder="Description (optional)"
                v-model="workflow.description"
              />
            </div>
          </div>
          <div class="!flex !items-center !gap-3">
            <Button label="Save workflow" icon="pi pi-save" @click="handleSave" :loading="loading_saving_workflow" />
          </div>
        </div>

        <!-- Canvas -->
        <div class="!flex-1 !overflow-y-auto !p-8 custom-scrollbar">
          <div class="!max-w-3xl !mx-auto !flex !flex-col !items-center !pb-20">
            
            <div class="!mb-4 !text-xs !font-bold !text-gray-400 !uppercase !tracking-widest">Start</div>

            <!-- Trigger Node -->
            <div v-if="workflow.trigger" class="!flex !flex-col !items-center !w-full">
               <div 
                @click="handleNodeClick({...workflow.trigger})"
                class="!relative !w-80 !p-4 !rounded-xl !border-2 !transition-all !cursor-pointer !shadow-sm !group !bg-white"
                :class="selectedNodeId === workflow.trigger.id ? '!border-gray-900 !ring-2 !ring-gray-900/20' : '!border-gray-200 hover:!border-[#001F3E]'"
              >
                <div class="!flex !items-start !gap-3">
                  <div class="!p-2.5 !rounded-lg !shrink-0 !text-white !shadow-sm" :class="getDefinition(workflow.trigger.definitionId)?.color">
                    <component :is="getIcon(getDefinition(workflow.trigger.definitionId)?.icon)" :size="20" />
                  </div>
                  <div class="!flex-1 !min-w-0">
                    <div class="!flex !items-center !justify-between">
                      <span class="!text-xs !font-bold !text-gray-400 !uppercase !tracking-wider !mb-0.5 !block">
                          Trigger
                      </span>
                      <div class="!w-2 !h-2 !rounded-full !bg-green-400"></div>
                    </div>
                    <h4 class="!font-semibold !text-gray-900 !truncate">{{ definitionId }}</h4>
                    <p class="!text-xs !text-gray-500 !truncate !mt-1">
                      {{ getPropertiesString(workflow.trigger.properties, workflow.trigger.definitionId) }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
            
            <button 
              v-else
              @click="addingType = 'TRIGGER'"
              class="!w-64 !h-24 !border-2 !border-dashed !border-gray-300 !rounded-xl !flex !flex-col !items-center !justify-center !text-gray-400 hover:!border-[#001F3E] hover:!text-[#001F3E] hover:!bg-gray-50 !transition-all !group"
            >
              <div class="!bg-gray-100 !p-2 !rounded-full group-hover:!bg-[#001F3E] group-hover:!text-white !mb-2 !transition-colors !text-gray-500">
                <component :is="getIcon('Zap')" :size="20" />
              </div>
              <span class="!text-sm !font-medium">Add Trigger</span>
            </button>

            <div class="!h-8 !w-0.5 !bg-gray-300 !my-1 !relative">
                 <component :is="getIcon('ArrowDown')" class="!absolute !-bottom-2 !-left-2.5 !text-gray-300" :size="20" />
            </div>

            <!-- Action Nodes -->
            <template v-for="(action, index) in workflow.actions" :key="action.id">
                <div 
                  @click="handleNodeClick(action)"
                  class="!relative !w-80 !p-4 !rounded-xl !border-2 !transition-all !cursor-pointer !shadow-sm !group !bg-white"
                  :class="selectedNodeId === action.id ? '!border-gray-900 !ring-2 !ring-gray-900/20' : '!border-gray-200 hover:!border-[#001F3E]'"
                >
                  <div class="!flex !items-start !gap-3">
                    <div class="!p-2.5 !rounded-lg !shrink-0 !text-white !shadow-sm" :class="getDefinition(action.definitionId)?.color">
                      <component :is="getIcon(getDefinition(action.definitionId)?.icon)" :size="20" />
                    </div>
                    <div class="!flex-1 !min-w-0">
                      <div class="!flex !items-center !justify-between">
                        <span class="!text-xs !font-bold !text-gray-400 !uppercase !tracking-wider !mb-0.5 !block">
                            Action
                        </span>
                        <div class="!w-2 !h-2 !rounded-full !bg-green-400"></div>
                      </div>
                      <h4 class="!font-semibold !text-gray-900 !truncate">{{ getDefinition(action.definitionId)?.label }}</h4>
                      <p class="!text-xs !text-gray-500 !truncate !mt-1">
                        {{ getPropertiesString(action.properties, action.definitionId) }}
                      </p>
                    </div>
                  </div>
                </div>

                <div v-if="index < workflow.actions.length - 1 || workflow.actions.length < 1" class="!h-8 !w-0.5 !bg-gray-300 !my-1 !relative">
                    <component :is="getIcon('ArrowDown')" class="!absolute !-bottom-2 !-left-2.5 !text-gray-300" :size="20" />
                </div>
            </template>

            <!-- Add Action Button -->
            <button 
              v-if="workflow.actions.length < 1"
              @click="addingType = 'ACTION'"
              class="!mt-2 !flex !items-center !gap-2 !bg-white !border !border-gray-300 !px-4 !py-2 !rounded-full !text-sm !font-medium !text-gray-600 !shadow-sm hover:!border-[#001F3E] hover:!text-[#001F3E] !transition-all"
            >
              <component :is="getIcon('Plus')" :size="16" />
              Add Action
            </button>

          </div>
        </div>
      </div>

      <!-- Right Panel: Configuration -->
      <div v-if="selectedNode && selectedDef" class="!w-96 !bg-white !border-l !border-gray-200 !h-full !flex !flex-col !shadow-xl !z-20 !transition-transform">
          <div class="!h-16 !border-b !border-gray-200 !flex !items-center !justify-between !px-6 !shrink-0 !z-30 !bg-white">
            <div class="!flex !items-center !gap-3 !flex-1 !mr-4">
                <div class="!p-2 !rounded-lg !text-white !shadow-sm !transition-colors" :class="selectedDef.color">
                   <component :is="getIcon(selectedDef.icon)" :size="18" />
                </div>
                
                <div class="!relative !flex-1" ref="typeSelectorRef">
                   <button
                      @click="isTypeSelectorOpen = !isTypeSelectorOpen"
                      class="!flex !items-center !gap-3 !font-bold !text-gray-900 hover:!bg-gray-100 !px-3 !py-2 !-ml-2 !rounded-xl !transition-colors !group !w-full !text-left"
                   >
                      <component :is="getIcon('ChevronDown')" stroke-width="3" :size="20" class="!text-gray-900 !transition-transform !duration-200 !shrink-0" :class="isTypeSelectorOpen ? '!rotate-180' : ''" />
                      
                      <div class="!flex !flex-col !min-w-0">
                          <span class="!truncate !text-lg !leading-tight">{{ selectedDef.label }}</span>
                          <div class="!text-[10px] !font-medium !text-gray-400 !uppercase !tracking-wide">
                             Switch {{ selectedDef.type === 'TRIGGER' ? 'Trigger' : 'Action' }}
                          </div>
                      </div>
                   </button>

                   <!-- Type Dropdown -->
                   <div v-if="isTypeSelectorOpen" class="!absolute !top-full !left-0 !mt-2 !w-64 !bg-white !rounded-xl !shadow-xl !border !border-gray-100 !z-50 !overflow-hidden !ring-1 !ring-black/5">
                       <div class="!p-1 !max-h-64 !overflow-y-auto custom-scrollbar">
                           <div class="!px-3 !py-2 !text-xs !font-semibold !text-gray-400 !uppercase !tracking-wider !bg-gray-50/50 !mb-1 !rounded-t-lg">
                              Switch {{ selectedDef.type === 'TRIGGER' ? 'Trigger' : 'Action' }}
                           </div>
                           <button 
                                v-for="def in NODE_DEFINITIONS.filter(d => d.type === selectedDef.type)"
                                :key="def.id"
                               @click="handleDefinitionChange(def.id); isTypeSelectorOpen = false"
                               class="!w-full !flex !items-center !gap-3 !p-2 !rounded-lg !text-sm !transition-all !text-left !mb-0.5"
                               :class="def.id === selectedDef.id ? '!bg-[#001F3E]/5 !text-[#001F3E] !font-medium' : '!text-gray-700 hover:!bg-gray-50'"
                           >
                                <div class="!p-1.5 !rounded-md !shrink-0 !text-white !scale-90 !shadow-sm" :class="def.color">
                                    <component :is="getIcon(def.icon)" :size="14" />
                                </div>
                                <span class="!flex-1 !truncate">{{ def.label }}</span>
                                <component :is="getIcon('Check')" v-if="def.id === selectedDef.id" :size="14" class="!text-[#001F3E] !shrink-0" />
                           </button>
                       </div>
                   </div>
                </div>
            </div>
            <button @click="selectedNodeId = null" class="!text-gray-400 hover:!text-gray-600 !shrink-0">
              <component :is="getIcon('X')" :size="20" />
            </button>
          </div>
          
          <div class="!flex-1 !overflow-y-auto !p-6 custom-scrollbar">
            <p class="!text-sm !text-gray-500 !mb-6 !border-l-2 !border-gray-200 !pl-3">{{ selectedDef.description }}</p>
            
            <!-- Dynamic Form -->
            <div v-if="selectedDef.schema.length > 0" class="!space-y-4">
              <template v-for="field in selectedDef.schema" :key="field.key">
                 <div v-if="checkCondition(field, selectedNode.properties)" class="!flex !flex-col !space-y-1">
                    <label class="!text-sm !font-medium !text-gray-700">
                      {{ field.label }} <span v-if="field.required" class="!text-red-500 !ml-1">*</span>
                    </label>

                    <!-- Text -->
                    <input 
                      v-if="field.type === 'TEXT'"
                      type="text" 
                      :value="selectedNode.properties[field.key]"
                      @input="updateNodeProperty(field.key, $event.target.value)"
                      :placeholder="field.placeholder"
                      class="!border !border-gray-300 !rounded-md !px-3 !py-2 !text-sm focus:!outline-none focus:!ring-2 focus:!ring-[#001F3E] focus:!border-transparent !transition-all !w-full"
                    />

                    <!-- Number -->
                    <input 
                      v-if="field.type === 'NUMBER'"
                      type="number" 
                      :value="selectedNode.properties[field.key]"
                      @input="updateNodeProperty(field.key, parseFloat($event.target.value))"
                      :placeholder="field.placeholder"
                      class="!border !border-gray-300 !rounded-md !px-3 !py-2 !text-sm focus:!outline-none focus:!ring-2 focus:!ring-[#001F3E] focus:!border-transparent !transition-all !w-full"
                    />

                    <!-- Select -->
                    <select
                      v-if="field.type === 'SELECT'"
                      :value="selectedNode.properties[field.key] || ''"
                      @change="updateNodeProperty(field.key, $event.target.value)"
                      class="!border !border-gray-300 !rounded-md !px-3 !py-2 !text-sm !bg-white focus:!outline-none focus:!ring-2 focus:!ring-[#001F3E] focus:!border-transparent !transition-all !w-full"
                    >
                      <option value="" disabled>Select an option</option>
                      <option 
                        v-for="opt in field.options" 
                        :key="typeof opt === 'object' ? opt.id : opt" 
                        :value="typeof opt === 'object' ? opt.id : opt"
                      >
                        {{ typeof opt === 'object' ? opt.title : opt }}
                      </option>
                    </select>

                    <!-- Multi Select (Enhanced) -->
                    <div v-if="field.type === 'MULTI_SELECT'" class="!relative" :ref="el => setMultiSelectRef(el, field.key)">
                        <!-- Trigger -->
                        <div 
                            @click="toggleMultiSelect(field.key)"
                            class="!border !border-gray-300 !rounded-md !px-3 !py-2 !text-sm !bg-white !cursor-pointer !flex !justify-between !items-center hover:!border-[#001F3E] !transition-colors !min-h-[38px]"
                        >
                             <span class="!truncate !pr-2" :class="(selectedNode.properties[field.key] || []).length ? '!text-gray-900' : '!text-gray-400'">
                                {{ (selectedNode.properties[field.key] || []).length > 0 ? `${(selectedNode.properties[field.key] || []).length} selected` : (field.placeholder || 'Select items...') }}
                             </span>
                             <component :is="getIcon('ChevronDown')" :size="16" class="!text-gray-400 !shrink-0" />
                        </div>

                        <!-- Dropdown -->
                        <div v-if="multiSelectOpen === field.key" class="!absolute !z-50 !mt-1 !w-full !bg-white !border !border-gray-300 !rounded-md !shadow-xl !max-h-64 !flex !flex-col !overflow-hidden">
                             <!-- Search -->
                             <div class="!p-2 !border-b !border-gray-100 !bg-gray-50">
                                 <div class="!relative">
                                     <component :is="getIcon('Search')" class="!absolute !left-2.5 !top-2 !text-gray-400" :size="14" />
                                     <input 
                                        type="text"
                                        v-model="multiSelectSearch"
                                        class="!w-full !pl-8 !pr-3 !py-1.5 !text-sm !border !border-gray-200 !rounded-md focus:!outline-none focus:!border-[#001F3E] focus:!ring-1 focus:!ring-[#001F3E] !bg-white"
                                        placeholder="Search..."
                                        @click.stop
                                        ref="multiSelectSearchInput"
                                     />
                                 </div>
                             </div>
                             <!-- Options -->
                             <div class="!overflow-y-auto !flex-1">
                                 <div 
                                    v-for="opt in filterOptions(field.options)" 
                                    :key="opt"
                                    @click="toggleMultiSelectOption(field.key, opt)"
                                    class="!flex !items-center !justify-between !px-3 !py-2.5 !cursor-pointer !transition-colors hover:!bg-gray-50"
                                    :class="isSelected(field.key, opt) ? '!bg-blue-50' : ''"
                                 >
                                     <div class="!flex !items-center !gap-3">
                                        <div class="!w-4 !h-4 !rounded !border !flex !items-center !justify-center !transition-colors"
                                             :class="isSelected(field.key, opt) ? '!bg-[#001F3E] !border-[#001F3E]' : '!border-gray-300 !bg-white'">
                                             <component :is="getIcon('Check')" v-if="isSelected(field.key, opt)" :size="10" class="!text-white" stroke-width="3" />
                                        </div>
                                        <span class="!text-sm" :class="isSelected(field.key, opt) ? '!text-gray-900 !font-medium' : '!text-gray-700'">{{ opt }}</span>
                                     </div>
                                 </div>
                                 <div v-if="filterOptions(field.options).length === 0" class="!p-4 !text-center !text-gray-400 !text-xs">No items found.</div>
                             </div>
                        </div>

                        <!-- Tags -->
                        <div v-if="(selectedNode.properties[field.key] || []).length > 0" class="!mt-2 !flex !flex-wrap !gap-1.5">
                             <span v-for="val in selectedNode.properties[field.key]" :key="val" class="!inline-flex !items-center !px-2 !py-0.5 !rounded !text-xs !font-medium !bg-gray-100 !text-gray-800 !border !border-gray-200">
                                 {{ val }}
                                 <button @click="toggleMultiSelectOption(field.key, val)" class="!ml-1.5 !text-gray-400 hover:!text-red-500">&times;</button>
                             </span>
                        </div>
                    </div>

                    <!-- Boolean -->
                    <div v-if="field.type === 'BOOLEAN'" class="!flex !items-center">
                        <input 
                            type="checkbox"
                            :checked="!!selectedNode.properties[field.key]"
                            @change="updateNodeProperty(field.key, $event.target.checked)"
                            class="!h-4 !w-4 !text-[#001F3E] focus:!ring-[#001F3E] !border-gray-300 !rounded"
                        />
                        <span class="!ml-2 !text-sm !text-gray-600">{{ field.description || 'Enable' }}</span>
                    </div>

                    <!-- Key Value List (Headers) -->
                    <div v-if="field.type === 'KEY_VALUE_LIST'" class="!flex !flex-col !gap-2">
                      <div v-for="(item, index) in (selectedNode.properties[field.key] || [])" :key="index" class="!flex !gap-2 !items-center grid">
                          <div class="col-5">
                            <input 
                                type="text" 
                                v-model="item.key" 
                                @input="updateNodeProperty(field.key, selectedNode.properties[field.key])"
                                placeholder="Key"
                                class="!flex-1 !border !border-gray-300 !rounded-md !px-3 !py-2 !text-sm focus:!outline-none focus:!ring-2 focus:!ring-[#001F3E] focus:!border-transparent w-full"
                            />
                          </div>
                          <div class="col-5">
                            <input 
                                type="text" 
                                v-model="item.value" 
                                @input="updateNodeProperty(field.key, selectedNode.properties[field.key])"
                                placeholder="Value"
                                class="!flex-1 !border !border-gray-300 !rounded-md !px-3 !py-2 !text-sm focus:!outline-none focus:!ring-2 focus:!ring-[#001F3E] focus:!border-transparent w-full"
                            />
                          </div>
                          <div class="col-1">
                              <button @click="removeKeyValue(field.key, index)" class="!text-gray-400 hover:!text-red-500 !p-1">
                                  <component :is="getIcon('Trash2')" :size="16" />
                              </button>
                          </div>
                      </div>
                      <button 
                          @click="addKeyValue(field.key)"
                          class="!mt-1 !flex !items-center !gap-2 !text-sm !font-medium !text-[#001F3E] hover:!text-blue-700 !transition-colors"
                      >
                          <component :is="getIcon('Plus')" :size="14" />
                          Add Header
                      </button>
                    </div>

                    <p v-if="field.description && field.type !== 'BOOLEAN'" class="!text-xs !text-gray-500">{{ field.description }}</p>

                 </div>
              </template>
            </div>
            <div v-else class="!text-gray-400 !text-sm !italic">No properties to configure.</div>

          </div>
          
          <div class="!p-4 !border-t !border-gray-200 !bg-gray-50">
             <button 
               @click="handleDeleteNode(selectedNode.id)"
               class="!w-full !py-2 !text-red-600 !border !border-red-200 !rounded-lg hover:!bg-red-50 !text-sm !font-medium !flex !items-center !justify-center !gap-2"
             >
                <component :is="getIcon('Trash2')" :size="16" />
                Delete Node
             </button>
          </div>
      </div>

      <!-- Adding Node Modal -->
      <div v-if="addingType" class="!absolute !inset-0 !bg-black/50 !backdrop-blur-sm !z-50 !flex !items-center !justify-center !p-4">
        <div class="!bg-white !rounded-xl !shadow-2xl !w-full !max-w-2xl !max-h-[80vh] !flex !flex-col !overflow-hidden">
          <div class="!p-5 !border-b !border-gray-100 !flex !justify-between !items-center">
            <h3 class="!text-lg !font-bold !text-[#001F3E]">
              Select {{ addingType === 'TRIGGER' ? 'Trigger' : 'Action' }}
            </h3>
            <button @click="addingType = null">
                <component :is="getIcon('X')" class="!text-gray-400 hover:!text-gray-600" />
            </button>
          </div>
          <div class="!p-5 !overflow-y-auto !grid !grid-cols-1 sm:!grid-cols-2 !gap-4">
            <button 
                v-for="def in NODE_DEFINITIONS.filter(d => d.type === addingType)"
                :key="def.id"
                @click="addingType === 'TRIGGER' ? handleSetTrigger(def.id) : handleAddAction(def.id)"
                class="!flex !flex-col !items-start !p-4 !border !border-gray-200 !rounded-xl hover:!border-[#001F3E] hover:!ring-1 hover:!ring-[#001F3E] !transition-all !text-left !bg-gray-50/50 hover:!bg-white"
              >
                <div class="!p-2 !rounded-lg !text-white !mb-3 !shadow-sm" :class="def.color">
                  <component :is="getIcon(def.icon)" :size="20" />
                </div>
                <div class="!font-semibold !text-gray-900 !mb-1">{{ def.label }}</div>
                <div class="!text-xs !text-gray-500 !leading-relaxed">{{ def.description }}</div>
            </button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted,getCurrentInstance } from 'vue';
import {Button} from 'primevue';
import * as LucideIcons from 'lucide-vue-next';
import axios from 'axios';
import { useRoute, useRouter } from 'vue-router'
import { useToast } from "primevue/usetoast";


const toast = useToast()
const {proxy} = getCurrentInstance()
const route = useRoute()
const router = useRouter()

const loading_saving_workflow = ref(false);

const definitionId = computed(() => {
  if (route.params.id != "") {
    return getDefinition(workflow.value.trigger.type)?.label 
  }
   return getDefinition(workflow.value.trigger.definitionId)?.label
});


// ==========================================
// CONSTANTS (Duplicated for self-containment)
// ==========================================
const NODE_DEFINITIONS = ref([
  {
    id: 'trigger_low_stock',
    type: 'TRIGGER',
    label: 'Low Stock Alert',
    description: 'Triggers when product inventory falls below a specific level.',
    icon: 'Package',
    color: '!bg-orange-500',
    schema: [
      { 
        key: 'monitor_type', 
        label: 'Monitor Type', 
        type: 'SELECT', 
        options: [
          { id: 'any_item', title: 'Any Item' },
          { id: 'specific_items', title: 'Specific Items' }
        ], 
        required: true 
      },
      { 
        key: 'product_ids', 
        label: 'Select Products', 
        type: 'MULTI_SELECT', 
        options: ['WIDGET-X', 'GADGET-Y', 'GIZMO-Z', 'SUPER-TOOL-2000', 'IPHONE-15'], 
        required: true,
        condition: { key: 'monitor_type', value: 'specific_items' }
      }
    ]
  },
  {
    id: 'trigger_manual',
    type: 'TRIGGER',
    label: 'Manual Trigger',
    description: 'Manually execute this workflow from the dashboard.',
    icon: 'MousePointerClick',
    color: '!bg-gray-900',
    schema: [
      { key: 'confirmation_message', label: 'Confirmation Message', type: 'TEXT', placeholder: 'Are you sure you want to run this?', description: 'Optional message to show before running.' }
    ]
  },
  {
    id: 'action_n8n_webhook',
    type: 'ACTION',
    label: 'n8n Webhook',
    description: 'Trigger an external n8n workflow.',
    icon: 'Webhook',
    color: '!bg-pink-600',
    schema: [
      { key: 'webhook_url', label: 'Webhook URL', type: 'TEXT', placeholder: 'https://primary.n8n.cloud/webhook/...', required: true },
      { key: 'method', label: 'Method', type: 'SELECT', options: ['GET', 'POST'], required: true },
      { key: 'headers', label: 'Headers', type: 'KEY_VALUE_LIST', description: 'Add custom headers (e.g. Authorization)', required: false }
    ]
  }
]);

const saveWorkflow = () => {

    loading_saving_workflow.value = true;

    if (route.params.id == "") {

      workflow.value.trigger.type = workflow.value.trigger.definitionId
      workflow.value.runs = workflow.value.runs || []

      workflow.value.actions.forEach((action, index) => {
        action.type = action.definitionId
        if (action.definitionId == "action_n8n_webhook"){
          var old_headers = workflow.value.actions[index].headers || []
          workflow.value.actions[index].headers = {}
          old_headers.forEach(header => {
            workflow.value.actions[index].headers[Object.entries(header)[0][0]] = Object.entries(header)[0][1]
          })
        }
      })

      axios.post(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/workflows`, {
        data: workflow.value,
      }, {
          headers: {
              Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
          }
      })
      .then(() => {
          loading_saving_workflow.value = false;
          toast.add({ severity: 'success', summary: 'Workflow Saved', detail: "Successfully saved workflow",group:'br',life:3000 });
      })
      .catch(() => {
          loading_saving_workflow.value = false;
          toast.add({ severity: 'error', summary: 'Error', detail: "Error saving workflow",group:'br',life:3000 }); 
      })

    }else {

      workflow.value.actions.forEach((action, index) => {
        action.type = action.definitionId
        if (action.definitionId == "action_n8n_webhook"){
          var old_headers = workflow.value.actions[index].headers || []
          workflow.value.actions[index].headers = {}
          old_headers.forEach(header => {
            workflow.value.actions[index].headers[Object.entries(header)[0][0]] = Object.entries(header)[0][1]
          })
        }
      })


      axios.patch(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/workflows/${workflow.value.id}`, {
        data: workflow.value,
      }, {
          headers: {
              Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
          }
      })
      .then(() => {
          loading_saving_workflow.value = false;
          toast.add({ severity: 'success', summary: 'Workflow Saved', detail: "Successfully saved workflow",group:'br',life:3000 });
          router.push('/console/workflows');
      })
      .catch(() => {
          loading_saving_workflow.value = false;
          toast.add({ severity: 'error', summary: 'Error', detail: "Error saving workflow",group:'br',life:3000 });
          router.push('/console/workflows');
      })
    }
}

const loadWorkflowData = async (workflowId) => {
    const response = await axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/workflows/${workflowId}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        return response
    })

    return response
}

if (route.params.id != "") {
    loadWorkflowData(route.params.id).then(result => {
        workflow.value = result.data.data
        workflow.value.trigger.definitionId = workflow.value.trigger.type
        workflow.value.trigger.id = `node-${Date.now()}`
        
        // Initialize properties and map trigger-specific fields
        if (!workflow.value.trigger.properties) {
            workflow.value.trigger.properties = {}
        }
        
        // Map trigger properties based on trigger type
        if (workflow.value.trigger.type === 'trigger_low_stock') {
            if (workflow.value.trigger.monitor_type) {
                workflow.value.trigger.properties.monitor_type = workflow.value.trigger.monitor_type
            }
            if (workflow.value.trigger.product_ids) {
                workflow.value.trigger.properties.product_ids = workflow.value.trigger.product_ids
            }
        }
        
        for (let i = 0; i < workflow.value.actions.length; i++) {
            workflow.value.actions[i].definitionId = workflow.value.actions[i].type
            workflow.value.actions[i].id = `node-${Date.now()+i+1}`
            // Initialize properties if it doesn't exist
            if (!workflow.value.actions[i].properties) {
                workflow.value.actions[i].properties = {}

                if (workflow.value.actions[i].type == 'action_n8n_webhook') {
                    workflow.value.actions[i].properties.webhook_url = workflow.value.actions[i].webhook_url
                    workflow.value.actions[i].properties.method = workflow.value.actions[i].method
                    workflow.value.actions[i].properties.headers = []

                    var old_headers = workflow.value.actions[i].headers
                    workflow.value.actions[i].headers = []

                    for (const [key, value] of Object.entries(old_headers)) {
                        workflow.value.actions[i].properties.headers.push({
                            key: key,
                            value: value
                        })
                    }

                    // workflow.value.actions[i].properties.payload = workflow.value.actions[i].payload
                }
            }
        }
    })
}


const loadInventory = async () => {


    const response = await axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/inventories`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        return response
    })

    return response
    
}

loadInventory().then(result => {
    NODE_DEFINITIONS.value[0].schema[1].options = result.data.data.map(item => item.name)
})


// ==========================================
// STATE & PROPS
// ==========================================
const props = defineProps({
  initialWorkflow: {
    type: Object,
    default: () => ({
      id: `wf-${Date.now()}`,
      name: 'New Workflow',
      description: '',
      status: 'draft',
      trigger: null,
      actions: [],
      createdAt: new Date().toISOString()
    })
  }
});

// Local mutable copy of workflow
const workflow = ref(JSON.parse(JSON.stringify(props.initialWorkflow)));

const selectedNodeId = ref(null);
const addingType = ref(null); // 'TRIGGER' | 'ACTION' | null
const isTypeSelectorOpen = ref(false);
const typeSelectorRef = ref(null);

// MultiSelect State
const multiSelectOpen = ref(null);
const multiSelectSearch = ref('');
const multiSelectRefs = ref({});
const multiSelectSearchInput = ref(null);

// ==========================================
// HELPERS
// ==========================================
const getIcon = (name) => LucideIcons[name] || LucideIcons.HelpCircle;

const getDefinition = (defId) => NODE_DEFINITIONS.value.find(d => d.id === defId);

const getPropertiesString = (props, definitionId) => {
    if (!props || Object.keys(props).length === 0) return 'Not configured';
    
    // Special handling for trigger_low_stock
    if (definitionId === 'trigger_low_stock') {
        const parts = [];
        if (props.monitor_type === 'any_item') {
            parts.push('Any item');
        } else if (props.monitor_type === 'specific_items' && props.product_ids && props.product_ids.length > 0) {
            parts.push(`Specific items: ${props.product_ids.join(', ')}`);
        }
        return parts.length > 0 ? parts.join(' | ') : 'Not configured';
    }
    
    return Object.values(props).join(', ');
};

const selectedNode = computed(() => {
    if (!selectedNodeId.value) return null;
    if (workflow.value.trigger && workflow.value.trigger.id === selectedNodeId.value) return workflow.value.trigger;
    return workflow.value.actions.find(a => a.id === selectedNodeId.value) || null;
});

const selectedDef = computed(() => {
    if (!selectedNode.value) return null;
    return getDefinition(selectedNode.value.definitionId);
});

// ==========================================
// ACTIONS
// ==========================================
const handleNodeClick = (node) => {
    selectedNodeId.value = node.id;
    isTypeSelectorOpen.value = false;
    multiSelectOpen.value = null;
};

const handleSetTrigger = (defId) => {
    const newId = `node-${Date.now()}`;
    workflow.value.trigger = {
        id: newId,
        definitionId: defId,
        properties: {}
    };
    selectedNodeId.value = newId;
    isTypeSelectorOpen.value = false;
    multiSelectOpen.value = null;
    addingType.value = null;
};

const handleAddAction = (defId) => {
    if (workflow.value.actions.length >= 1) {
        alert("Only 1 action is allowed.");
        return;
    }
    const newId = `node-${Date.now()}`;
    workflow.value.actions.push({
        id: newId,
        definitionId: defId,
        properties: {}
    });
    selectedNodeId.value = newId;
    isTypeSelectorOpen.value = false;
    multiSelectOpen.value = null;
    addingType.value = null;
};

const handleDeleteNode = (nodeId) => {
    if (workflow.value.trigger && workflow.value.trigger.id === nodeId) {
        workflow.value.trigger = null;
    } else {
        workflow.value.actions = workflow.value.actions.filter(a => a.id !== nodeId);
    }
    if (selectedNodeId.value === nodeId) {
        selectedNodeId.value = null;
    }
};

const handleDefinitionChange = (newDefId) => {
    if (!selectedNode.value) return;
    selectedNode.value.definitionId = newDefId;
    selectedNode.value.properties = {}; // Reset props
};

const updateNodeProperty = (key, value) => {
    if (!selectedNode.value) return;
    selectedNode.value.properties[key] = value;
    
    // Sync properties to workflow object for backend compatibility
    if (workflow.value.trigger && workflow.value.trigger.id === selectedNode.value.id) {
        // Update trigger properties
        if (workflow.value.trigger.definitionId === 'trigger_low_stock') {
            if (key === 'monitor_type') {
                workflow.value.trigger.monitor_type = value;
            } else if (key === 'product_ids') {
                workflow.value.trigger.product_ids = value;
            }
        }
    } else {
        // Update action properties
        const actionIndex = workflow.value.actions.findIndex(a => a.id === selectedNode.value.id);
        if (actionIndex !== -1) {
            if (workflow.value.actions[actionIndex].definitionId === 'action_n8n_webhook') {
                if (key === 'webhook_url') {
                    workflow.value.actions[actionIndex].webhook_url = value;
                } else if (key === 'method') {
                    workflow.value.actions[actionIndex].method = value;
                } else if (key === 'headers') {
                    workflow.value.actions[actionIndex].headers = value.map(header => {
                        return {
                            [header.key] : header.value
                        };
                    });
                }
            }
        }
    }
};

const addKeyValue = (key) => {
    if (!selectedNode.value.properties[key]) {
        selectedNode.value.properties[key] = [];
    }
    selectedNode.value.properties[key].push({ key: '', value: '' });
    updateNodeProperty(key, selectedNode.value.properties[key]);
};

const removeKeyValue = (propKey, index) => {
     if (selectedNode.value.properties[propKey]) {
        selectedNode.value.properties[propKey].splice(index, 1);
        updateNodeProperty(propKey, selectedNode.value.properties[propKey]);
     }
};

const checkCondition = (field, data) => {
    if (!field.condition) return true;
    return data[field.condition.key] === field.condition.value;
};

const handleSave = () => {
    if (!workflow.value.name) return alert("Please give your workflow a name.");
    if (!workflow.value.trigger) return alert("A trigger is required.");
    saveWorkflow()
};

// ==========================================
// MULTI SELECT LOGIC
// ==========================================
const setMultiSelectRef = (el, key) => {
  if (el) multiSelectRefs.value[key] = el;
};

const toggleMultiSelect = (key) => {
  if (multiSelectOpen.value === key) {
    multiSelectOpen.value = null;
  } else {
    multiSelectOpen.value = key;
    multiSelectSearch.value = '';
    // Focus search next tick
    setTimeout(() => {
       if (multiSelectSearchInput.value && typeof multiSelectSearchInput.value.focus === 'function') {
         multiSelectSearchInput.value.focus();
       }
    }, 50);
  }
};

const isSelected = (key, opt) => {
  const vals = selectedNode.value.properties[key] || [];
  return vals.includes(opt);
};

const toggleMultiSelectOption = (key, opt) => {
  const current = selectedNode.value.properties[key] || [];
  let newValue;
  if (current.includes(opt)) {
    newValue = current.filter(v => v !== opt);
  } else {
    newValue = [...current, opt];
  }
  updateNodeProperty(key, newValue);
};

const filterOptions = (options) => {
   if (!options) return [];
   return options.filter(o => o.toLowerCase().includes(multiSelectSearch.value.toLowerCase()));
};

// Close dropdown on outside click
const handleClickOutside = (event) => {
    // Type Selector
    if (typeSelectorRef.value && !typeSelectorRef.value.contains(event.target)) {
        isTypeSelectorOpen.value = false;
    }

    // Multi Select
    if (multiSelectOpen.value) {
       const el = multiSelectRefs.value[multiSelectOpen.value];
       if (el && !el.contains(event.target)) {
           multiSelectOpen.value = null;
       }
    }
};

onMounted(() => {
    document.addEventListener('mousedown', handleClickOutside);
});

onUnmounted(() => {
    document.removeEventListener('mousedown', handleClickOutside);
});

</script>

<style scoped>
.workflow-editor-wrapper {
  all: unset !important;
  display: block !important;
  width: 100% !important;
  height: 100% !important;
  box-sizing: border-box !important;
  font-family: 'Inter', sans-serif !important;
}

.workflow-editor-wrapper * {
  box-sizing: border-box !important;
}

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