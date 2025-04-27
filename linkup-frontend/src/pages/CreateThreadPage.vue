<template>
  <transition name="fade" appear>
    <div class="create-thread-page">
      <div class="create-thread-card">
        <button class="back-btn" @click="$router.back()">← Back</button>
        <h1>Create New Thread</h1>
        <form @submit.prevent="submitThread" class="thread-form">
          <div class="form-group">
            <label for="content">Content:</label>
            <textarea
                id="content"
                v-model="content"
                required
                placeholder="What's on your mind?"
            ></textarea>
          </div>
          <div class="form-group">
            <label for="media">Media (optional):</label>
            <input
                id="media"
                type="file"
                @change="onFileChange"
                accept="image/jpeg,image/png,video/mp4"
            />
            <p v-if="fileError" class="error-text">{{ fileError }}</p>
          </div>
          <button type="submit" :disabled="loading">
            {{ loading ? 'Posting...' : 'Post Thread' }}
          </button>
        </form>
      </div>
    </div>
  </transition>
</template>

<script>
import { toast } from 'vue3-toastify'

export default {
  name: 'CreateThreadPage',
  data() {
    return {
      content: '',
      file: null,
      fileError: '',
      loading: false
    }
  },
  methods: {
    onFileChange(e) {
      this.fileError = ''
      const file = e.target.files[0] || null
      if (!file) {
        this.file = null
        return
      }

      const allowedTypes = ['image/jpeg', 'image/png', 'video/mp4']
      if (!allowedTypes.includes(file.type)) {
        this.fileError = 'Only JPG, PNG images and MP4 videos are allowed.'
        toast.error(this.fileError)
        this.file = null
        e.target.value = ''
        return
      }

      const url = URL.createObjectURL(file)
      const allowedRatios = [16/9, 4/3, 3/4, 1]
      const tolerance = 0.02

      const checkRatio = (w, h) => {
        const r = w / h
        return allowedRatios.some(a => Math.abs(r - a) / a < tolerance)
      }

      if (file.type.startsWith('image/')) {
        const img = new Image()
        img.onload = () => {
          if (!checkRatio(img.width, img.height)) {
            this.fileError = 'Allowed image aspect ratios: 16:9, 4:3, 3:4, 1:1'
            toast.error(this.fileError)
            this.file = null
          } else {
            this.file = file
          }
          URL.revokeObjectURL(url)
        }
        img.onerror = () => {
          this.fileError = 'Invalid image file'
          toast.error(this.fileError)
          this.file = null
          URL.revokeObjectURL(url)
        }
        img.src = url
      } else if (file.type === 'video/mp4') {
        const vid = document.createElement('video')
        vid.preload = 'metadata'
        vid.onloadedmetadata = () => {
          const duration = vid.duration
          const width = vid.videoWidth
          const height = vid.videoHeight
          if (duration < 1 || duration > 30) {
            this.fileError = 'Video must be between 1 and 30 seconds.'
            toast.error(this.fileError)
            this.file = null
          } else if (!checkRatio(width, height)) {
            this.fileError = 'Allowed video aspect ratios: 16:9, 4:3, 3:4, 1:1'
            toast.error(this.fileError)
            this.file = null
          } else {
            this.file = file
          }
          URL.revokeObjectURL(url)
        }
        vid.onerror = () => {
          this.fileError = 'Invalid video file'
          toast.error(this.fileError)
          this.file = null
          URL.revokeObjectURL(url)
        }
        vid.src = url
      }
    },
    async submitThread() {
      this.loading = true
      try {
        let mediaUrl = ''
        if (this.file) {
          const form = new FormData()
          form.append('file', this.file)
          form.append('upload_preset', import.meta.env.VITE_CLOUDINARY_UPLOAD_PRESET)
          const res = await this.$axios.post(
              `https://api.cloudinary.com/v1_1/${import.meta.env.VITE_CLOUDINARY_CLOUD_NAME}/upload`,
              form
          )
          mediaUrl = res.data.secure_url
        }
        await this.$axios.post(
            `${import.meta.env.VITE_API_URL}/user/threads`,
            { content: this.content, media_url: mediaUrl },
            { headers: { Authorization: `Bearer ${localStorage.getItem('token')}` } }
        )
        toast.success('Thread posted')
        this.$router.push({ name: 'home' })
      } catch (err) {
        console.error('Create thread error', err)
        toast.error('Failed to post thread')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700&display=swap');

.create-thread-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(to bottom right, #e3f2fd, #f1f8e9);
  padding: 20px;
}

.create-thread-card {
  font-family: 'Nunito', sans-serif;
  width: 100%;
  max-width: 700px;
  background: #ffffff;
  padding: 30px 40px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.create-thread-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.back-btn {
  background: none;
  border: none;
  color: #3498db;
  cursor: pointer;
  font-size: 1rem;
  margin-bottom: 20px;
}

h1 {
  margin-bottom: 24px;
  font-weight: 700;
  color: #2c3e50;
  font-size: 2rem;
}

.thread-form .form-group {
  margin-bottom: 20px;
}
.thread-form label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #555;
}
.thread-form textarea,
.thread-form input[type="file"] {
  width: 100%;
  padding: 12px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 1rem;
  font-family: inherit;
}
.thread-form textarea {
  min-height: 140px;
  resize: vertical;
}

.error-text {
  color: #e74c3c;
  font-size: 0.9rem;
  margin-top: 6px;
}

.thread-form button {
  width: 100%;
  padding: 14px;
  background: #3498db;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  text-transform: uppercase;
  cursor: pointer;
  transition: background 0.3s, transform 0.2s;
}
.thread-form button:hover:not(:disabled) {
  background: #2980b9;
}
.thread-form button:active:not(:disabled) {
  transform: scale(0.98);
}
.thread-form button:disabled {
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

@media (max-width: 600px) {
  .create-thread-card {
    padding: 20px;
    border-radius: 12px;
  }
  h1 {
    font-size: 1.6rem;
    margin-bottom: 20px;
  }
  .thread-form textarea,
  .thread-form input[type="file"] {
    font-size: 0.95rem;
  }
  .thread-form button {
    font-size: 0.95rem;
    padding: 12px;
  }
}
</style>
