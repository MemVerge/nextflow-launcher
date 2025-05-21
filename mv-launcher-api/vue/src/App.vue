<script>
export default {
  name: 'App'
}
</script>

<template>
  <div class="app">
    <nav class="navbar">
      <div class="nav-content">
        <div class="nav-logo">
          <img src="./assets/logo.png" alt="Memory Machine Logo" class="logo">
        </div>
        <div class="nav-links">
          <router-link to="/" class="nav-link" active-class="active">
            <span class="nav-icon">+</span>
            Create Job
          </router-link>
          <router-link to="/jobs" class="nav-link" active-class="active">
            <span class="nav-icon">📋</span>
            Job List
          </router-link>
        </div>
      </div>
    </nav>
    <main class="main-content">
      <div class="content-wrapper">
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<style>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f8f9fa;
}

.navbar {
  background-color: white;
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.nav-content {
  width: 100%;
  margin: 0;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-logo {
  height: 40px;
  display: flex;
  align-items: center;
}

.logo {
  height: 32px;
  width: auto;
  object-fit: contain;
}

.nav-links {
  display: flex;
  gap: 1rem;
  margin: 0;
  padding: 0;
}

.nav-link {
  color: #4a5568;
  text-decoration: none;
  font-weight: 500;
  padding: 0.75rem 1.25rem;
  border-radius: 8px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  position: relative;
  font-size: 0.95rem;
}

.nav-link:hover {
  background-color: #f0f4f8;
  color: #0066cc;
  transform: translateY(-1px);
}

.nav-link.active {
  color: #0066cc;
  background-color: #e6f0ff;
  font-weight: 600;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 3px;
  background-color: #0066cc;
  border-radius: 2px;
}

.nav-icon {
  font-size: 1.1rem;
  opacity: 0.9;
}

.main-content {
  flex: 1;
  padding-top: 80px;
  width: 100%;
  min-height: calc(100vh - 80px);
}

.content-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  width: 100%;
  height: 100%;
}

@media (max-width: 768px) {
  .nav-content {
    padding: 0 1rem;
  }
  
  .nav-link {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }
  
  .nav-icon {
    font-size: 1rem;
  }
}
</style>