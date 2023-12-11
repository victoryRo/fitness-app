<template>
  <div>
    <ul class="bg-amber-300 overflow-hidden">
      <li class="float-left p-3.5 bg-amber-300 text-center hover:bg-amber-400">
        <router-link :to="{name: 'Home'}">Home</router-link>
      </li>
      <li class="float-left p-3.5 bg-amber-300 text-center hover:bg-amber-400">
        <router-link :to="{name: 'Login'}">Login</router-link>
      </li>
    </ul>
    <router-view></router-view>
  </div>
</template>

