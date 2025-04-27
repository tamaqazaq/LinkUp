<template>
  <div id="app">
    <Navbar />
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
import Navbar from './components/Navbar.vue'
export default { name: 'App', components: { Navbar } }
</script>

<style>
* { margin:0; padding:0; box-sizing: border-box; }
body { font-family: Arial, sans-serif; background: #f4f4f4; }
main { padding: 20px; }
</style>