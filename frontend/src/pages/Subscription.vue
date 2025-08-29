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
                                    <Badge class="ml-auto" :value="`${store.subscription.subscription_plan.charAt(0).toUpperCase()}${store.subscription.subscription_plan.slice(1)}`" size="xlarge" :style="`background-color: ${store.subscription.subscription == 'free' ?'silver' : 'var(--p-primary-color)'};color:${store.subscription.subscription == 'free' ?'black' : 'var(--p-button-primary-color)'}`" />
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
                            <Card class="w-20rem" v-if="store.subscription.subscription_plan != 'standard'">
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
                                    <Button label="Subscribe" severity="primary" class="w-full" @click="requestSusbcription('standard')" />
                                </template>
                            </Card>
                            <Card class="w-20rem" v-if="store.subscription.subscription_plan != 'gold'">
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
                                    <Button label="Subscribe" severity="primary" class="w-full" style="background-color:#E1C05C;border-color:gold;color:black" @click="requestSusbcription('gold')" />
                                </template>
                            </Card>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <Dialog v-model:visible="loading_subscription_request" pt:root:class="!border-0 !bg-transparent" pt:mask:class="backdrop-blur-sm">
            <template #container>
                <Card>
                    <template #content>
                        <div class="flex flex-column gap-3 justify-content-center align-items-center p-3">
                            <p class="m-0">
                                {{ subscrition_request_message }}
                            </p>
                            <ProgressSpinner style="width: 35px; height: 35px;" strokeWidth="4" fill="transparent" animationDuration=".5s" aria-label="Custom ProgressSpinner" />
                        </div>
                    </template>
                </Card>
            </template>
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import {getCurrentInstance,ref} from 'vue'
import {Card,Badge, Button,Dialog, ProgressSpinner} from 'primevue';
import axios from 'axios';
import { useToast } from "primevue/usetoast";
import { globalStore } from '@/stores';


const store = globalStore()


const loading_subscription_request = ref(false);
const toast = useToast();
const { proxy } = getCurrentInstance();
const subscrition_request_message = ref("Requesting subscription...");

const requestSusbcription = (plan: string) => {
    loading_subscription_request.value = true;
    axios.post(`${import.meta.env.VITE_APP_BACKEND_HOST}/v1/api/subscriptions/request`, 
       {
           data: {
               plan: plan,
           }
       },
       {
           headers: {
               Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
           }
       }
   )
   .then((response)=>{
       console.log("response: "+response.data);
       subscrition_request_message.value = "Subscription request successful. Redirecting to the payment gateway...";
       setTimeout(() => {
           window.location.href = response.data.data.payment_url;
       }, 1000);       
   })
   .catch((err) => {
       toast.add({ severity: 'error', summary: 'Failed', detail: err, life: 3000,group:'br' });  
       console.log(err)
       loading_subscription_request.value = false;
   });
}

</script>