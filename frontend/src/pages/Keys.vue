<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('API Keys')}}</h3>
                    </div>
                    <div class="col-12 flex flex-column">
                        <div class="col-12 flex justify-content-end">
                            <Button label="Add key" icon="pi pi-plus" @click="$router.push({name:'newkey'})" />
                        </div>
                        
                        <Card>
                            <template #title>Branch #1</template>
                            <template #content>
                                <p class="m-0">
                                    Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis, culpa ratione quam perferendis esse, cupiditate neque
                                    quas!
                                </p>
                            </template>
                        </Card>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import {Card, Button} from 'primevue'

</script>