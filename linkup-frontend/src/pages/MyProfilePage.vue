<template>
  <div class="my-profile-page">
    <button class="back-btn" @click="$router.back()">← Back</button>

    <div v-if="user" class="profile-container">
      <div class="profile-header-card">
        <div class="avatar-name-wrapper">
          <div class="avatar-wrapper">
            <img :src="userAvatar" alt="avatar" class="avatar-large" />
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
        <button class="edit-profile-btn" @click="goToEdit">Update</button>
      </div>


      <h2 class="threads-title">My Threads</h2>

      <div v-if="threads.length === 0" class="no-threads">
        <p>No threads yet.</p>
        <button class="create-thread-btn" @click="$router.push({ name: 'create-thread' })">
          Create Thread
        </button>
      </div>

      <transition-group name="fade" tag="div" class="threads-list" v-else>
        <div
            v-for="thr in threads"
            :key="thr.id"
            class="thread-card"
            @click="goToThread(thr.id)"
        >
          <div class="thread-header">
            <span>{{ formatDate(thr.created_at) }}</span>
            <span>{{ thr.comment_count }} {{ thr.comment_count === 1 ? 'comment' : 'comments' }}</span>
          </div>
          <p class="thread-content">{{ thr.content }}</p>
          <div v-if="thr.media_url" class="thread-media">
            <img :src="thr.media_url" alt="thread media" />
          </div>

          <div class="thread-actions">
            <button class="edit-btn" @click.stop="startEdit(thr)">Edit</button>
            <button class="delete-btn" @click.stop="deleteThread(thr)">Delete</button>
          </div>

          <div v-if="editingId === thr.id" class="edit-area">
            <textarea v-model="editingContent" class="edit-textarea"></textarea>
            <button class="save-btn" @click="saveEdit(thr)" :disabled="saving">
              {{ saving ? 'Saving...' : 'Save' }}
            </button>
            <button class="cancel-btn" @click="cancelEdit()">Cancel</button>
          </div>
        </div>
      </transition-group>
    </div>

    <div v-else class="loading">Loading profile…</div>
  </div>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'MyProfilePage',
  data() {
    return {
      user: null,
      threads: [],
      editingId: null,
      editingContent: '',
      saving: false
    }
  },
  computed: {
    userAvatar() {
      if (!this.user) return ''
      if (this.user.avatar_url) {
        return this.user.avatar_url
      }
      const name = encodeURIComponent(this.user.name || 'User')
      return `https://ui-avatars.com/api/?name=${name}&background=random`
    }
  },
  async created() {
    await this.fetchMe()
    if (this.user && this.user.user_id) {
      await this.fetchMyThreads()
      await this.fetchCommentCounts()
    }
  },
  methods: {
    formatDate(d) {
      return new Date(d).toLocaleDateString()
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
    async fetchMe() {
      const token = localStorage.getItem('token')
      try {
        const { data } = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/me`,
            { headers: { Authorization: `Bearer ${token}` } }
        )
        this.user = data
      } catch (error) {
        this.handleUnauthorized(error)
      }
    },
    async fetchMyThreads() {
      const token = localStorage.getItem('token')
      try {
        const {data} = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/users/${this.user.user_id}/threads`,
            {headers: {Authorization: `Bearer ${token}`}}
        )
        this.threads = Array.isArray(data) ? data : []
      } catch (error) {
        this.handleUnauthorized(error)
      }
    },
    async fetchCommentCounts() {
      const token = localStorage.getItem('token')
      await Promise.all(this.threads.map(async thr => {
        try {
          const {data} = await this.$axios.get(
              `${import.meta.env.VITE_API_URL}/user/comments/${thr.id}`,
              {headers: {Authorization: `Bearer ${token}`}}
          )
          thr.comment_count = Array.isArray(data) ? data.length : 0
        } catch (error) {
          this.handleUnauthorized(error)
        }
      }))
    },
    goToEdit() {
      this.$router.push({name: 'edit-profile'})
    },
    startEdit(thr) {
      this.editingId = thr.id
      this.editingContent = thr.content
    },
    cancelEdit() {
      this.editingId = null
      this.editingContent = ''
    },
    async saveEdit(thr) {
      this.saving = true
      const token = localStorage.getItem('token')
      try {
        await this.$axios.put(
            `${import.meta.env.VITE_API_URL}/user/threads/${thr.id}`,
            {content: this.editingContent},
            {headers: {Authorization: `Bearer ${token}`}}
        )
        await this.fetchMyThreads()
        await this.fetchCommentCounts()
        this.cancelEdit()
        toast.success('Thread updated')
      } catch (error) {
        this.handleUnauthorized(error)
      } finally {
        this.saving = false
      }
    },
    async deleteThread(thr) {
      if (!confirm('Delete this thread?')) return
      const token = localStorage.getItem('token')
      try {
        await this.$axios.delete(
            `${import.meta.env.VITE_API_URL}/user/threads/${thr.id}`,
            {headers: {Authorization: `Bearer ${token}`}}
        )
        this.threads = this.threads.filter(t => t.id !== thr.id)
        toast.success('Thread deleted')
      } catch (error) {
        this.handleUnauthorized(error)
      }
    },
    goToThread(id) {
      this.$router.push({name: 'thread-detail', params: {id}})
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600;700&display=swap');

.my-profile-page {
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

.user-info .joined-date {
  font-size: 0.95rem;
  color: #777;
}

.edit-profile-btn {
  background: #3498db;
  color: white;
  padding: 12px 24px; /* Больше внутренний отступ */
  border-radius: 999px; /* Полная круглая кнопка */
  font-weight: bold;
  font-size: 1.1rem; /* Увеличенный размер шрифта */
  border: none;
  cursor: pointer;
  transition: background 0.3s, transform 0.2s;
}

.edit-profile-btn:hover {
  background: #2980b9;
  transform: translateY(-2px);
}


.edit-profile-btn:hover {
  background: #2980b9;
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

.threads-title {
  margin: 24px 0 16px;
  font-weight: 700;
  font-size: 1.6rem;
}

.no-threads {
  text-align: center;
  margin-top: 20px;
}

.create-thread-btn {
  background: #3498db;
  color: white;
  padding: 10px 18px;
  border-radius: 8px;
  font-weight: bold;
  border: none;
  margin-top: 10px;
  cursor: pointer;
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
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.thread-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
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

.thread-actions {
  display: flex;
  gap: 12px;
  margin-top: 10px;
}

.edit-btn, .delete-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 999px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 1.2rem;
}

.edit-btn {
  background: #3498db;
  color: white;
}

.delete-btn {
  background: #e74c3c;
  color: white;
}

.edit-btn:hover {
  background: #2980b9;
}

.delete-btn:hover {
  background: #c0392b;
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
@media (max-width: 600px) {
  .details-card {
    padding: 20px;
    gap: 16px;
  }

  .detail-label {
    font-size: 0.85rem;
  }

  .detail-content {
    font-size: 1rem;
  }
}
.edit-area {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.edit-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-family: 'Poppins', sans-serif;
  font-size: 1rem;
  resize: vertical;
  min-height: 100px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.edit-textarea:focus {
  border-color: #3498db;
  box-shadow: 0 0 8px rgba(52, 152, 219, 0.3);
  outline: none;
}

.save-btn, .cancel-btn {
  padding: 10px 20px;
  font-size: 1rem;
  font-weight: bold;
  border: none;
  border-radius: 999px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.save-btn {
  background: #2ecc71;
  color: #fff;
}

.cancel-btn {
  background: #e74c3c;
  color: #fff;
}

.save-btn:hover {
  background: #27ae60;
}

.cancel-btn:hover {
  background: #c0392b;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Адаптив для редактирования */
@media (max-width: 600px) {
  .edit-textarea {
    font-size: 0.95rem;
    padding: 10px;
  }

  .save-btn, .cancel-btn {
    font-size: 0.95rem;
    padding: 10px 18px;
  }
}

</style>
