<template>
  <div class="verified-page">
    <div class="verified-card">
      <h1>🎉 Email Verified!</h1>
      <p>Your email address has been successfully verified.</p>
      <router-link to="/login" class="login-button">Go to Login</router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EmailVerifiedPage',
}
</script>

<style scoped>
.verified-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
  font-family: 'Nunito', sans-serif;
}

.verified-card {
  background: #fff;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-align: center;
  width: 100%;
  max-width: 400px;
}

.verified-card h1 {
  color: #2ecc71;
  margin-bottom: 16px;
}

.verified-card p {
  color: #555;
  margin-bottom: 24px;
}

.login-button {
  display: inline-block;
  padding: 10px 20px;
  background-color: #3498db;
  color: white;
  border-radius: 4px;
  text-decoration: none;
  font-weight: bold;
}

.login-button:hover {
  background-color: #2980b9;
}
</style>
