<template>
        <Card class="w-full">
            <template #title><span class="fa fa-rocket"></span> N8n webhook action</template>
            <template #content>
                <div class="flex align-items-center gap-3 w-full">
                    <label for="link" class="font-semibold">Webhook link</label>
                    <InputText id="webhook_link" v-model="webhook_link" class="flex-auto" />
                </div>
            </template>
        </Card>
</template>

<script setup>
import { ref } from 'vue';
import { InputText, Card } from 'primevue';

const webhook_link = ref('https://nutrixpos.app.n8n.cloud/webhook-test/a324955c-ea3f-41a2-92ac-5c34e61b56df');
</script>