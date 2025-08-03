<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('koptan')}}</h3>
                    </div>
                    <div class="col-12 flex flex-column gap-3 w-full">
                        <h5>Suggestions</h5>
                        <div class="flex flex-column w-full">
                            <DataTable @page="updatSalesTableRowsPerPage" :lazy="true" :totalRecords="suggestionsTableTotalRecords" :loading="isSuggestionsTableLoading" paginatorPosition="both"  paginator :rows="salesTableRowsPerPage" :rowsPerPageOptions="[7, 14, 30, 90]" :value="suggestions" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
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


const salesTableRowsPerPage = ref(7)
const suggestionsTableTotalRecords = ref(0)
const suggestions = ref([])
const isSuggestionsTableLoading = ref(true)
const salesTableFirstIndex = ref(0)

const {proxy} = getCurrentInstance()


const updatSalesTableRowsPerPage = (event) => {

    const { first, rows } = event;
    loadSuggestions(first,rows)
}


const loadSuggestions = (first=salesTableFirstIndex.value,rows=salesTableRowsPerPage.value) => {


    let page_number = Math.floor(first/rows) + 1
    isSuggestionsTableLoading.value = true


    axios.get(`${import.meta.env.VITE_APP_BACKEND_HOST}/${import.meta.env.VITE_APP_BACKEND_VERSION}/api/koptan/suggestions?page[number]=${page_number}&page[size]=${rows}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        suggestions.value = response.data.data
        suggestionsTableTotalRecords.value = response.data.meta.total_records
        isSuggestionsTableLoading.value = false
    })
    
}

loadSuggestions()

</script>