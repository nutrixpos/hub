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
import zitadelAuth from "@/services/zitadelAuth";



// Import PrimeVue CSS
import 'primeicons/primeicons.css';                           // icons

const secure_routes = [
    {
        path: '/', alias: ['/console'],
        meta: {
            authName: zitadelAuth.oidcAuth.authName
        },
        component: ()=>{
            if (zitadelAuth.hasRole("admin") || zitadelAuth.hasRole("cashier") ) {
                return import('@/pages/Console.vue')
            }
            return import('@/pages/NoAccessView.vue')
        },
        children: [
            {
                path: 'sales',
                component: () => import('@/pages/Sales.vue')
            },
            {
                path: 'inventory',
                component: () => import('@/pages/Inventory.vue')
            },
            {
                path: 'koptan',
                component: () => import('@/pages/Koptan.vue')
            },
        ],
    },
    {
        path: '/login',
        meta: {
            authName: zitadelAuth.oidcAuth.authName
        },
        component: ()=>{
            return import('@/pages/Login.vue')
        },
    },
  ]


const insecure_routes = [
    {
        path: '/', alias: ['/console'],
        component: ()=>{
            return import('@/pages/Console.vue')
        },
        children: [
            {
                path: 'sales',
                component: () => import('@/pages/Sales.vue')
            },
            {
                path: 'inventory',
                component: () => import('@/pages/Inventory.vue')
            },
            {
                path: 'koptan',
                component: () => import('@/pages/Koptan.vue')
            },
        ],
    },
    {
        path: '/login',
        component: ()=>{
            return import('@/pages/Login.vue')
        },
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

  const secureRouter = createRouter({
    history: createWebHistory(),
    routes: secure_routes,
  })
  
  const insecureRouter = createRouter({
    history: createWebHistory(),
    routes: insecure_routes,
  })
  


  if (import.meta.env.VITE_APP_ZITADEL_ENABLED === 'true'){
    zitadelAuth.oidcAuth.useRouter(secureRouter)
  
    zitadelAuth.oidcAuth.startup().then(ok => {
      if (ok) {
            const app = createApp(App).use(createPinia())
            app.config.globalProperties.$zitadel = zitadelAuth


            app
            .use(secureRouter)
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
            .mount('#app')
      } else {
          console.error('Zitadel startup was not ok')
      }
    })
  } else {
    const app = createApp(App).use(createPinia())
  
    app
    .use(insecureRouter)
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
    .mount('#app')
  }
