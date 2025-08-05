<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('inventory')}}</h3>
                    </div>
                    <div class="col-12 flex flex-column gap-3 w-full">
                        <h5>Suggestions</h5>
                        <div class="flex flex-column w-full">
                            <DataTable @page="updatSalesTableRowsPerPage" :lazy="true" :totalRecords="inventoryTableTotalRecords" :loading="isInventoryTableLoading" paginatorPosition="both"  paginator :rows="inventoryTableRowsPerPage" :rowsPerPageOptions="[7, 14, 30, 90]" :value="inventory_items" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                                    <Column sortable field="content" :header="$t('content')"></Column>
                                    <Column sortable field="type" :header="$t('type')"></Column>
                            </DataTable>
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
import {Badge, Chip} from 'primevue';


const inventoryTableRowsPerPage = ref(7)
const inventoryTableTotalRecords = ref(0)
const inventory_items = ref([])
const isInventoryTableLoading = ref(true)
const salesTableFirstIndex = ref(0)

const {proxy} = getCurrentInstance()


const updatSalesTableRowsPerPage = (event) => {

    const { first, rows } = event;
    loadInventory(first,rows)
}


const loadInventory = (first=salesTableFirstIndex.value,rows=inventoryTableRowsPerPage.value) => {


    let page_number = Math.floor(first/rows) + 1
    isInventoryTableLoading.value = true


    axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/inventories`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        inventory_items.value = response.data.data
        inventoryTableTotalRecords.value = response.data.length
        isInventoryTableLoading.value = false
    })
    
}

loadInventory()

</script>