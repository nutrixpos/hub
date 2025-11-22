<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3 style="font-size:2rem">Pipelines</h3>
                    </div>
                    <div class="col-12">
                        <Button icon="pi pi-plus" class="w-auto" label="Add pipeline" @click="new_pipelines_dialog=true"></Button>
                    </div>
                    <div class="col-12">
                        <Panel header="Pipelines">
                            <div class="flex gap-3 align-items-center">
                                <p style="font-size:1.2rem;">#1</p>
                                <Breadcrumb :model="pipelines_recap">
                                    <template #item="{ item }">
                                        <div class="cursor-pointer flex gap-2 align-items-center">
                                            <span :class="item.icon"></span>
                                            <span>{{ item.label }}</span>
                                        </div>
                                    </template>
                                </Breadcrumb>
                                <ButtonGroup>
                                    <Button icon="pi pi-pencil" size="small" severity="secondary" />
                                    <Button icon="pi pi-times" size="small" severity="secondary" />
                                </ButtonGroup>
                            </div>
                            <Divider />
                        </Panel>
                        <Dialog v-model:visible="new_pipelines_dialog" modal header="Add new pipeline" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                            <div class="flex flex-column align-items-start gap-4">
                                <div class="flex align-items-center gap-3">
                                    <label for="new_pipeline_name" class="font-semibold">Name</label>
                                    <InputText id="new_pipeline_name" v-model="new_pipeline_name" class="flex-auto" />
                                </div>
                                <div class="flex gap-3 align-items-center">
                                    <span class="fa fa-bolt" style="color:gold"></span> <p style="font-size:1.5rem;font-weight:bold;">Event:</p>
                                    <Select v-model="selectedEvent" :options="events" optionLabel="name" placeholder="Select" class="w-full md:w-56" />
                                </div>
                                <event-inventory-low-stock/>
                                <div class="flex gap-3 align-items-center">
                                    <span class="fa fa-rocket" style="color:cornflowerblue"></span> <p style="font-size:1.5rem;font-weight:bold">Action:</p>
                                    <Select v-model="selectedAction" :options="actions" optionLabel="name" placeholder="Select" class="w-full md:w-56" />
                                </div>
                                <action-n8n />
                            </div>
                            <template #footer>
                                <div class="mx-5 my-3 gap-1 flex">
                                    <Button label="Cancel" text severity="secondary" />
                                    <Button label="Test" @click="testPipeline" severity="warn" class="w-5rem"  />
                                    <Button label="Save" class="w-5rem"/>
                                </div>
                            </template>
                        </Dialog>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import {Button, Panel, Breadcrumb, Divider, Dialog, Select, InputText, ButtonGroup} from 'primevue';
import EventInventoryLowStock from '../components/EventInventoryLowStock.vue';
import ActionN8n from '../components/ActionN8n.vue';
import axios from 'axios';

const new_pipelines_dialog = ref(false);
const new_pipeline_name = ref("")

const selectedEvent = ref({ name: 'Inventory Low stock', code: 'inventory_low_stock' });
const events = ref([
    { name: 'Inventory Low stock', code: 'inventory_low_stock' },
]);

const pipelines_recap = ref([
    { label: 'low_stock', icon:'fa fa-bolt' },
    { label: 'n8n_webhook', icon:'fa fa-rocket' },
]);


const selectedAction = ref({ name: 'n8n', code: 'n8n' })
const actions = ref([
    { name: 'n8n_webhook', code: 'n8n_webhook' },
]);


const testPipeline = () => {
    console.log("Testing pipeline...")
    axios.post("https://nutrixpos.app.n8n.cloud/webhook-test/a324955c-ea3f-41a2-92ac-5c34e61b56df", {
        "name": "Elmawardy",
    })
    .then(response => {
        console.log(response.data);
    })
    .catch(error => {
        console.log(error);
    });
}

</script>