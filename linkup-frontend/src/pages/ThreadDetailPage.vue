<template>
  <div class="thread-detail">
    <button class="back-btn" @click="$router.back()">← Back</button>

    <transition name="fade" appear>
      <div v-if="thread" class="thread-card">
        <div class="thread-header">
          <h2>{{ thread.user_name }}</h2>
        </div>
        <span class="thread-date">Posted at: {{ formatDate(thread.created_at) }}</span>
        <p class="thread-content">{{ thread.content }}</p>

        <div v-if="thread.media_url" class="thread-media">
          <img :src="thread.media_url" alt="thread media" />
        </div>

        <div class="like-section">
          <button
              class="like-btn"
              :class="{
              liked: liked,
              pulse: pulse,
              loading: likeLoading
            }"
              @click="toggleLike"
              :disabled="likeLoading"
          >
            {{ liked ? '❤️ Liked' : '🤍 Like' }}
          </button>
          <span class="like-count">{{ likes.length }} like{{ likes.length === 1 ? '' : 's' }}</span>
        </div>

        <hr />

        <h3>Comments ({{ comments.length }})</h3>
      </div>
    </transition>

    <transition-group name="fade" tag="ul" class="comments-list">
      <li v-for="c in comments" :key="c.id" class="comment-card">
        <div class="comment-header">
          <strong>{{ c.user_name }}</strong>
          <small class="comment-date">{{ formatDate(c.created_at) }}</small>
        </div>
        <p class="comment-content">{{ c.content }}</p>
      </li>
    </transition-group>

    <div class="new-comment">
      <textarea
          v-model="newComment"
          placeholder="Write a comment..."
          :disabled="posting"
      ></textarea>
      <button @click="postComment" :disabled="posting || !newComment.trim()">
        {{ posting ? 'Posting...' : 'Post Comment' }}
      </button>
    </div>
  </div>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'ThreadDetailPage',
  props: ['id'],
  data() {
    return {
      thread: null,
      comments: [],
      newComment: '',
      posting: false,
      likes: [],
      liked: false,
      pulse: false,
      likeLoading: false
    }
  },
  async created() {
    if (!localStorage.getItem('token')) {
      toast.error('Please login to view thread')
      this.$router.push('/login')
      return
    }

    try {
      await Promise.all([
        this.fetchThread(),
        this.fetchComments(),
        this.fetchLikes()
      ])
    } catch (error) {
      toast.error('Failed to load thread data')
    }
  },
  methods: {
    formatDate(d) {
      return new Date(d).toLocaleString()
    },

    async fetchThread() {
      try {
        const { data } = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/threads/${this.id}`,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`
              }
            }
        )
        this.thread = data
      } catch (err) {
        console.error('Thread fetch error:', err)
        toast.error('Failed to load thread')
      }
    },

    async fetchComments() {
      try {
        const { data } = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/comments/${this.id}`,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`
              }
            }
        )
        this.comments = Array.isArray(data) ? data : []
      } catch (err) {
        console.error('Comments fetch error:', err)
        this.comments = []
      }
    },

    async fetchLikes() {
      try {
        const { data } = await this.$axios.get(
            `${import.meta.env.VITE_API_URL}/user/likes/${this.id}`,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
              },
            }
        );
        console.log('Likes fetched:', data);

        this.likes = Array.isArray(data) ? data : [];

        const userId = localStorage.getItem('user_id');

        if (!userId) {
          console.warn('User ID is not available in localStorage.');
          this.liked = false;
          return;
        }
        this.liked = this.likes.some((like) => {
          const normalizedUserId = like.user_id || like.UserID;
          if (!normalizedUserId) {
            console.warn('Invalid like object, missing user_id/UserID:', like);
            return false;
          }
          return normalizedUserId.toString() === userId.toString();
        });
      } catch (err) {
        console.error('Likes fetch error:', err);
        this.likes = [];
        this.liked = false;
      }
    },

    async toggleLike() {
      if (this.likeLoading) return

      this.likeLoading = true
      try {
        const token = localStorage.getItem('token')
        const userId = localStorage.getItem('user_id')

        if (!userId) {
          toast.error('Please login to like threads')
          return
        }

        if (this.liked) {
          await this.$axios.delete(
              `${import.meta.env.VITE_API_URL}/user/likes/${this.id}`,
              { headers: { Authorization: `Bearer ${token}` } }
          )
        } else {
          await this.$axios.post(
              `${import.meta.env.VITE_API_URL}/user/likes/${this.id}`,
              {},
              { headers: { Authorization: `Bearer ${token}` } }
          )
          this.pulse = true
          setTimeout(() => (this.pulse = false), 500)
        }

        await this.fetchLikes()
        toast.success(this.liked ? 'Like removed' : 'Liked!')
      } catch (err) {
        console.error('Like error:', err)
        toast.error(err.response?.data?.error || 'Like action failed')
      } finally {
        this.likeLoading = false
      }
    },

    async postComment() {
      if (!this.newComment.trim()) return
      this.posting = true
      try {
        await this.$axios.post(
            `${import.meta.env.VITE_API_URL}/user/comments/${this.id}`,
            { content: this.newComment },
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`
              }
            }
        )
        this.newComment = ''
        await this.fetchComments()
        toast.success('Comment posted')
      } catch (err) {
        console.error('Comment post error:', err)
        toast.error(err.response?.data?.error || 'Failed to post comment')
      } finally {
        this.posting = false
      }
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700&display=swap');

.thread-detail {
  font-family: 'Nunito', sans-serif;
  max-width: 700px;
  margin: auto;
  padding: 20px;
  background: #f0f2f5;
}

.back-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  color: #3498db;
  cursor: pointer;
  margin-bottom: 16px;
  padding: 8px 16px;
  border-radius: 4px;
  transition: background 0.2s;
}

.back-btn:hover {
  background: rgba(52, 152, 219, 0.1);
}

.thread-card {
  background: #fff;
  border-radius: 10px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
}

.thread-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.thread-header h2 {
  margin: 0;
  font-weight: 700;
  font-size: 1.5rem;
  color: #2c3e50;
}

.thread-date {
  color: #7f8c8d;
  font-size: 0.9rem;
}

.thread-content {
  color: #34495e;
  line-height: 1.6;
  font-size: 1.1rem;
  margin-bottom: 16px;
}

.thread-media {
  width: 100%;
  overflow: hidden;
  border-radius: 8px;
  margin: 16px 0;
}

.thread-media img {
  width: 100%;
  height: auto;
  border-radius: 8px;
  border: 1px solid #ecf0f1;
}

.like-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 20px 0;
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
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

.like-count {
  color: #7f8c8d;
  font-size: 0.95rem;
}

.comments-list {
  list-style: none;
  padding: 0;
  margin: 20px 0;
}

.comment-card {
  background: white;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-header strong {
  color: #2c3e50;
}

.comment-date {
  color: #95a5a6;
  font-size: 0.8rem;
}

.comment-content {
  color: #34495e;
  line-height: 1.5;
}

.new-comment textarea {
  width: 100%;
  min-height: 100px;
  padding: 12px;
  border: 1px solid #bdc3c7;
  border-radius: 8px;
  margin-bottom: 12px;
  font-family: 'Nunito', sans-serif;
  resize: vertical;
}

.new-comment button {
  background: #3498db;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.new-comment button:hover:not(:disabled) {
  background: #2980b9;
}

.new-comment button:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 600px) {
  .thread-detail {
    padding: 16px;
  }

  .thread-card {
    padding: 16px;
  }

  .thread-header h2 {
    font-size: 1.3rem;
  }

  .thread-content {
    font-size: 1rem;
  }
}
</style>