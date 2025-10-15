<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed } from 'vue'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarHeader,
  SidebarFooter
} from '@/components/ui/sidebar'
import { Home, FileText } from 'lucide-vue-next'

const route = useRoute()

const menuItems = [
  {
    title: 'PC Discovery',
    icon: Home,
    to: '/discovery'
  },
  {
    title: 'Screenshots',
    icon: FileText,
    to: '/screenshots'
  }
]

const isActive = (path: string) => {
  return route.path === path
}
</script>

<template>
  <Sidebar>
    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Pages</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in menuItems" :key="item.to">
              <SidebarMenuButton
                as="a"
                :href="`#${item.to}`"
                :data-active="isActive(item.to)"
              >
                <component :is="item.icon" />
                <span>{{ item.title }}</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>
    
    <SidebarFooter class="border-t p-4">
      <p class="text-xs text-muted-foreground">Nexus Support v1.0</p>
    </SidebarFooter>
  </Sidebar>
</template>
