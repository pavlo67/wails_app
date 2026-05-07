<script setup>
import { ref } from 'vue'

const apiBase = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:34116'
const loading = ref(false)
const error = ref('')
const result = ref('')

async function getHello() {
  loading.value = true
  error.value = ''
  result.value = ''

  try {
    const response = await fetch(`${apiBase}/api/hello`)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    const data = await response.json()
    result.value = `${data.message} / ${data.time}`
  } catch (e) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="container py-4 py-md-5">
    <div class="row justify-content-center">
      <div class="col-12 col-md-10 col-lg-8 col-xl-6">
        <div class="card shadow-sm">
          <div class="card-body p-4 p-md-5">
            <h1 class="h3 mb-3">Привіт 👋</h1>
            <p class="text-secondary mb-4">
              Wails + Vue + Go. Кнопка нижче робить GET-запит до Go-бека.
            </p>

            <button class="btn btn-primary btn-lg w-100" :disabled="loading" @click="getHello">
              {{ loading ? 'Питаю бек...' : 'GET /api/hello' }}
            </button>

            <div v-if="result" class="alert alert-success mt-4 mb-0">
              {{ result }}
            </div>

            <div v-if="error" class="alert alert-danger mt-4 mb-0">
              {{ error }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>
