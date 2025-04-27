<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Register</h1>
      <form @submit.prevent="register">
        <input
            v-model="name"
            type="text"
            placeholder="Name"
            required
        />
        <input
            v-model="email"
            type="email"
            placeholder="Email"
            required
        />
        <input
            v-model="password"
            type="password"
            placeholder="Password"
            required
        />
        <button type="submit" :disabled="loading">
          {{ loading ? 'Registering…' : 'Register' }}
        </button>
      </form>
      <p class="redirect">
        Already have an account?
        <router-link :to="{ name: 'login' }">Login</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'RegisterPage',
  data() {
    return {
      name: '',
      email: '',
      password: '',
      loading: false
    }
  },
  methods: {
    async register() {
      this.loading = true
      try {
        await this.$axios.post(
            `${import.meta.env.VITE_API_URL}/register`,
            {
              name: this.name,
              email: this.email,
              password: this.password
            }
        )
        toast.success('Registration successful! Please check your email to verify your account.')
      } catch (err) {
        console.error('Register error', err)
        const msg = err.response?.data?.error
            || err.message
            || 'Registration failed'
        toast.error(msg)
      } finally {
        this.loading = false
      }
    }

  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700&display=swap');

/* box-sizing */
*, *::before, *::after {
  box-sizing: border-box;
}

.auth-page {
  font-family: 'Nunito', sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
  padding: 16px;
}

.auth-card {
  background: #fff;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-width: 360px;
  width: 100%;
  transition: transform 0.2s, box-shadow 0.2s;
}
.auth-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.auth-card h1 {
  margin-bottom: 24px;
  font-size: 1.75rem;
  color: #333;
  text-align: center;
}

.auth-card form {
  display: flex;
  flex-direction: column;
}

.auth-card input {
  margin-bottom: 16px;
  padding: 12px;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  outline: none;
  transition: border-color 0.2s;
}
.auth-card input:focus {
  border-color: #3498db;
}

.auth-card button {
  padding: 12px;
  font-size: 1rem;
  background: #3498db;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s, transform 0.1s;
}
.auth-card button:hover:not(:disabled) {
  background: #2980b9;
}
.auth-card button:active:not(:disabled) {
  transform: translateY(1px);
}
.auth-card button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.redirect {
  margin-top: 16px;
  text-align: center;
  color: #555;
  font-size: 0.9rem;
}
.redirect a {
  color: #3498db;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s;
}
.redirect a:hover {
  color: #2980b9;
  text-decoration: underline;
}

/* mobile */
@media (max-width: 480px) {
  .auth-card {
    padding: 24px;
  }
  .auth-card h1 {
    font-size: 1.5rem;
  }
  .auth-card input,
  .auth-card button {
    font-size: 0.95rem;
  }
}
</style>
