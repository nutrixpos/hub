<template>
    <Card class="w-full">
        <template #title><span class="fa fa-bolt"></span> Low Inventory Stock</template>
        <template #content>
            <div class="flex flex-wrap gap-3 flex-column">
                <div class="flex items-center gap-2">
                    <RadioButton v-model="items_selection_type" inputId="anything" name="items_selection_type" value="anything" />
                    <label for="anything">Any Item</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton v-model="items_selection_type" inputId="selected" name="items_selection_type" value="selected_items" />
                    <label for="selected">Selected Items</label>
                </div>
            </div>
            <div class="flex gap-2 align-items-center px-5">
                <p style="  white-space: nowrap;">Triggering items: </p>
                <MultiSelect :disabled="items_selection_type != 'selected_items'" v-model="selected_items" display="chip" :options="available_items" optionLabel="name" filter placeholder="Select" :maxSelectedLabels="3" class="w-full md:w-80" />
            </div>
        </template>
    </Card>
</template>

<script setup>
import { ref } from "vue";
import {Card, MultiSelect, RadioButton } from 'primevue';

const items_selection_type = ref('anything');

const selected_items = ref();
const available_items = ref([
    { name: 'Motzarilla', id: '67adfa6ccadb32c437c44e09' },
    { name: 'Olive', id:'681340bb4089da41c24cec6d' },
]);



</script>