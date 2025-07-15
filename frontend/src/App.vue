<script setup lang="ts">
import { ref, onMounted } from 'vue'

const name = ref('')
const greeting = ref('')
const isLoading = ref(false)
const error = ref('')
const isServerHealthy = ref(false)

async function sendGreeting() {
  if (!name.value.trim()) {
    error.value = 'Пожалуйста, введите имя11'
    return
  }

  isLoading.value = true
  error.value = ''
  greeting.value = ''

  try {
    // const response = await fetch(`/api/greet?name=${encodeURIComponent(name.value)}`, {
    //   headers: {
    //     'Accept': 'application/json',
    //   }
    // })
    
    // if (!response.ok) {
    //   const errorData = await response.json().catch(() => ({}))
    //   throw new Error(errorData.detail || `Ошибка сервера: ${response.status}`)
    // }

    // const data = await response.json()
    // greeting.value = data.message
  } catch (err) {
    console.error('Ошибка запроса:', err)
    if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'Произошла неизвестная ошибка'
    }
    
    // Проверяем здоровье сервера после ошибки
    await checkHealth()
  } finally {
    isLoading.value = false
  }
}

async function checkHealth(): Promise<boolean> {
  try {
    // const response = await fetch('/api/health', {
    //   signal: AbortSignal.timeout(3000) // Таймаут 3 секунды
    // })
    
    // if (!response.ok) {
    //   throw new Error('Сервер не отвечает')
    // }
    
    isServerHealthy.value = true
    return true
  } catch (err) {
    console.error('Health check failed:', err)
    isServerHealthy.value = false
    
    if (err instanceof Error) {
      error.value = 'Серверная часть недоступна: ' + err.message
    }
    
    return false
  }
}

// Проверка соединения при загрузке компонента
onMounted(async () => {
  await checkHealth()
})
</script>

<template>
  <div class="container">
    <h1>Приветственное приложение</h1>
    
    <div v-if="!isServerHealthy" class="server-warning">
      ⚠ Серверная часть недоступна. Некоторые функции могут не работать.
    </div>
    
    <div class="input-group">
      <input
        v-model.trim="name"
        placeholder="Введите ваше имяffff"
        :disabled="isLoading || !isServerHealthy"
        @keyup.enter="sendGreeting"
      />
      <button
        @click="sendGreeting"
        :disabled="isLoading || !name.trim() || !isServerHealthy"
      >
        <span v-if="isLoading">Отправка...</span>
        <span v-else>Отправитьff98</span>
      </button>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="greeting" class="greeting-message">
      {{ greeting }}
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}

.server-warning {
  padding: 10px;
  margin-bottom: 15px;
  background-color: #fff3cd;
  color: #856404;
  border-radius: 4px;
  border-left: 4px solid #ffeeba;
}

.input-group {
  display: flex;
  gap: 10px;
  margin: 20px 0;
}

input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

button {
  padding: 8px 16px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.2s;
}

button:hover:not(:disabled) {
  background-color: #369f6b;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background-color: #ffebee;
  color: #f44336;
  border-radius: 4px;
  border-left: 4px solid #f44336;
}

.greeting-message {
  margin-top: 15px;
  padding: 10px;
  background-color: #e8f5e9;
  color: #2e7d32;
  border-radius: 4px;
  border-left: 4px solid #2e7d32;
}
</style>