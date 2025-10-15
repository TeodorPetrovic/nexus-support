<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, Monitor, User, HardDrive, Cpu, Video } from 'lucide-vue-next'
import { DiscoverClients, GetClientInfo } from '../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'

const router = useRouter()
import type { main } from '../../wailsjs/go/models'

type ClientInfo = main.ClientInfo
type SystemInfo = main.SystemInfo

const clients = ref<ClientInfo[]>([])
const isDiscovering = ref(false)
const selectedClient = ref<ClientInfo | null>(null)
const clientSystemInfo = ref<SystemInfo | null>(null)
const isLoadingInfo = ref(false)

const discoverPCs = async () => {
  isDiscovering.value = true
  try {
    const discovered = await DiscoverClients()
    clients.value = discovered || []
  } catch (error) {
    console.error('Discovery failed:', error)
  } finally {
    isDiscovering.value = false
  }
}

const selectClient = async (client: ClientInfo) => {
  selectedClient.value = client
  isLoadingInfo.value = true
  try {
    const info = await GetClientInfo(client.ip)
    clientSystemInfo.value = info
  } catch (error) {
    console.error('Failed to get client info:', error)
  } finally {
    isLoadingInfo.value = false
  }
}

const viewStream = (client: ClientInfo) => {
  router.push(`/stream/${client.ip}/${client.hostname}`)
}

onMounted(() => {
  discoverPCs()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">PC Discovery</h1>
        <p class="text-muted-foreground mt-2">Discover and manage connected PCs</p>
      </div>
      <Button @click="discoverPCs" :disabled="isDiscovering">
        <Loader2 v-if="isDiscovering" class="mr-2 h-4 w-4 animate-spin" />
        <Monitor v-else class="mr-2 h-4 w-4" />
        {{ isDiscovering ? 'Discovering...' : 'Discover PCs' }}
      </Button>
    </div>

    <div class="grid gap-6 lg:grid-cols-3">
      <!-- Client List -->
      <div class="lg:col-span-2">
        <Card>
          <CardHeader>
            <CardTitle>Available PCs ({{ clients.length }})</CardTitle>
            <CardDescription>Click on a PC to view details</CardDescription>
          </CardHeader>
          <CardContent>
            <div v-if="clients.length === 0" class="text-center py-8">
              <Monitor class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
              <p class="text-muted-foreground">No PCs found. Make sure client is running on target machines.</p>
            </div>
            <div v-else class="space-y-2">
              <div 
                v-for="client in clients" 
                :key="client.id"
                @click="selectClient(client)"
                class="flex items-center justify-between p-3 border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                :class="selectedClient?.id === client.id ? 'bg-accent border-primary' : ''"
              >
                <div class="flex items-center gap-3">
                  <Monitor class="h-5 w-5" />
                  <div>
                    <p class="font-medium">{{ client.hostname }}</p>
                    <p class="text-sm text-muted-foreground">{{ client.ip }}</p>
                  </div>
                </div>
                <Badge>{{ client.status }}</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Client Details -->
      <div>
        <Card>
          <CardHeader>
            <CardTitle>PC Details</CardTitle>
            <CardDescription>System information</CardDescription>
          </CardHeader>
          <CardContent>
            <div v-if="!selectedClient" class="text-center py-8">
              <p class="text-muted-foreground">Select a PC to view details</p>
            </div>
            <div v-else-if="isLoadingInfo" class="text-center py-8">
              <Loader2 class="mx-auto h-8 w-8 animate-spin mb-2" />
              <p class="text-muted-foreground">Loading system info...</p>
            </div>
            <div v-else-if="clientSystemInfo" class="space-y-4">
              <div class="flex items-center gap-2">
                <Monitor class="h-4 w-4" />
                <span class="font-medium">{{ clientSystemInfo.hostname }}</span>
              </div>
              <div class="flex items-center gap-2">
                <User class="h-4 w-4" />
                <span>{{ clientSystemInfo.user }}</span>
              </div>
              <div class="flex items-center gap-2">
                <HardDrive class="h-4 w-4" />
                <span class="text-sm">{{ clientSystemInfo.os }}</span>
              </div>
              <div class="flex items-center gap-2">
                <Cpu class="h-4 w-4" />
                <span>{{ clientSystemInfo.cpu }} cores</span>
              </div>
              <div class="flex items-center gap-2">
                <HardDrive class="h-4 w-4" />
                <span>{{ clientSystemInfo.memMB }}MB RAM</span>
              </div>
              
              <div class="pt-4 space-y-2">
                <Button 
                  @click="viewStream(selectedClient)" 
                  class="w-full"
                  variant="default"
                >
                  <Video class="mr-2 h-4 w-4" />
                  View Live Stream
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
