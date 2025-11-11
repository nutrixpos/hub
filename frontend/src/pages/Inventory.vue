<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3 style="font-size:2rem">{{$t('inventory')}}</h3>
                    </div>
                    <div class="col-12 flex flex-column gap-3 w-full">
                        <div class="flex flex-column w-full">
                            <DataTable @page="updatSalesTableRowsPerPage" :lazy="true" :totalRecords="inventoryTableTotalRecords" :loading="isInventoryTableLoading" :value="inventory_items" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                                    <Column sortable field="name" :header="$t('name')"></Column>
                                    <Column sortable field="quantity" :header="$t('quantity')"></Column>
                                    <Column field="labels" header="Labels">
                                        <template #body="slotProps">
                                            <div class="flex gap-2">
                                                <Chip v-for="(label,index) in slotProps.data.labels" :key="index" :label="label" style="height: 1.5rem;" class="m-1" />
                                            </div>
                                        </template>
                                    </Column>
                                    <Column sortable field="unit" :header="$t('unit')"></Column>
                                    <Column :header="$t('status')">
                                        <template #body="slotProps">
                                            <Tag :value="slotProps.data.quantity > slotProps.data.settings.alert_threshold ? 'INSTOCK' : 'LOWSTOCK'" :severity="slotProps.data.quantity > slotProps.data.settings.alert_threshold ? 'success' : 'danger'" />
                                        </template>
                                    </Column>
                                    <Column :header="$t('actions')" style="width:30rem">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-cog" severity="secondary" aria-label="Settings" @click="material_settings = slotProps.data;  alert_threshold = slotProps.data.settings.alert_threshold; material_settings_dialog=true"  />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #footer v-if="total_inventory_items_count_in_db > 1 && store.subscription.subscription_plan == 'free'">
                                        <div class="free text-center flex-column">
                                            <p>{{ total_inventory_items_count_in_db - 1 }} more items/s to show</p>
                                            <Button class="mt-2" style="background-color:#E1C05C;border-color:gold;color:black">
                                                <RouterLink to="/console/subscription">Upgrade to GOLD to unlock</RouterLink>
                                            </Button>
                                        </div>
                                    </template>
                            </DataTable>
                            <Dialog v-model:visible="material_settings_dialog" modal :header="`Settings for  ${material_settings?.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                                <div class="flex align-items-center">
                                    <h4>stock_alert_treshold</h4>
                                    <InputText type="number" class="ml-2" id="stock_alert_treshold" v-model.number="alert_threshold" aria-describedby="stock_alert_treshold" />
                                </div>
                                <template #footer>
                                    <Button label="Cancel" @click="material_settings_dialog=false" severity="secondary" aria-label="Save"  />
                                    <Button severity="primary" label="Save" aria-label="Save" @click="patchInventory"/>
                                </template>
                            </Dialog>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import DataTable from "primevue/datatable";
import Column from 'primevue/column'
import {getCurrentInstance, ref} from 'vue'
import axios from 'axios'
import { $dt } from '@primevue/themes';
import {Badge, Button, ButtonGroup, Chip, Tag, Dialog, InputText} from 'primevue';
import { globalStore } from '../stores';


import {useToast} from 'primevue/usetoast';
const toast = useToast()
const store = globalStore()


const inventory_items = ref([])
const isInventoryTableLoading = ref(true)
const material_settings_dialog = ref(false)
const material_settings = ref({})
const total_inventory_items_count_in_db = ref(0)

const alert_threshold = ref(0)

const {proxy} = getCurrentInstance()


const loadInventory = () => {


    isInventoryTableLoading.value = true


    axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/inventories`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        inventory_items.value = response.data.data
        total_inventory_items_count_in_db.value = response.data.meta.total_records
        isInventoryTableLoading.value = false
    })
    
}


const patchInventory = () => {
        axios.patch(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/inventories`, 
        {
            data: {
                alert_threshold: alert_threshold.value
            },
            meta: {
                    id: material_settings.value.id,
            }
        },
        {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        toast.add({ 
            severity: 'success', 
            summary: 'Inventory Updated', 
            detail: response.data.message || 'Inventory updated successfully!', 
            group: 'br' 
        });
        material_settings_dialog.value = false
        loadInventory()
    })
    .catch(error => {
        toast.add({ 
            severity: 'error', 
            summary: 'Inventory Update Failed', 
            detail: error.response.data.message || 'Failed to update inventory!', 
            group: 'br' 
        });
    })
}

loadInventory()

</script>