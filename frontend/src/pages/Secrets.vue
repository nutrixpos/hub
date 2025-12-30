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
                        <Button label="New workflow" icon="pi pi-plus" @click="$router.push('/console/workflows/put')" />
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
                                <InputText v-model="data.value" class="w-full" type="password" />
                            </template>
                        </Column>
                        <Column field="type" header="Type" style="width: 15rem">
                            <template #body="{ data }">
                                <Select
                                    v-model="data.type"
                                    :options="['string', 'number', 'boolean', 'object', 'array', 'null', 'undefined']"
                                    optionLabel="name"
                                    optionValue="value"
                                    placeholder="Select type"
                                    class="w-full"
                                />
                            </template>
                        </Column>
                    </DataTable>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<script setup>
import {ref} from 'vue'
import {DataTable,Column, InputText, Select} from 'primevue'


const vars = ref([
    {
        name: 'name',
        value: 'value',
        type: 'string'
    }
])
</script>