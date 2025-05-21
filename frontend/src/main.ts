import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import ToastService from 'primevue/toastservice';
import { createPinia } from 'pinia'
import {  createWebHistory, createRouter } from 'vue-router'
import { createI18n } from 'vue-i18n'
import { definePreset } from '@primeuix/themes';


// Import PrimeVue CSS
import 'primeicons/primeicons.css';                           // icons


const routes = [
    {
        path: '/console',
        component: ()=>{
            return import('@/pages/Console.vue')
        },
        children: [
            {
                path: 'sales',
                component: () => import('@/pages/Sales.vue')
            },
        ],
    },
    {
        path: '/login',
        component: ()=>{
            return import('@/pages/Login.vue')
        },
    },
    { 
        path: '/', alias:['/home'], 
        component: () => {
            return import('@/pages/Login.vue')
        }
    },
  ]


  const i18n = createI18n({
    legacy: false,
    locale: 'en',
    fallbackLocale: 'en',
    messages: {
      en: {
        "cashier":"Cashier",
        "kitchen":"Kitchen",
        "admin":"Admin",
        "inventory":"Inventory",
        "product": "Product | Products",
        "order":"Order | Orders",
        "order_items":"Order Items",
        "total":"Total",
        "subtotal":"Subtotal",
        "discount":"Discount",
        "egp":"EGP",
        "search":"Search",
        "signout":"Signout",
        "notifications":"Notifications",
        "clear_all":"Clear All",
        "stashed_orders":"Stashed Orders",
        "chats":"Chats",
        "messages":"Messages",
        "write_message":"Write Message",
        "paylater_orders":"Paylater Orders",
        "checkout":"Checkout",
        "category":"Category | Categories",
        "add_component":"Add Component",
        "name":"Name",
        "quantity":"Quantity",
        "unit":"Unit",
        "status":"Status",
        "actions":"Actions",
        "history":"History",
        "list":"List",
        "report":"Report | Reports",
        "settings":"Settings",
        "language":"Language | Languages",
        "sales":"Sales"
      }
    }
  })

  const preset = definePreset(Aura, {
    semantic: {
        primary: {
            50: '{zinc.50}',
            100: '{zinc.100}',
            200: '{zinc.200}',
            300: '{zinc.300}',
            400: '{zinc.400}',
            500: '{zinc.500}',
            600: '{zinc.600}',
            700: '{zinc.700}',
            800: '{zinc.800}',
            900: '{zinc.900}',
            950: '{zinc.950}'
        },
        colorScheme: {
            light: {
                primary: {
                    color: '#2e4762',
                    inverseColor: '#FFDC00',
                    hoverColor: '#365473',
                    activeColor: '#263a51'
                },
                highlight: {
                    background: '#fff6c7',
                    focusBackground: '#FFDC00',
                    color: '#173350',
                    focusColor: '#173350'
                }
            },
            dark: {
                primary: {
                    color: '#FFDC00',
                    inverseColor: '#2e4762',
                    hoverColor: '#ffec54',
                    activeColor: '#ffce1e'
                },
                highlight: {
                    background: '#fff6c7',
                    focusBackground: '#FFDC00',
                    color: '#173350',
                    focusColor: '#173350'
                }
            },
        }
    }
  });


const router = createRouter({
    history: createWebHistory(),
    routes: routes,
  })

createApp(App)
.use(router)
.use(PrimeVue, {
    theme: {
        preset: preset,
        options: {
            darkModeSelector: '.my-app-dark',
        }
    }
 })
 .use(ToastService)
 .use(i18n)
 .use(createPinia())
 .mount('#app')
