<template>
    <div class="w-full h-full flex flex-column" style="height: calc(100vh - 6rem);"> <!-- Adjust height to fit within layout -->
        <!-- Header -->
        <div class="flex justify-content-between align-items-center px-4 py-2 border-bottom-1 surface-border">
            <div class="flex align-items-center gap-2">
                <h3 class="m-0 font-bold text-xl">Koptan AI</h3>
                <Tag value="Beta" severity="info" class="text-xs"></Tag>
            </div>
        </div>

        <!-- Main Content Area -->
        <div class="flex-grow-1 flex flex-column overflow-hidden relative">
            
            <!-- Free Plan State -->
            <div v-if="store.subscription.subscription_plan == 'free'" class="flex flex-column align-items-center justify-content-center h-full gap-4">
                <i class="pi pi-lock text-5xl text-400"></i>
                <div class="text-center">
                    <h2 class="mt-0 mb-2">Unlock Koptan AI to full potential</h2>
                    <p class="text-color-secondary m-0">Upgrade to Gold plan to access personalized suggestions and chat.</p>
                </div>
                <Button label="Upgrade to Gold" icon="pi pi-star-fill" severity="warning" @click="$router.push('/console/subscription')" />
            </div>

            <!-- Chat Interface (Gold Plan) -->
            <div v-else class="flex flex-column h-full">
                <!-- Messages Area -->
                <div class="flex-grow-1 overflow-y-auto p-4 flex flex-column gap-4">
                    
                    <!-- Welcome / Intro Message -->
                    <div class="flex gap-3">
                        <Avatar icon="pi pi-sparkles" shape="circle" class="bg-primary-50 text-primary" />
                        <div class="flex flex-column gap-1" style="max-width: 80%;">
                            <span class="font-bold text-sm">Koptan</span>
                            <div class="p-3 border-round-2xl surface-card shadow-1 line-height-3">
                                Hello! I'm Koptan, your AI assistant. Here are some suggestions based on your recent sales:
                            </div>
                        </div>
                    </div>

                    <div v-for="(chat, index) in chats" :key="index" class="flex gap-3 fadein animation-duration-500">
                        <Avatar icon="pi pi-sparkles" shape="circle" v-if="chat.source != 'You'" class="`bg-primary-50 text-primary" />
                        <div :class="`flex flex-column gap-1 ${chat.source == 'Koptan' ? '' : 'ml-auto'}`" style="max-width: 80%;">
                            <span class="font-bold text-sm" v-if="chat.source != 'You'">{{chat.source}}</span>
                            <div class="p-3 border-round-2xl surface-card shadow-1 line-height-3">
                                <div class="font-medium mb-1 text-primary-600 text-xs uppercase" v-if="chat.type">{{ chat.type }}</div>
                                {{ chat.content }}
                            </div>
                        </div>
                    </div>

                    <div v-if="is_loading" class="flex gap-3">
                        <Avatar icon="pi pi-sparkles" shape="circle" class="bg-primary-50 text-primary" />
                        <div class="flex align-items-center p-3">
                             <ProgressSpinner style="width:1.5rem; height: 1.5rem;"/>
                        </div>
                    </div>
                    
                </div>

                <!-- Input Area -->
                <div class="px-4 mb-6 pb-4 pt-2">
                    <div class="flex flex-column align-items-center w-full">
                        <div class="w-full relative shadow-1 border-round-3xl surface-100" style="max-width: 800px;">
                             <InputText placeholder="Reply to Koptan..." class="w-full border-none py-3 pl-4 pr-6 text-lg shadow-none outline-none" v-model="userInput" style="border-radius: inherit;" />
                             <div class="absolute right-0 top-0 h-full flex align-items-center pr-2">
                                <Button icon="pi pi-send" rounded text class="text-color-secondary hover:text-primary transition-colors" @click="addUserChat(userInput)" :loading="is_loading"/>
                             </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import {getCurrentInstance, ref, onMounted, watch} from 'vue'
import axios from 'axios'
import { globalStore } from '../stores';
import Avatar from 'primevue/avatar';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import ProgressSpinner from 'primevue/progressspinner';

const chats = ref<any[]>([])

const suggestions = ref([])
const is_loading = ref(true)

// Keep pagination logic implicitly for fetching, or just fetch a clean batch
// Since it's a chat flow, maybe we just load the first batch for now.
const pageNumber = ref(1)
const pageSize = ref(10) 

const userInput = ref('')

const {proxy} = getCurrentInstance()
const store = globalStore()

const loadSuggestions = () => {
    is_loading.value = true
    axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/koptan/suggestions?page[number]=${pageNumber.value}&page[size]=${pageSize.value}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        suggestions.value = response.data.data

        response.data.data.forEach((suggestion:any) => {
            chats.value.push({
                type: suggestion.type,
                content: suggestion.content,
                source: "Koptan"
            })
        })

        is_loading.value = false
    })
    .catch(() => {
        is_loading.value = false
    })
}

const addUserChat = (content:string) => {

    is_loading.value = true

    chats.value.push({
        content: content,
        source: "You"
    })
}


if (store.subscription.subscription_plan == 'gold') {
    loadSuggestions()
}

watch(
  () => store.subscription,
  (newValue, _) => {
    if (newValue.subscription_plan == 'gold') {
        loadSuggestions()
    }
  }
)

</script>

<style scoped>
/* Add any custom transitions or overrides here */
.surface-card {
    background-color: var(--surface-card);
}
</style>