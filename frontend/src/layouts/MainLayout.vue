<script setup lang="ts">
import { SidebarProvider, SidebarInset, SidebarTrigger } from '@/components/ui/sidebar'
import SidebarNav from '@/components/SidebarNav.vue'
import TopNav from '@/components/TopNav.vue'
</script>

<template>
  <SidebarProvider>
    <SidebarNav />
    <SidebarInset>
      <TopNav />
      <main class="flex-1 overflow-auto p-6">
        <router-view />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>

<style scoped>
main {
  min-height: calc(100vh - 4rem);
}
</style>
