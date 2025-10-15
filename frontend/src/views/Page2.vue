<script setup lang="ts">
import { ref } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Camera, Loader2, Download } from 'lucide-vue-next'
import { GetScreenshot } from '../../wailsjs/go/main/App'

const targetIP = ref('')
const screenshot = ref('')
const isCapturing = ref(false)
const lastCaptured = ref('')

const captureScreenshot = async () => {
  if (!targetIP.value) {
    alert('Please enter an IP address')
    return
  }
  
  isCapturing.value = true
  try {
    const base64Image = await GetScreenshot(targetIP.value)
    if (base64Image.startsWith('Error:')) {
      alert(base64Image)
    } else {
      screenshot.value = base64Image
      lastCaptured.value = new Date().toLocaleString()
    }
  } catch (error) {
    alert('Failed to capture screenshot: ' + error)
  } finally {
    isCapturing.value = false
  }
}

const downloadScreenshot = () => {
  if (!screenshot.value) return
  
  const link = document.createElement('a')
  link.href = `data:image/jpeg;base64,${screenshot.value}`
  link.download = `screenshot_${targetIP.value}_${new Date().getTime()}.jpg`
  link.click()
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold">Remote Screenshot</h1>
      <p class="text-muted-foreground mt-2">Capture screenshots from remote PCs</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Screenshot Controls</CardTitle>
        <CardDescription>Enter the IP address of the target PC</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid w-full max-w-sm items-center gap-1.5">
          <Label for="ip-address">Target IP Address</Label>
          <Input
            id="ip-address"
            v-model="targetIP"
            placeholder="192.168.1.100"
            type="text"
          />
        </div>
        <div class="flex gap-2">
          <Button @click="captureScreenshot" :disabled="isCapturing">
            <Loader2 v-if="isCapturing" class="mr-2 h-4 w-4 animate-spin" />
            <Camera v-else class="mr-2 h-4 w-4" />
            {{ isCapturing ? 'Capturing...' : 'Capture Screenshot' }}
          </Button>
          <Button 
            v-if="screenshot" 
            @click="downloadScreenshot" 
            variant="outline"
          >
            <Download class="mr-2 h-4 w-4" />
            Download
          </Button>
        </div>
        <div v-if="lastCaptured" class="text-sm text-muted-foreground">
          Last captured: {{ lastCaptured }}
        </div>
      </CardContent>
    </Card>

    <Card v-if="screenshot">
      <CardHeader>
        <CardTitle>Screenshot Preview</CardTitle>
        <CardDescription>Remote screen from {{ targetIP }}</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="border rounded-lg overflow-hidden">
          <img 
            :src="`data:image/jpeg;base64,${screenshot}`" 
            alt="Remote Screenshot"
            class="w-full h-auto max-h-96 object-contain"
          />
        </div>
      </CardContent>
    </Card>
  </div>
</template>
