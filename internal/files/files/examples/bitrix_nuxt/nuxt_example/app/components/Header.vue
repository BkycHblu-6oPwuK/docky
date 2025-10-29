<!-- components/Header.vue -->
<template>
  <header class="bg-white border-b border-gray-200 shadow-sm">
    <div class="container mx-auto flex items-center justify-between px-4 py-4">
      <NuxtLink to="/" class="text-xl font-semibold text-indigo-600 hover:text-indigo-800 transition-colors">
        MyNuxtApp
      </NuxtLink>

      <nav class="flex items-center gap-6">
        <NuxtLink
          to="/"
          class="text-gray-700 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 font-medium"
        >
          Главная
        </NuxtLink>

        <NuxtLink
          to="/catalog"
          class="text-gray-700 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 font-medium"
        >
          Каталог
        </NuxtLink>

        <NuxtLink
          to="/profile"
          class="text-gray-700 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 font-medium"
        >
          Профиль
        </NuxtLink>
      </nav>
    </div>
  </header>
</template>
