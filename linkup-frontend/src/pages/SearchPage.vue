<template>
  <div class="search-page">
    <h1>Find Users</h1>

    <div class="search-box">
      <input
          v-model="query"
          type="text"
          placeholder="Search by username..."
      />
      <button @click="search" :disabled="loading">
        {{ loading ? 'Searching...' : 'Search' }}
      </button>
    </div>

    <div v-if="results.length">
      <transition-group name="fade" tag="ul" class="results-list">
        <li
            v-for="user in results"
            :key="user.user_id"
            class="result-item"
            @click="goToProfile(user.user_id)"
        >
          <img
              :src="user.avatar_url || defaultAvatar(user.name)"
              :alt="`Avatar of ${user.name}`"
              class="avatar-small"
          />
          <div class="user-info">
            <span class="username">{{ user.name }}</span>
          </div>
        </li>
      </transition-group>
    </div>

    <p v-else-if="!loading && queried" class="no-results">
      No users found.
    </p>
  </div>
</template>

<script>
let debounceTimeout = null

export default {
  name: 'SearchPage',
  data() {
    return {
      query: '',
      results: [],
      loading: false,
      queried: false
    }
  },
  watch: {
    query(newQuery) {
      clearTimeout(debounceTimeout)
      if (!newQuery.trim()) {
        this.results = []
        this.queried = false
        return
      }
      debounceTimeout = setTimeout(() => {
        this.search()
      }, 400)
    }
  },
  methods: {
    async search() {
      this.loading = true
      this.queried = true
      try {
        const { data } = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/search`,
            { params: { query: this.query } }
        )
        this.results = data.users
      } catch (err) {
        console.error('Search error', err)
        alert('Search failed')
      } finally {
        this.loading = false
      }
    },
    goToProfile(userId) {
      const meId = this.$root.$data?.meId || localStorage.getItem('user_id')
      if (userId === meId) {
        this.$router.push({ name: 'my-profile' })
      } else {
        this.$router.push({ name: 'profile', params: { id: userId } })
      }
    },
    defaultAvatar(name) {
      return `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=random`
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700&display=swap');

/* Global box-sizing */
*, *::before, *::after {
  box-sizing: border-box;
}

.search-page {
  font-family: 'Nunito', sans-serif;
  max-width: 600px;
  margin: auto;
  padding: 16px;
  background: linear-gradient(to bottom right, #f9f9f9, #e6f0fa);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

h1 {
  margin-bottom: 24px;
  color: #2c3e50;
  text-align: center;
  font-size: 1.8rem;
  font-weight: 700;
}

.search-box {
  display: flex;
  align-items: center;
  background: #fff;
  padding: 10px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.search-box input {
  flex: 1;
  padding: 12px 16px;
  font-size: 1rem;
  border: none;
  outline: none;
  border-radius: 8px;
}

.search-box input:focus {
  background: #f0f8ff;
}

.search-box button {
  margin-left: 12px;
  padding: 12px 20px;
  background: #3498db;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s ease;
}

.search-box button:hover:not(:disabled) {
  background: #1d7ecb;
}

.search-box button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.results-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.result-item {
  display: flex;
  align-items: center;
  background: #fff;
  padding: 16px;
  margin-bottom: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.result-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.avatar-small {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  margin-right: 16px;
  object-fit: cover;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 1.1rem;
  color: #2c3e50;
  font-weight: 600;
}

.username:hover {
  text-decoration: underline;
}

.user-id {
  font-size: 0.85rem;
  color: #7f8c8d;
}

.no-results {
  text-align: center;
  color: #888;
  margin-top: 24px;
  font-size: 1rem;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
  transform: translateY(0);
}

@media (max-width: 480px) {
  .search-page {
    padding: 12px;
  }

  h1 {
    font-size: 1.6rem;
  }

  .result-item {
    padding: 12px;
  }

  .search-box input {
    font-size: 0.9rem;
    padding: 10px 12px;
  }

  .search-box button {
    font-size: 0.9rem;
    padding: 10px 16px;
  }
}
</style>