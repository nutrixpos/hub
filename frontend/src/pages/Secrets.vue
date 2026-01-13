<template>
<div class="w-full">
    <div class="grid mx-2">
        <div class="col-12 flex">
            <div class="gird w-full">
                <div class="col-12 flex w-full justify-content-between align-items-center">
                    <div>
                        <h3 style="font-size:2rem" class="font-bold">Secrets and variables</h3>
                        <p>Add secrets and variables to your workflows</p>
                    </div>
                    <div class="flex px-8 mx-8">
                        <Button label="New environment variable" icon="pi pi-plus" @click="displayModal = true" />
                    </div>
                </div>
                <div class="col-12 flex justify-content-center align-items-center w-full">
                    <DataTable :value="vars" class="w-full mt-4">
                        <Column field="name" header="Name">
                            <template #body="{ data }">
                                <InputText v-model="data.name" class="w-full" />
                            </template>
                        </Column>
                        <Column field="value" header="Value">
                            <template #body="{ data }">
                                <InputText v-model="data.value" class="w-full" :type="data.type == 'secret' ? 'password' : 'text' " />
                            </template>
                        </Column>
                        <Column field="type" header="Type" style="width: 15rem">
                            
                            <template #body="{ data }">
                                {{data.type}}
                            </template>
                        </Column>
                        <Column :header="$t('actions')" style="width:30rem">
                            <template #body="{data}">
                                <ConfirmPopup></ConfirmPopup>
                                <ButtonGroup>
                                    <Button icon="pi pi-save" severity="secondary" aria-label="Save" @click="patchVar(data.name,data)" />
                                    <Button icon="pi pi-trash" severity="danger" aria-label="Delete" @click="confirmDeleteVar($event,data.name)" />
                                </ButtonGroup>
                            </template>
                        </Column>
                    </DataTable>
                </div>
            </div>
        </div>
    </div>
    <Dialog v-model:visible="displayModal" header="New Variable" :modal="true" :style="{ width: '30vw' }">
        <div class="flex flex-column gap-3">
            <div class="flex flex-column gap-2">
                <label for="name">Name</label>
                <InputText id="name" v-model="newVar.name" />
            </div>
            <div class="flex flex-column gap-2">
                <label for="type">Type</label>
                <Select v-model="newVar.type" :options="['secret', 'var']" placeholder="Select a type" class="w-full" />
            </div>
            <div class="flex flex-column gap-2">
                <label for="value">Value</label>
                <InputText id="value" v-model="newVar.value" :type="newVar.type === 'secret' ? 'password' : 'text'" />
            </div>
        </div>
        <template #footer>
            <Button label="Cancel" icon="pi pi-times" text @click="displayModal = false" />
            <Button label="Save" icon="pi pi-check" @click="saveNewVar" />
        </template>
    </Dialog>
</div>
</template>

<script setup>
import {ref, getCurrentInstance} from 'vue'
import {DataTable,Column, InputText, Select, ButtonGroup, Button, Dialog, ConfirmPopup, useConfirm} from 'primevue'
import axios from 'axios'
import {useToast} from 'primevue/usetoast';
const toast = useToast()
const {proxy} = getCurrentInstance()
const confirm = useConfirm();

const vars = ref([])

const displayModal = ref(false)
const newVar = ref({
    name: '',
    value: '',
    type: 'var'
})


const confirmDeleteVar = (event,var_name) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to delete this variable ?',
        icon: 'pi pi-exclamation-triangle',
        rejectProps: {
            label: 'Cancel',
            severity: 'secondary',
            outlined: true
        },
        acceptProps: {
            label: 'Yes'
        },
        accept: () => {


            deleteVar(var_name)

         
        },
        reject: () => {
        }
    });
  }


const patchVar = (var_name, data) => {
    axios.patch(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/env_vars/${var_name}`, 
        {
            data: data,
        },
        {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
    
        toast.add({ 
            severity: 'success', 
            summary: 'Done', 
            detail: response.data.message || 'Successfully updated environment variable !', 
            group: 'br' 
        });
        loadVars()
    })
    .catch(error => {
        toast.add({ 
            severity: 'error', 
            summary: 'Failed', 
            detail: error.response.data.message || 'Failed to update environment variable !', 
            group: 'br' 
        });
    })
}

const deleteVar = (var_name) => {
    axios.delete(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/env_vars/${var_name}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        toast.add({ 
            severity: 'success', 
            summary: 'Done', 
            detail: response.data.message || 'Successfully deleted environment variable !', 
            group: 'br' 
        });
        loadVars()
    })
    .catch(error => {
        toast.add({ 
            severity: 'error', 
            summary: 'Failed', 
            detail: error.response.data.message || 'Failed to delete environment variable !', 
            group: 'br' 
        });
    })
}

const loadVars = () => {

    vars.value = []

    axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/env_vars`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        vars.value = response.data.data.map((item) => {
            return {
                name: item.name,
                value: item.value,
                type: item.is_secret ? 'secret' : 'var'
            }
        }) || []
    })
    .catch(error => {
        toast.add({ 
            severity: 'error', 
            summary: 'Failed', 
            detail: error.response.data.message || 'Failed to load environment variables !', 
            group: 'br' 
        });
    })
}

loadVars()

const saveNewVar = () => {
    displayModal.value = false

    axios.patch(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/env_vars/${newVar.value.name}`, 
        {
            data: {
                name: newVar.value.name,
                value: newVar.value.value,
                is_secret: newVar.value.type === 'secret'
            },
        },
        {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
    
        toast.add({ 
            severity: 'success', 
            summary: 'Done', 
            detail: response.data.message || 'Successfully added new environment variable !', 
            group: 'br' 
        });
        displayModal.value = false
        loadVars()
    })
    .catch(error => {
        toast.add({ 
            severity: 'error', 
            summary: 'Failed', 
            detail: error.response.data.message || 'Failed to insert new environment variable !', 
            group: 'br' 
        });
    })
}
</script>