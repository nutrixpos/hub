import { defineStore } from 'pinia'

export const globalStore = defineStore('global', {
  state: () => ({ 
    count: 0,
    orientation:'ltr',
    subscription: {
      subscription_plan: 'free'
    }
  }),
  getters: {
    double: state => state.count * 2,
    currentOrientation(state) {
      return state.orientation;
    }
  },
  actions: {
    increment() {
      this.count++
    },
    setSubscription(subscription: string) {
      this.subscription = subscription;
    },
    setOrientation(orientation:string){
        this.orientation = orientation;
    }
  },
})
