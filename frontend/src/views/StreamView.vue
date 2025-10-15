<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, Play, Square, RefreshCw, Maximize2, Minimize2, X } from 'lucide-vue-next'
import { StartScreenStream, StopScreenStream, GetStreamFrame } from '../../wailsjs/go/main/App'

interface Props {
  ip: string
  hostname: string
}

const props = defineProps<Props>()
const router = useRouter()

const isStreaming = ref(false)
const streamFrame = ref('')
const isLoading = ref(false)
const lastUpdated = ref('')
const streamInterval = ref<NodeJS.Timeout | null>(null)
const frameCount = ref(0)
const framerate = ref<number>(10)
const actualFPS = ref(0)
const fpsCounter = ref(0)
const lastFPSCheck = ref(Date.now())
const isFullscreen = ref(false)
const streamContainer = ref<HTMLElement | null>(null)

const startStream = async () => {
  isLoading.value = true
  try {
    await StartScreenStream(props.ip, parseInt(framerate.value.toString()))
    isStreaming.value = true
    frameCount.value = 0
    fpsCounter.value = 0
    lastFPSCheck.value = Date.now()
    
    // Calculate interval from framerate (1000ms / fps)
    const interval = Math.max(16, Math.floor(1000 / framerate.value)) // Min 16ms (60fps max)
    
    streamInterval.value = setInterval(async () => {
      try {
        const frame = await GetStreamFrame(props.ip)
        if (!frame.startsWith('Error:')) {
          streamFrame.value = frame
          lastUpdated.value = new Date().toLocaleTimeString()
          frameCount.value++
          fpsCounter.value++
          
          // Calculate actual FPS every second
          const now = Date.now()
          if (now - lastFPSCheck.value >= 1000) {
            actualFPS.value = Math.round(fpsCounter.value * 1000 / (now - lastFPSCheck.value))
            fpsCounter.value = 0
            lastFPSCheck.value = now
          }
        }
      } catch (error) {
        console.error('Frame error:', error)
      }
    }, interval)
    
  } catch (error) {
    console.error('Failed to start stream:', error)
  } finally {
    isLoading.value = false
  }
}

const stopStream = async () => {
  isLoading.value = true
  try {
    if (streamInterval.value) {
      clearInterval(streamInterval.value)
      streamInterval.value = null
    }
    await StopScreenStream(props.ip)
    isStreaming.value = false
    frameCount.value = 0
  } catch (error) {
    console.error('Failed to stop stream:', error)
  } finally {
    isLoading.value = false
  }
}

const goBack = () => {
  if (isStreaming.value) {
    stopStream()
  }
  router.push('/discovery')
}

const toggleFullscreen = () => {
  if (!isFullscreen.value) {
    // Enter fullscreen
    if (streamContainer.value?.requestFullscreen) {
      streamContainer.value.requestFullscreen()
    }
  } else {
    // Exit fullscreen
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
  }
}

const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement
}

const closeApp = () => {
  if (confirm('Are you sure you want to close the application?')) {
    window.close()
  }
}

