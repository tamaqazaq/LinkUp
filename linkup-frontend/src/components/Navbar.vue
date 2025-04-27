<template>
  <header>
    <nav class="navbar">
      <div class="nav-left">
        <button class="menu-btn" @click="menuOpen = !menuOpen">
          ☰
        </button>
        <transition name="menu-fade">
          <ul
              v-if="menuOpen"
              class="menu-list"
              @click.self="menuOpen = false"
          >
            <li><router-link :to="{ name: 'home' }">Home</router-link></li>
            <li><router-link :to="{ name: 'create-thread' }">New Thread</router-link></li>
            <li><router-link :to="{ name: 'search' }">Search</router-link></li>
            <li v-if="!user"><router-link :to="{ name: 'login' }">Login</router-link></li>
            <li v-if="!user"><router-link :to="{ name: 'register' }">Register</router-link></li>
            <li v-if="user"><button @click="logout">Logout</button></li>
          </ul>
        </transition>
      </div>

      <div class="nav-right">
        <img
            v-if="user"
            :src="avatarSrc"
            alt="avatar"
            class="avatar"
            @error="onAvatarError"
            @click="goProfile"
        />
      </div>
    </nav>
  </header>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      user: null,
      menuOpen: false
    }
  },
  computed: {
    avatarSrc() {
      if (this.user?.avatar_url) {
        return this.user.avatar_url
      }
      const name = encodeURIComponent(this.user?.name || 'User')
      return `https://ui-avatars.com/api/?name=${name}&background=random&color=fff&size=128`
    }
  },
  async created() {
    await this.loadUser()
  },
  watch: {
    $route() {
      this.menuOpen = false
      this.loadUser()
    }
  },
  methods: {
    async loadUser() {
      const token = localStorage.getItem('token')
      if (!token) {
        this.user = null
        return
      }
      try {
        const res = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/me`,
            {headers: {Authorization: `Bearer ${token}`}}
        )
        this.user = res.data
      } catch {
        this.user = null
      }
    },
    logout() {
      localStorage.removeItem('token')
      this.user = null
      this.$router.push({name: 'login'})
    },
    goProfile() {
      if (this.user) {
        this.$router.push({name: 'my-profile'})
      }
    },
    onAvatarError(e) {
      const initials = encodeURIComponent(this.user?.name || 'U')
      e.target.src = `https://ui-avatars.com/api/?name=${initials}&background=888888&color=fff&size=128`
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600;700&display=swap');

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: 'Poppins', sans-serif;
  background: #ffffff;
  padding: 14px 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.nav-left {
  position: relative;
}

.menu-btn {
  font-size: 1.8rem;
  background: none;
  border: none;
  cursor: pointer;
  color: #333;
}

.menu-list {
  position: absolute;
  top: 110%;
  left: 0;
  background: #ffffff;
  border: 1px solid #ddd;
  border-radius: 10px;
  margin-top: 8px;
  list-style: none;
  padding: 8px 0;
  min-width: 180px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
  z-index: 100;
}

.menu-list li {
  border-bottom: 1px solid #f0f0f0;
}

.menu-list li:last-child {
  border-bottom: none;
}

.menu-list a,
.menu-list button {
  display: block;
  padding: 12px 20px;
  text-decoration: none;
  color: #333;
  font-size: 1.1rem;
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.2s;
}

.menu-list a:hover,
.menu-list button:hover {
  background: #f5f7fa;
}

.nav-right .avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  cursor: pointer;
  border: 2px solid #dfe6e9;
  transition: transform 0.3s;
}

.nav-right .avatar:hover {
  transform: scale(1.05);
}

.menu-fade-enter-active,
.menu-fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.menu-fade-enter-from,
.menu-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.menu-fade-enter-to,
.menu-fade-leave-from {
  opacity: 1;
  transform: translateY(0);
}

/* Адаптив */
@media (max-width: 600px) {
  .menu-list {
    min-width: 150px;
  }

  .menu-btn {
    font-size: 1.6rem;
  }

  .menu-list a,
  .menu-list button {
    font-size: 1rem;
    padding: 10px 14px;
  }

  .nav-right .avatar {
    width: 42px;
    height: 42px;
  }
}
</style>
