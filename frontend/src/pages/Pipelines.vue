<template>
  <div class="w-full">
    <div class="grid mx-2">
        <div class="col-12 flex">
            <div class="gird w-full">
                <div class="col-12 flex w-full justify-content-between align-items-center">
                    <div>
                        <h3 style="font-size:2rem" class="font-bold">Pipelines</h3>
                        <p>Automate your workflows effortlessly.</p>
                    </div>
                    <div class="flex px-8 mx-8">
                        <Button label="New workflow" icon="pi pi-plus" @click="$router.push('/pipelines/put')" />
                    </div>
                </div>
                <div class="col-12 mt-4">
                    <div class="w-20rem">
                        <workflowCard :workflow="{status:'active',actions:[{id:'123',definitionId:'123'}]}" />
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import WorkflowCard from '@/components/WorkflowCard.vue';
import {Button} from 'primevue';
</script>