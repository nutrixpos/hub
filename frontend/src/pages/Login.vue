<template>
  <div class="grid bg-gray-100 ">
    <div class="col-12 lg:col-6 lg:col-offset-3 min-h-screen flex justify-content-center align-items-center bg-gray-100">
      <Card class="h-auto lg:w-30rem">
        <template #title>
          <div class="text-center">
            {{ isLogin ? 'Login' : 'Register' }}
          </div>
        </template>
        <template #content>
          <form @submit.prevent="handleSubmit">
            <!-- Email Field -->
            <div class="field mb-4">
              <label for="email" class="block text-900 font-medium mb-2">Email</label>
              <InputText
                v-model="form.email"
                type="email"
                id="email"
                class="w-full"
                required
              />
            </div>
            
            <!-- Password Field -->
            <div class="field mb-6">
              <label for="password" class="block text-900 font-medium mb-2">Password</label>
              <Password
                v-model="form.password"
                id="password"
                :feedback="!isLogin"
                toggleMask
                class="w-full"
                :required="true"
                :minLength="isLogin ? 6 : 8"
              />
            </div>
            
            <!-- Name Field (only for register) -->
            <div v-if="!isLogin" class="field mb-4">
              <label for="name" class="block text-900 font-medium mb-2">Full Name</label>
              <InputText
                v-model="form.name"
                type="text"
                id="name"
                class="w-full"
                required
              />
            </div>
            
            <!-- Submit Button -->
            <Button
              type="submit"
              :label="isLogin ? 'Login' : 'Register'"
              class="w-full"
            />
            
            <!-- Toggle between Login/Register -->
            <div class="text-center mt-4">
              <span class="text-600">
                {{ isLogin ? "Don't have an account?" : "Already have an account?" }}
              </span>
              <Button
                type="button"
                @click="toggleAuthMode"
                :label="isLogin ? 'Register' : 'Login'"
                link
                class="p-0 ml-2"
              />
            </div>
          </form>
        </template>
      </Card>
    </div>
  </div>
</template>

<script>
import Card from 'primevue/card';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';

export default {
  components: {
    Card,
    InputText,
    Password,
    Button
  },
  data() {
    return {
      isLogin: true,
      form: {
        email: '',
        password: '',
        name: ''
      }
    }
  },
  methods: {
    toggleAuthMode() {
      this.isLogin = !this.isLogin
      this.form.password = ''
    },
    handleSubmit() {
      if (this.isLogin) {
        this.login()
      } else {
        this.register()
      }
    },
    login() {
      console.log('Logging in with:', this.form.email)
      this.$toast.add({
        severity: 'success',
        summary: 'Login Attempt',
        detail: `Login attempt for ${this.form.email}`,
        life: 3000
      });
    },
    register() {
      console.log('Registering:', this.form)
      this.$toast.add({
        severity: 'success',
        summary: 'Registration',
        detail: `Registration for ${this.form.name} (${this.form.email})`,
        life: 3000
      });
      this.isLogin = true
      this.form.password = ''
    }
  }
}
</script>

<style scoped>
/* Add any custom styles here if needed */
</style>