<template>
  <div class="w-full">
    <div class="grid mx-2">
        <div class="col-12 flex">
            <div class="gird w-full">
                <div class="col-12 flex w-full justify-content-between align-items-center">
                    <div>
                        <h3 style="font-size:2rem" class="font-bold">Workflows</h3>
                        <p>Automate your workflows effortlessly.</p>
                    </div>
                    <div class="flex px-8 mx-8">
                        <Button label="New workflow" icon="pi pi-plus" @click="$router.push('/console/workflows/put')" />
                    </div>
                </div>
                <div class="col-12 mt-4 flex gap-3 flex-wrap">
                    <div class="w-20rem" v-for="(workflow,index) in workflows" :key="index">
                        <workflowCard @edit="(workflow) => editWorkflow(workflow)" :workflow="workflow" />
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import {getCurrentInstance,ref} from 'vue';
import WorkflowCard from '@/components/WorkflowCard.vue';
import {Button} from 'primevue';
import axios from 'axios';
import { useToast } from "primevue/usetoast";

const toast = useToast()
const workflows = ref([]);

const { proxy } = getCurrentInstance();

const editWorkflow = (workflow) => {
    axios.patch(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/workflows/${workflow.id}`, {
        data: {
            ...workflow,
        }
    }, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then((response)=>{
        toast.add({ severity: 'success', summary: 'Workflow updated successfully!', detail: response.data.data,group:'br' });
        getWorkflows();
    })
    .catch((err) => {
        toast.add({ severity: 'error', summary: 'Failed updating the workflow!', detail: err.data,group:'br' });
        console.log(err)
    });
}

const getWorkflows = () => {
   axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/workflows`, {
       headers: {
           Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
       },
   })
   .then((response)=>{
    console.log(response.data.data)
    workflows.value = response.data.data;
   })
   .catch((err) => {
       console.log(err)
   });
}

getWorkflows();

</script>