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
                            <Card v-if="store.subscription.subscription_plan == 'free'" class="w-20rem">
                                <template #title>
                                    <div class="flex gap-1 flex-column justify-content-stasrt align-items-start">
                                        <Badge size="xlarge" :value="store.subscription.subscription_plan.toUpperCase()"  :style="`background-color: ${store.subscription.subscription_plan.toUpperCase() == 'FREE' ?'silver' : '#E1C05C'};color:black`"/>
                                        <p class="my-2 p-2" style="font-size:1rem;font-weight: 400;">
                                            Get started with the basics. See your daily profit and manage your inventory in one place. No credit card required.
                                        </p>
                                    </div>
                                </template>
                                <template #content>
                                    <div class="pl-0 ml-0" style="list-style-type: none;">
                                        <div class="flex">
                                            ✅&nbsp; <span>Total sales, cost, profit / day</span>
                                        </div>
                                        <div class="flex">
                                            ✅&nbsp; <span>Inventory stock quantities</span>
                                        </div>
                                        <div class="flex">
                                            ✅&nbsp; <span>Best selling products</span>
                                        </div>
                                        <div class="flex">
                                            ✅&nbsp; <span>Sales, cost, profit, refunds / day trend chart</span>
                                        </div>
                                        <div class="flex">
                                            ❌&nbsp; <span>Per order cost and profit</span>
                                        </div>
                                        <div class="flex">
                                            ❌&nbsp; <span>Per order components quantity & cost</span>
                                        </div>
                                        <div class="flex">
                                            ❌&nbsp; <span>AI powered offer generation system</span>
                                        </div>
                                        <div class="flex">
                                            ❌&nbsp; <span>Refunds details</span>
                                        </div>
                                    </div>                                    
                                </template>
                            </Card>
                            <Card v-if="store.subscription.subscription_plan == 'gold'" class="w-20rem">
                                <template #title>
                                    <div class="flex gap-2 flex-column">
                                        <div class="flex justify-content-between align-items-end">
                                            <Badge size="xlarge" :value="store.subscription.subscription_plan.toUpperCase()"  :style="`background-color: ${store.subscription.subscription_plan.toUpperCase() == 'FREE' ?'silver' : '#E1C05C'};color:black`"/>
                                            <Badge size="small" :value="store.subscription.status"  :style="`background-color: ${store.subscription.status == 'active' ?'green' : 'silver'};color:${store.subscription.status == 'active' ?'white' : 'black'}`"/>
                                        </div>
                                        <p style="font-size:0.9rem;">Renewal date: <span>{{ new Date(store.subscription.end_date).toISOString().split('T')[0] }}</span></p>
                                    </div>
                                </template>
                                <template #content>
                                    <div class="pl-0 ml-0" style="list-style-type: none;">
                                        <div class="pl-0 ml-0" style="list-style-type: none;">
                                            <div class="flex">
                                                ✅&nbsp; <span>Total sales, cost, profit / day</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>Inventory stock quantities</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>Best selling products</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>Sales, cost, profit, refunds / day trend chart</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>Per order cost and profit</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>Per order components quantity & cost</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>AI powered offer generation system</span>
                                            </div>
                                            <div class="flex">
                                                ✅&nbsp; <span>Refunds details</span>
                                            </div>
                                        </div>   
                                    </div>                                    
                                </template>
                                <template #footer v-if="store.subscription.status == 'active'">
                                    <div class="flex justify-content-end">
                                        <Button size="small" class="mt-4" text label="Request cancellation" @click="confirm_subscription_cancellation" />
                                    </div>
                                </template>
                            </Card>
                        </div>
                        <h1 style="font-size:2rem;" class="mt-4" v-if="store.subscription.subscription_plan != 'gold'">Upgrade subscription</h1>
                        <div class="flex gap-3" v-if="store.subscription.subscription_plan != 'gold'">
                            <Card class="w-20rem">
                                <template #title>
                                    <Badge class="ml-auto" value="GOLD" size="xlarge" style="background-color: #E1C05C;color:black" />
                                </template>
                                <template #content>
                                    <div class="flex-column flex gap-3 mt-3">
                                        <p class="text-center" style="font-size:1rem;text-decoration: line-through;line-height: 0px;margin-top:-0.5rem;">$60/month</p>
                                        <p class="text-center" style="font-size:1.5rem;line-height: 1rem;">$50/month</p>
                                        <p class="my-2 p-2">
                                            Turn insight into action. Get in depth insights on each order down to components cost, pinpoint profitable dishes, and use AI to create offers that fill your seats.
                                        </p>
                                        <div class="pl-0 ml-0" style="list-style-type: none;">
                                            <div class="pl-0 ml-0" style="list-style-type: none;">
                                                <div class="flex">
                                                    ✅&nbsp; <span>Total sales, cost, profit / day</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>Inventory stock quantities</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>Best selling products</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>Sales, cost, profit, refunds / day trend chart</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>Per order cost and profit</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>Per order components quantity & cost</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>AI powered offer generation system</span>
                                                </div>
                                                <div class="flex">
                                                    ✅&nbsp; <span>Refunds details</span>
                                                </div>
                                            </div>   
                                        </div>    
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
    <ConfirmDialog></ConfirmDialog>
</template>

<script setup lang="ts">
import {getCurrentInstance,ref} from 'vue'
import {ConfirmDialog, Card,Badge, Button,Dialog, ProgressSpinner} from 'primevue';
import axios from 'axios';
import { useToast } from "primevue/usetoast";
import { globalStore } from '@/stores';
import { useConfirm } from "primevue/useconfirm";


const confirm = useConfirm();
const store = globalStore()


const loading_subscription_request = ref(false);
const toast = useToast();
const { proxy } = getCurrentInstance();
const subscrition_request_message = ref("Requesting subscription...");

const confirm_subscription_cancellation = () => {
    confirm.require({
        message: 'Are you sure you want to cancel the subscription?',
        header: 'Confirmation',
        icon: 'pi pi-exclamation-triangle',
        rejectProps: {
            label: 'Cancel',
            severity: 'secondary',
            outlined: true
        },
        acceptProps: {
            label: 'Confirm'
        },
        accept: () => {
            axios.post(`${import.meta.env.VITE_APP_BACKEND_HOST}/v1/api/subscriptions/request_cancellation`, {}, {
                headers: {
                    Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
                }
            }).then((res) => {
                console.log(res);
                toast.add({ severity: 'success', summary: 'Success', detail: 'Subscription cancellation pending', life: 3000, group:'br' });

                let subscription = store.subscription
                subscription.status = "cancellation_pending"
                store.setSubscription(subscription)
                
            }).catch((err) => {
                console.log(err);
                toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000, group:'br' });
            });
        },
        reject: () => {
            toast.add({ severity: 'error', summary: 'Rejected', detail: 'You have rejected', life: 3000 });
        }
    });
};

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