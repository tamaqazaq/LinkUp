<template>
  <div class="profile-page">
    <button class="back-btn" @click="$router.back()">← Back</button>

    <div v-if="user" class="profile-container">
      <div class="profile-header-card">
        <div class="avatar-name-wrapper">
          <div class="avatar-wrapper">
            <img :src="avatarSrc" alt="avatar" class="avatar-large" />
          </div>
          <div class="user-info">
            <h2 class="username">{{ user.name }}</h2>
          </div>
        </div>
      </div>

      <div class="details-card">
        <div class="detail-item">
          <span class="detail-label">Bio:</span>
          <p class="detail-content">{{ user.bio || '—' }}</p>
        </div>

        <div class="detail-item">
          <span class="detail-label">Location:</span>
          <p class="detail-content">{{ user.location || '—' }}</p>
        </div>

        <div class="detail-item">
          <span class="detail-label">Social Links:</span>
          <p class="detail-content">
            <a v-if="user.social_links" :href="user.social_links" target="_blank" class="social-link">
              {{ user.social_links }}
            </a>
            <span v-else>—</span>
          </p>
        </div>

        <div class="detail-item">
          <span class="detail-label">Email:</span>
          <p class="detail-content">{{ user.email }}</p>
        </div>

        <div class="detail-item">
          <span class="detail-label">Joined:</span>
          <p class="detail-content">{{ formatDate(user.created_at) }}</p>
        </div>
      </div>

      <h2 class="threads-title">Threads</h2>

      <div v-if="threads.length === 0" class="no-threads">
        <p>No threads yet.</p>
      </div>

      <transition-group name="fade" tag="div" class="threads-list" v-else>
        <router-link
            v-for="thr in threads"
            :key="thr.id"
            :to="{ name: 'thread-detail', params: { id: thr.id } }"
            class="thread-card"
        >
          <div class="thread-header">
            <span>{{ formatDate(thr.created_at) }}</span>
            <span>{{ thr.comment_count }} {{ thr.comment_count === 1 ? 'comment' : 'comments' }}</span>
          </div>
          <p class="thread-content">{{ thr.content }}</p>
          <div v-if="thr.media_url" class="thread-media">
            <img :src="thr.media_url" alt="thread media" />
          </div>
        </router-link>
      </transition-group>
    </div>

    <div v-else class="loading">Loading profile…</div>
  </div>
</template>
<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'ProfilePage',
  props: ['id'],
  data() {
    return {
      meId: null,
      user: null,
      threads: [],
      avatarSrc: '',
    }
  },
  async created() {
    await this.fetchMeId()
    try {
      if (this.id === this.meId) {
        await this.fetchMe()
      } else {
        await this.fetchUserById()
      }
      await this.fetchThreads()
      await this.fetchCommentCounts()
    } catch (error) {
      this.handleUnauthorized(error)
    }
  },
  methods: {
    formatDate(d) {
      return new Date(d).toLocaleDateString()
    },
    onAvatarError() {
      if (this.user?.name) {
        const name = encodeURIComponent(this.user.name)
        this.avatarSrc = `https://ui-avatars.com/api/?name=${name}&background=random`
      }
    },
    handleUnauthorized(error) {
      if (error.response?.status === 401) {
        localStorage.removeItem('token')
        toast.error('Session expired. Please login again.')
        this.$router.push({ name: 'login' })
      } else {
        console.error(error)
      }
    },
    async fetchMeId() {
      const token = localStorage.getItem('token')
      const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/me`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.meId = data.user_id
    },
    async fetchMe() {
      const token = localStorage.getItem('token')
      const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/me`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.user = data
      this.avatarSrc = data.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(data.name)}&background=random`
    },
    async fetchUserById() {
      const token = localStorage.getItem('token')
      const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/${this.id}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.user = data
      this.avatarSrc = data.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(data.name)}&background=random`
    },
    async fetchThreads() {
      const token = localStorage.getItem('token')
      const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/users/${this.id}/threads`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.threads = Array.isArray(data) ? data : []
    },
    async fetchCommentCounts() {
      const token = localStorage.getItem('token')
      await Promise.all(this.threads.map(async thr => {
        const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/comments/${thr.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        thr.comment_count = Array.isArray(data) ? data.length : 0
      }))
    }
  }
}
</script>
<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600;700&display=swap');

.profile-page {
  font-family: 'Poppins', sans-serif;
  max-width: 800px;
  margin: auto;
  padding: 24px;
  background: #fafafa;
}

.back-btn {
  background: none;
  border: none;
  color: #3498db;
  font-size: 1.2rem;
  cursor: pointer;
  margin-bottom: 20px;
}

.profile-container {
  display: flex;
  flex-direction: column;
}

.profile-header-card {
  background: #fff;
  padding: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;
}

.avatar-name-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar-wrapper {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.avatar-large {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-info .username {
  font-size: 1.5rem;
  font-weight: bold;
}

.details-card {
  background: #fff;
  padding: 24px;
  border-radius: 12px;
  margin-bottom: 30px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-label {
  font-size: 0.9rem;
  font-weight: 600;
  color: #888;
  margin-bottom: 6px;
}

.detail-content {
  font-size: 1.1rem;
  color: #333;
  word-break: break-word;
}

.social-link {
  color: #3498db;
  text-decoration: none;
  max-width: 100%;
  display: inline-block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.social-link:hover {
  text-decoration: underline;
  color: #2980b9;
}

.threads-title {
  margin: 24px 0 16px;
  font-weight: 700;
  font-size: 1.6rem;
}

.no-threads {
  text-align: center;
  margin-top: 20px;
}

.threads-list {
  display: grid;
  gap: 20px;
}

.thread-card {
  background: #fff;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: transform 0.3s ease;
  cursor: pointer;
  text-decoration: none;
}
.thread-card:hover {
  transform: translateY(-4px);
}

.thread-header {
  display: flex;
  justify-content: space-between;
  color: #777;
  font-size: 0.9rem;
  margin-bottom: 8px;
}

.thread-content {
  font-size: 1.1rem;
  color: #333;
  margin-bottom: 8px;
}

.thread-media img {
  width: 100%;
  border-radius: 8px;
  margin-top: 8px;
}

/* Fade animation */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.fade-enter-to, .fade-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
