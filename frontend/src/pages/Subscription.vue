<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h1 style="font-size:2rem;">Current Subscription</h1>
                    </div>
                    <div class="col-12 flex flex-column gap-3 w-full">
                        <div class="flex flex-column w-full">
                            <Card>
                                <template #title>
                                    <Badge class="ml-auto" value="Free" size="xlarge" style="background-color: silver;color:black" />
                                </template>
                                <template #content>
                                    <p class="m-0">
                                        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis, culpa ratione quam perferendis esse, cupiditate neque
                                        quas!
                                    </p>
                                </template>
                            </Card>
                        </div>
                        <h1 style="font-size:2rem;" class="mt-4 ">Other options</h1>
                        <div class="flex gap-3">
                            <Card class="w-20rem">
                                <template #title>
                                    <Badge class="ml-auto" value="Standard" size="xlarge" />
                                </template>
                                <template #content>
                                    <div class="flex-column flex gap-3 mt-3">
                                        <p class="text-center" style="font-size:1.5rem">$30/month</p>
                                        <p class="m-0 p-2">
                                            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis, culpa ratione quam perferendis esse, cupiditate neque
                                            quas!
                                        </p>
                                    </div>
                                </template>
                                <template #footer>
                                    <Button label="Subscribe" severity="primary" class="w-full" />
                                </template>
                            </Card>
                            <Card class="w-20rem">
                                <template #title>
                                    <Badge class="ml-auto" value="GOLD" size="xlarge" style="background-color: #E1C05C;color:black" />
                                </template>
                                <template #content>
                                    <div class="flex-column flex gap-3 mt-3">
                                        <p class="text-center" style="font-size:1rem;text-decoration: line-through;line-height: 0px;margin-top:-0.5rem;">$60/month</p>
                                        <p class="text-center" style="font-size:1.5rem;line-height: 1rem;">$50/month</p>
                                        <p class="my-2 p-2">
                                            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis, culpa ratione quam perferendis esse, cupiditate neque
                                            quas!
                                        </p>
                                    </div>
                                </template>
                                <template #footer>
                                    <Button label="Subscribe" severity="primary" class="w-full" style="background-color:#E1C05C;border-color:gold;color:black" />
                                </template>
                            </Card>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import {ref} from "vue";
import {Card,Badge, Button} from 'primevue';
</script>