onMounted(() => {
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

onUnmounted(() => {
  if (streamInterval.value) {
    clearInterval(streamInterval.value)
  }
  if (isStreaming.value) {
    StopScreenStream(props.ip)
  }
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
})
</script>

<template>
  <div class="space-y-6" :class="isFullscreen ? 'fixed inset-0 z-50 bg-black' : ''">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <Button variant="outline" @click="goBack">
          <ArrowLeft class="mr-2 h-4 w-4" />
          Back to Discovery
        </Button>
        <div>
          <h1 class="text-3xl font-bold">Live Stream</h1>
          <p class="text-muted-foreground mt-2">{{ hostname }} ({{ ip }})</p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <Badge v-if="isStreaming" variant="default" class="bg-green-500">
          <div class="w-2 h-2 bg-white rounded-full mr-2 animate-pulse"></div>
          Live
        </Badge>
        <Badge v-else variant="outline">
          Stopped
        </Badge>
        <Button variant="outline" size="sm" @click="toggleFullscreen">
          <Maximize2 v-if="!isFullscreen" class="h-4 w-4" />
          <Minimize2 v-else class="h-4 w-4" />
        </Button>
        <Button variant="outline" size="sm" @click="closeApp">
          <X class="h-4 w-4" />
        </Button>
      </div>
    </div>

    <div class="grid gap-6 lg:grid-cols-4" v-if="!isFullscreen">
      <!-- Stream Control ---->
      <div class="lg:col-span-1">
        <Card>
          <CardHeader>
            <CardTitle>Stream Controls</CardTitle>
            <CardDescription>Manage live screen streaming</CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div v-if="!isStreaming" class="space-y-3">
              <div class="space-y-2">
                <label class="text-sm font-medium">Framerate: {{ framerate }} FPS</label>
                <input 
                  type="range" 
                  v-model.number="framerate" 
                  min="1" 
                  max="60" 
                  step="1"
                  class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                />
                <div class="flex justify-between text-xs text-muted-foreground">
                  <span>1 FPS</span>
                  <span>60 FPS</span>
                </div>
              </div>
              <Button 
                @click="startStream" 
                :disabled="isLoading"
                class="w-full"
              >
                <Play class="mr-2 h-4 w-4" />
                Start Stream
              </Button>
            </div>
            <Button 
              v-else 
              @click="stopStream" 
              :disabled="isLoading"
              variant="destructive"
              class="w-full"
            >
              <Square class="mr-2 h-4 w-4" />
              Stop Stream
            </Button>
            
            <div v-if="isStreaming" class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span>Status:</span>
                <Badge variant="default" class="bg-green-500">Active</Badge>
              </div>
              <div class="flex justify-between">
                <span>Target FPS:</span>
                <span>{{ framerate }}</span>
              </div>
              <div class="flex justify-between">
                <span>Actual FPS:</span>
                <span>{{ actualFPS }}</span>
              </div>
              <div class="flex justify-between">
                <span>Frames:</span>
                <span>{{ frameCount }}</span>
              </div>
              <div class="flex justify-between">
                <span>Last Update:</span>
                <span>{{ lastUpdated }}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Stream Display -->
      <div class="lg:col-span-3">
        <Card>
          <CardHeader v-if="!isFullscreen">
            <CardTitle>Remote Desktop</CardTitle>
            <CardDescription>Live view of {{ hostname }}</CardDescription>
          </CardHeader>
          <CardContent :class="isFullscreen ? 'p-0' : ''">
            <div 
              ref="streamContainer"
              class="border rounded-lg overflow-hidden bg-black flex items-center justify-center relative"
              :class="isFullscreen ? 'h-screen w-screen border-none rounded-none' : 'min-h-96'"
            >
              <!-- Fullscreen controls overlay -->
              <div v-if="isFullscreen" class="absolute top-4 right-4 z-10 flex gap-2">
                <Button variant="outline" size="sm" @click="toggleFullscreen" class="bg-black/50 hover:bg-black/70">
                  <Minimize2 class="h-4 w-4 text-white" />
                </Button>
                <Button variant="outline" size="sm" @click="closeApp" class="bg-black/50 hover:bg-black/70">
                  <X class="h-4 w-4 text-white" />
                </Button>
              </div>
              
              <div v-if="!isStreaming && !streamFrame" class="text-center text-muted-foreground">
                <Play class="mx-auto h-16 w-16 mb-4 opacity-50" />
                <p class="text-lg">Click "Start Stream" to begin live streaming</p>
              </div>
              <div v-else-if="isLoading" class="text-center text-white">
                <RefreshCw class="mx-auto h-8 w-8 mb-2 animate-spin" />
                <p>{{ isStreaming ? 'Loading stream...' : 'Starting stream...' }}</p>
              </div>
              <img 
                v-else-if="streamFrame"
                :src="`data:image/jpeg;base64,${streamFrame}`" 
                alt="Live Stream"
                :class="isFullscreen ? 'w-full h-auto max-h-screen object-contain' : 'w-full h-auto max-h-96 object-contain'"
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
    
    <!-- Fullscreen Stream Display -->
    <div v-if="isFullscreen" class="fixed inset-0 z-50 bg-black flex items-center justify-center">
      <!-- Fullscreen controls overlay -->
      <div class="absolute top-4 right-4 z-10 flex gap-2">
        <Button variant="outline" size="sm" @click="toggleFullscreen" class="bg-black/50 hover:bg-black/70 text-white border-white/30">
          <Minimize2 class="h-4 w-4" />
        </Button>
        <Button variant="outline" size="sm" @click="closeApp" class="bg-black/50 hover:bg-black/70 text-white border-white/30">
          <X class="h-4 w-4" />
        </Button>
      </div>
      
      <div v-if="!isStreaming && !streamFrame" class="text-center text-white">
        <Play class="mx-auto h-16 w-16 mb-4 opacity-50" />
        <p class="text-lg">Click "Start Stream" to begin live streaming</p>
      </div>
      <div v-else-if="isLoading" class="text-center text-white">
        <RefreshCw class="mx-auto h-8 w-8 mb-2 animate-spin" />
        <p>{{ isStreaming ? 'Loading stream...' : 'Starting stream...' }}</p>
      </div>
      <img 
        v-else-if="streamFrame"
        :src="`data:image/jpeg;base64,${streamFrame}`" 
        alt="Live Stream"
        class="w-full h-auto max-h-screen object-contain"
      />
    </div>
  </div>
</template>

<style scoped>
.animate-pulse {
  animation: pulse 1s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>