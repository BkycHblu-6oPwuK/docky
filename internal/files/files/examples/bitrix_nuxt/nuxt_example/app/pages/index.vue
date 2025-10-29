<template>
  <section class="text-center py-20">
    <h1 class="text-5xl font-bold mb-4 text-indigo-700">Главная страница</h1>
    <p class="text-lg text-gray-600 mb-8">Добро пожаловать в Nuxt + Tailwind приложение</p>
    <img :src="'/upload/catalog/183/70qmt8o9moicedrf4hawhxi6duuwefxs.jpg'" alt="">
    <NuxtLink
      to="/catalog"
      class="inline-block bg-indigo-600 text-white px-6 py-3 rounded-lg shadow hover:bg-indigo-700 transition-all"
    >
      Перейти в каталог
    </NuxtLink>
  </section>
</template>
