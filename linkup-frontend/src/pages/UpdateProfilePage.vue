<template>
  <div class="update-page">
    <div class="card">
      <h1>Edit Your Profile</h1>
      <form @submit.prevent="updateProfile" class="form">
        <div class="form-group">
          <label for="name">Name</label>
          <input id="name" v-model="form.name" required placeholder="Enter your name" />
        </div>
        <div class="form-group">
          <label for="bio">Bio</label>
          <textarea id="bio" v-model="form.bio" placeholder="Write something about yourself"></textarea>
        </div>
        <div class="form-group">
          <label for="location">Location</label>
          <input id="location" v-model="form.location" placeholder="City, Country" />
        </div>
        <div class="form-group">
          <label for="social-links">Social Links</label>
          <input id="social-links" v-model="form.social_links" placeholder="Add your social media links" />
        </div>
        <div class="form-group">
          <label for="avatar">Avatar Image</label>
          <input
              id="avatar"
              type="file"
              @change="onFileChange"
              accept="image/jpeg,image/png"
          />
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? 'Updating…' : 'Save Changes' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UpdateProfilePage',
  data() {
    return {
      form: {
        name: '',
        bio: '',
        location: '',
        social_links: '',
        avatar_url: ''
      },
      file: null,
      loading: false
    }
  },
  async created() {
    try {
      const token = localStorage.getItem('token')
      const { data } = await this.$axios.get(
          `${import.meta.env.VITE_API_URL}/user/me`,
          { headers: { Authorization: `Bearer ${token}` } }
      )
      Object.assign(this.form, data)
    } catch (err) {
      console.error(err)
    }
  },
  methods: {
    onFileChange(e) {
      const file = e.target.files[0] || null
      if (file) {
        const allowed = ['image/jpeg', 'image/png']
        if (!allowed.includes(file.type)) {
          alert('Please upload a JPG or PNG image.')
          this.file = null
          e.target.value = ''
          return
        }
      }
      this.file = file
    },
    async updateProfile() {
      this.loading = true
      try {
        let avatarUrl = this.form.avatar_url
        if (this.file) {
          const formData = new FormData()
          formData.append('file', this.file)
          formData.append('upload_preset', import.meta.env.VITE_CLOUDINARY_UPLOAD_PRESET)
          const { data } = await this.$axios.post(
              `https://api.cloudinary.com/v1_1/${import.meta.env.VITE_CLOUDINARY_CLOUD_NAME}/upload`,
              formData
          )
          avatarUrl = data.secure_url
        }
        const token = localStorage.getItem('token')
        await this.$axios.put(
            `${import.meta.env.VITE_API_URL}/user/me`,
            { ...this.form, avatar_url: avatarUrl },
            { headers: { Authorization: `Bearer ${token}` } }
        )
        this.$router.push({ name: 'my-profile' })
      } catch (err) {
        console.error(err)
        alert('Failed to update profile')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

*, *::before, *::after {
  box-sizing: border-box;
}

.update-page {
  font-family: 'Inter', sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background: linear-gradient(to bottom right, #e3f2fd, #f1f8e9);
  height: 100vh;
}

.card {
  background: #ffffff;
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 500px;
  text-align: center;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.card h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 24px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group label {
  display: block;
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 8px;
  color: #555;
  text-align: left;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 8px;
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: #42a5f5;
  box-shadow: 0 0 8px rgba(66, 165, 245, 0.3);
}

textarea {
  resize: none;
  height: 100px;
}

button {
  padding: 14px;
  font-size: 1rem;
  background: #42a5f5;
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  text-transform: uppercase;
  transition: background 0.3s ease, transform 0.1s ease;
}

button:hover:not(:disabled) {
  background: #1e88e5;
}

button:active:not(:disabled) {
  transform: scale(0.98);
}

button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .card {
    padding: 24px;
  }

  .card h1 {
    font-size: 1.8rem;
  }

  button,
  .form-group input,
  .form-group textarea {
    font-size: 0.9rem;
  }

  textarea {
    height: 80px;
  }
}
</style>
