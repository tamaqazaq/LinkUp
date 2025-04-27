<template>
  <div class="comment-page">
    <button class="back-btn" @click="$router.back()">← Back</button>

    <transition-group name="fade" tag="ul" class="comments-list">
      <li
          v-for="c in comments"
          :key="c.id"
          class="comment-card"
      >
        <div class="comment-header">
          <strong>{{ c.user_name }}</strong>
          <div class="header-right">
            <small class="comment-date">{{ formatDate(c.created_at) }}</small>
            <button
                v-if="c.user_id === meId"
                class="delete-comment"
                @click="deleteComment(c.id)"
            >Delete</button>
          </div>
        </div>
        <p class="comment-content">{{ c.content }}</p>
      </li>
    </transition-group>

    <div class="new-comment">
      <textarea
          v-model="newComment"
          placeholder="Write a comment..."
      ></textarea>
      <button @click="postComment" :disabled="posting">
        {{ posting ? 'Posting...' : 'Post Comment' }}
      </button>
    </div>
  </div>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'CommentPage',
  props: ['threadId'],
  data() {
    return {
      comments: [],
      newComment: '',
      posting: false,
      meId: null
    }
  },
  async created() {
    await this.fetchMeId()
    await this.fetchComments()
  },
  methods: {
    formatDate(dateStr) {
      return new Date(dateStr).toLocaleString()
    },
    async fetchMeId() {
      try {
        const token = localStorage.getItem('token')
        const { data } = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/me`,
            { headers: { Authorization: `Bearer ${token}` } }
        )
        this.meId = data.user_id
      } catch (err) {
        console.error('fetchMeId error', err)
        toast.error('Failed to get your user data')
      }
    },

    async fetchComments() {
      try {
        const token = localStorage.getItem('token')
        const res = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/comments/${this.threadId}`,
            { headers: { Authorization: `Bearer ${token}` } }
        )
        this.comments = res.data
      } catch (err) {
        console.error('Fetch comments error', err)
        toast.error('Failed to load comments')
      }
    },

    async postComment() {
      if (!this.newComment.trim()) return
      this.posting = true
      try {
        const token = localStorage.getItem('token')
        await this.$axios.post(
            `${import.meta.env.VITE_API_URL}/user/comments/${this.threadId}`,
            { content: this.newComment },
            { headers: { Authorization: `Bearer ${token}` } }
        )
        this.newComment = ''
        await this.fetchComments()
        toast.success('Comment posted')
      } catch (err) {
        console.error('Post comment error', err)
        toast.error('Failed to post comment')
      } finally {
        this.posting = false
      }
    },

    async deleteComment(commentId) {
      if (!confirm('Delete this comment?')) return
      try {
        const token = localStorage.getItem('token')
        await this.$axios.delete(
            `${import.meta.env.VITE_API_URL}/user/comments/${commentId}`,
            { headers: { Authorization: `Bearer ${token}` } }
        )
        // Удаляем из массива
        this.comments = this.comments.filter(c => c.id !== commentId)
        toast.success('Comment deleted')
      } catch (err) {
        console.error('Delete comment error', err)
        toast.error('Failed to delete comment')
      }
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700&display=swap');

.comment-page {
  font-family: 'Nunito', sans-serif;
  max-width: 600px;
  margin: auto;
  padding: 20px;
  background: #f0f2f5;
}

.back-btn {
  background: none;
  border: none;
  font-size: 1rem;
  color: #3498db;
  cursor: pointer;
  margin-bottom: 16px;
}

.comments-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.comment-card {
  background: #fff;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.comment-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.header-right {
  display: flex;
  align-items: center;
}

.comment-date {
  color: #999;
  font-size: 0.8rem;
  margin-right: 8px;
}

.delete-comment {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 0.8rem;
}

.comment-content {
  color: #333;
  line-height: 1.5;
}

.new-comment {
  margin-top: 16px;
}

.new-comment textarea {
  width: 100%;
  min-height: 80px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  resize: vertical;
  font-family: inherit;
  font-size: 1rem;
}

.new-comment button {
  margin-top: 8px;
  padding: 8px 16px;
  background: #3498db;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
}

.new-comment button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
  .comment-page {
    padding: 16px;
  }

  .comment-card {
    padding: 12px;
  }

  .new-comment textarea {
    font-size: 0.9rem;
  }

  .new-comment button {
    width: 100%;
    font-size: 0.95rem;
  }
}
</style>
