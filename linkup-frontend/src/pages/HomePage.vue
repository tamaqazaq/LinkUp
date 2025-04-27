<template>
  <div class="home-page">
    <h1>All Threads</h1>

    <transition-group name="fade-up" tag="div" class="threads-container">
      <div
          v-for="thr in threads"
          :key="thr.id"
          class="thread"
          @click="goToThread(thr.id)"
      >
        <div class="thread-header">
          <img
              :src="avatarUrl(thr)"
              class="avatar"
              @click.stop="goToProfile(thr.user_id)"
          />
          <div class="user-info">
            <router-link
                :to="linkFor(thr.user_id)"
                class="username"
                @click.stop
            >
              {{ thr.user_name }}
            </router-link>
            <span class="thread-date">{{ formatDate(thr.created_at) }}</span>
          </div>
        </div>

        <div class="thread-body">
          <p class="thread-content">
            {{ thr.content }}
          </p>

          <div v-if="thr.media_url" class="thread-media">
            <img :src="thr.media_url" alt="thread media" />
          </div>
        </div>

        <div class="thread-footer">
          <div class="comments-section">
            <span class="comments-link">
              💬 {{ thr.comment_count }} {{ thr.comment_count === 1 ? 'Comment' : 'Comments' }}
            </span>
          </div>
          <div class="likes-section">
            <button
                class="like-btn"
                :class="{ liked: thr.liked, pulse: thr.pulse }"
                @click.stop="toggleLike(thr)"
                :disabled="thr.likeLoading"
            >
              {{ thr.liked ? '❤️ Liked' : '🤍 Like' }}
            </button>
            <span class="like-count">
              {{ thr.likes.length }} like{{ thr.likes.length === 1 ? '' : 's' }}
            </span>
          </div>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'HomePage',
  data() {
    return {
      threads: [],
      meId: null
    }
  },
  async created() {
    await this.fetchMe()
    await this.fetchThreads()
    await this.fetchCommentCounts()
    await this.fetchLikes()
  },
  methods: {
    formatDate(d) {
      return new Date(d).toLocaleString()
    },
    avatarUrl(thr) {
      if (thr.avatar_url) {
        return thr.avatar_url
      } else {
        return `https://ui-avatars.com/api/?name=${encodeURIComponent(thr.user_name)}&background=random`
      }
    },
    async fetchMe() {
      const token = localStorage.getItem('token')
      if (!token) return
      try {
        const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/me`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.meId = data.user_id
      } catch (error) {
        this.handleUnauthorized(error)
      }
    },
    async fetchThreads() {
      const token = localStorage.getItem('token')
      if (!token) return
      try {
        const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/threads`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.threads = Array.isArray(data)
            ? data.map(t => ({
              ...t,
              comment_count: 0,
              likes: [],
              liked: false,
              likeLoading: false,
              pulse: false
            }))
            : []
      } catch (err) {
        console.error('Fetch threads error', err)
        this.threads = []
      }
    },
    async fetchCommentCounts() {
      const token = localStorage.getItem('token')
      if (!token) return
      await Promise.all(this.threads.map(async thr => {
        try {
          const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/comments/${thr.id}`, {
            headers: { Authorization: `Bearer ${token}` }
          })
          thr.comment_count = Array.isArray(data) ? data.length : 0
        } catch (error) {
          this.handleUnauthorized(error)
        }
      }))
    },
    async fetchLikes() {
      const token = localStorage.getItem('token')
      if (!token) return
      await Promise.all(this.threads.map(async thr => {
        try {
          const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/likes/${thr.id}`, {
            headers: { Authorization: `Bearer ${token}` }
          })
          thr.likes = Array.isArray(data) ? data : []
          thr.liked = thr.likes.some(like => {
            const normalizedUserId = like.user_id || like.UserID // Handle inconsistent property names
            if (!normalizedUserId) {
              console.warn('Invalid like object, missing user_id/UserID:', like)
              return false
            }
            return normalizedUserId.toString() === this.meId.toString()
          })
        } catch (err) {
          console.error('Fetch likes error:', err)
          thr.likes = []
          thr.liked = false
        }
      }))
    },
    async toggleLike(thread) {
      if (thread.likeLoading) return
      thread.likeLoading = true
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          toast.error('Please login to like threads')
          return
        }

        if (thread.liked) {
          await this.$axios.delete(`${import.meta.env.VITE_API_URL}/user/likes/${thread.id}`, {
            headers: { Authorization: `Bearer ${token}` }
          })
        } else {
          await this.$axios.post(`${import.meta.env.VITE_API_URL}/user/likes/${thread.id}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          })
          thread.pulse = true // Trigger animation
          setTimeout(() => (thread.pulse = false), 500)
        }

        const { data } = await this.$axios.get(`${import.meta.env.VITE_API_URL}/user/likes/${thread.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        thread.likes = Array.isArray(data) ? data : []
        thread.liked = thread.likes.some(like => {
          const normalizedUserId = like.user_id || like.UserID // Handle inconsistent property names
          return normalizedUserId && normalizedUserId.toString() === this.meId.toString()
        })
        toast.success(thread.liked ? 'Liked!' : 'Like removed')
      } catch (err) {
        console.error('Toggle like error:', err)
        toast.error('Failed to update like status')
      } finally {
        thread.likeLoading = false
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
    linkFor(userId) {
      return userId === this.meId ? { name: 'my-profile' } : { name: 'profile', params: { id: userId } }
    },
    goToProfile(userId) {
      this.$router.push(this.linkFor(userId))
    },
    goToThread(threadId) {
      this.$router.push({ name: 'thread-detail', params: { id: threadId } })
    }
  }
}
</script>

<style scoped>
.home-page {
  font-family: 'Poppins', sans-serif;
  background-color: #f7f9fb;
  padding: 40px 20px;
}

h1 {
  font-weight: 700;
  font-size: 2.4rem;
  margin-bottom: 30px;
  color: #333;
  text-align: center;
}

.threads-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.thread {
  background: #ffffff;
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.thread:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.thread-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.avatar {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #ddd;
  cursor: pointer;
  margin-right: 14px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.username {
  font-weight: 700;
  text-decoration: none;
  color: #2c3e50;
  font-size: 1.1rem;
}

.username:hover {
  text-decoration: underline;
}

.thread-date {
  color: #888;
  font-size: 0.85rem;
}

.thread-body {
  margin-top: 8px;
}

.thread-content {
  font-size: 1.2rem;
  color: #444;
  line-height: 1.7;
  margin-bottom: 14px;
}

.thread-media img {
  width: 100%;
  border-radius: 12px;
  margin-top: 8px;
}

/* Comments and Likes Section */
.thread-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.comments-section {
  font-size: 1rem;
  color: #3498db;
  font-weight: 600;
  cursor: pointer;
}

.comments-section:hover {
  text-decoration: underline;
}

.likes-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.like-btn {
  padding: 8px 20px;
  background: #f8f9fa;
  border: 2px solid #3498db;
  border-radius: 20px;
  color: #3498db;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.like-btn:hover:not(:disabled) {
  transform: scale(1.05);
}

.like-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.like-btn.liked {
  background: #3498db;
  color: white;
}

.like-btn.pulse {
  animation: pulse 0.5s ease;
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
  }
}

.like-count {
  color: #7f8c8d;
  font-size: 0.95rem;
}

/* Responsive Styling */
@media (max-width: 768px) {
  .home-page {
    padding: 30px 16px;
  }

  h1 {
    font-size: 2rem;
  }

  .avatar {
    width: 46px;
    height: 46px;
  }

  .username {
    font-size: 1rem;
  }

  .thread-content {
    font-size: 1.05rem;
  }

  .thread-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .likes-section {
    align-self: flex-end;
  }
}
</style>