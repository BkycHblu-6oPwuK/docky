<!-- layouts/default.vue -->
<template>
  <div class="flex flex-col min-h-screen bg-gray-50 text-gray-900">
    <Header />

    <main class="flex-1 container mx-auto px-4 py-8">
      <slot />
    </main>

    <Footer />
  </div>
</template>

<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
</script>
