<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Login</h1>
      <form @submit.prevent="login">
        <input
            v-model="name"
            type="text"
            placeholder="Name"
            required
        />
        <input
            v-model="password"
            type="password"
            placeholder="Password"
            required
        />
        <button type="submit" :disabled="loading">
          {{ loading ? 'Logging in…' : 'Login' }}
        </button>
      </form>
      <p class="redirect">
        Don't have an account?
        <router-link to="/register">Register</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'LoginPage',
  data() {
    return {
      name: '',
      password: '',
      loading: false
    }
  },
  methods: {
    async login() {
      this.loading = true
      try {
        const res = await this.$axios.post(
            `${import.meta.env.VITE_API_URL}/login`,
            { name: this.name, password: this.password }
        )
        const token = res.data.token
        localStorage.setItem('token', token)

        const me = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/me`,
            { headers: { Authorization: `Bearer ${token}` } }
        )
        const user = me.data
        localStorage.setItem('user_id', user.user_id)
        localStorage.setItem('username', user.name)

        toast.success('Logged in successfully')
        this.$router.push({ name: 'home' })
      } catch (error) {
        console.error(error)
        const msg = error.response?.data?.error || error.message || 'Login failed'
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
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 360px;
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
  padding: 10px 12px;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  outline: none;
}

.auth-card input:focus {
  border-color: #3498db;
}

.auth-card button {
  padding: 10px;
  font-size: 1rem;
  background: #3498db;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
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
}

.redirect a:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .auth-card {
    padding: 24px;
  }
  .auth-card h1 {
    font-size: 1.5rem;
  }
  .auth-card input, .auth-card button {
    font-size: 0.95rem;
  }
}
</style>
