<template>
    <div v-if="!loading" class="h-full">
        <div class="grid p-0 m-0 h-full">
            <div class="col-12 p-0">
                <Toolbar style="border-radius: 0px;flex-shrink: 0;background-color:#001529;border: 0px;color:white;" class="py-1 lg:py-2">
                    <template #start>
                        <router-link to="/">
                            <span style="color:var(--p-button-secondary-background) !important;" class="text-xl font-bold">nutrixhub</span>
                        </router-link>
                    </template>
                </Toolbar>
            </div>
            <div class="col-12 h-full">
                <div class="grid h-full">
                    <div class="col-3 xl:col-2 p-0 m-0 h-full">                        
                        <Menu :model="items" class="w-full h-full md:w-60" style="border-radius: 0px;">
                            <template #submenulabel="{ item }">
                                <span class="text-primary font-bold">{{ item.label }}</span>
                            </template>
                            <template #item="{ item, props }">
                                <router-link class="flex items-center" :to="item.link" v-bind="props.action">
                                    <span :class="item.icon" />
                                    <span>{{ item.label }}</span>
                                    <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
                                    <span v-if="item.shortcut" class="ml-auto border border-surface rounded bg-emphasis text-muted-color text-xs p-1">{{ item.shortcut }}</span>
                                </router-link>
                            </template>
                            <template #end>
                                <button v-ripple class="flex align-items-end relative overflow-hidden w-full border-0 bg-transparent flex items-start p-2 pl-4 hover:bg-surface-100 dark:hover:bg-surface-800 rounded-none cursor-pointer transition-colors duration-200">
                                    <Avatar icon="pi pi-user" class="mr-2" shape="circle" />
                                    <span class="inline-flex flex-col items-start">
                                        <span class="font-bold">{{ user?.name }}</span>
                                        <span class="text-sm">
                                            <Chip v-for="(role,index) in roles" :key="index" :label="role" style="height: 1.5rem;" class="m-1" />
                                        </span>
                                    </span>
                                </button>
                                <Button class="w-full text-start" icon="pi pi-sign-out" severity="secondary" text aria-label="Signout" :label=" $t('signout')" @click="proxy.$zitadel?.oidcAuth.signOut()" />                            
                            </template>
                        </Menu>


                    </div>
                    <div class="col-9 xl:col-10 flex p-0 pt-3" style="background-color: white;background-color: var(--p-gray-100);">
                        <RouterView />
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div style="width:100vw;height:100vh;display:flex;justify-content:center;align-items:center" v-if="loading">
      <ProgressSpinner style="width: 35px; height: 35px;" strokeWidth="4" fill="transparent"
      animationDuration=".5s" aria-label="Custom ProgressSpinner" />
    </div>
</template>

<script setup lang="ts">
import {ref,getCurrentInstance,computed} from "vue";
import { Toolbar } from "primevue";
import Tree from "primevue/tree";
import Button from "primevue/button";
import { useI18n } from 'vue-i18n'
import { globalStore } from '@/stores';
import axios from "axios";
import ProgressSpinner from "primevue/progressspinner";
import {Menu,Avatar} from 'primevue';


const { proxy } = getCurrentInstance();
const store = globalStore()
const user_profile_op = ref();

const user : any = computed(() => {

    return proxy.$zitadel?.oidcAuth.userProfile

})

const sidemenuNodeSelect = (node) => {
    if (node.link) {
        proxy.$router.push(node.link);
    }
}


// const selected_list_item = ref ({ name: 'Inventory', icon:'inbox', link:'inventory' })

const user_profile_toggle = (event: any) => {
    user_profile_op.value.toggle(event);
}

const expandNode = (node) => {
    if (node.children && node.children.length) {
        expandedKeys.value[node.key] = true;

        for (let child of node.children) {
            expandNode(child);
        }
    }
};

const expandedKeys = ref({});

const items = ref([
    {
        separator: false
    },
    {
        label: 'Observe',
        items: [
            {
                label: 'Sales',
                icon: 'pi pi-chart-line',
                link: '/console/sales',
            },
        ]
    },
    {
        separator: true
    }
]);

  const loading = ref(true)
const { locale,setLocaleMessage } = useI18n({ useScope: 'global' })

const loadLanguage = async () => {

    await axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/settings`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then(async (response)=>{
        await axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/languages/${response.data.data.language.code}`, {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        })
        .then(response2 => {

            setLocaleMessage(response2.data.data.code, response2.data.data.pack);
            locale.value = response2.data.data.code
            store.setOrientation(response2.data.data.orientation)
            loading.value = false
        })
        .catch((err) => {
            console.log(err)
        });
        loading.value = false
    })
    .catch((err) => {
        console.log(err)
    });

}


loadLanguage()
</script>

<style>
html,
body {
height: 100%;
margin:0;
}

.p-progressspinner-circle {
    stroke: black !important;
}

.p-menu-start{
    display:none
}

#pv_id_5_0 {
    display: none;
}
</style>