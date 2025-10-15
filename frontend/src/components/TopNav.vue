<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { BrowseForPC } from '../../wailsjs/go/main/App'

const router = useRouter()
const pcName = ref('')
const isDialogOpen = ref(false)

const handleBrowseForPC = async () => {
  try {
    const clients = await BrowseForPC(pcName.value)
    console.log('Found PCs:', clients)
    isDialogOpen.value = false
    pcName.value = ''
    // Navigate to discovery page to show results
    router.push('/discovery')
  } catch (error) {
    console.error('Browse for PC failed:', error)
    alert('Failed to browse for PCs: ' + error)
  }
}
</script>

<template>
  <header class="flex h-16 shrink-0 items-center gap-2 border-b px-6">
    <div class="flex flex-1 items-center gap-2">
      <h1 class="text-xl font-semibold">Nexus Support</h1>
    </div>
    <div class="flex items-center gap-4">
      <Dialog v-model:open="isDialogOpen">
        <DialogTrigger as-child>
          <Button variant="default">
            Browse for PC
          </Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Browse for PC</DialogTitle>
            <DialogDescription>
              Enter a name to browse for PCs on the network.
            </DialogDescription>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid gap-2">
              <Label for="pc-name">PC Name</Label>
              <Input
                id="pc-name"
                v-model="pcName"
                placeholder="Enter PC name..."
                @keyup.enter="handleBrowseForPC"
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="isDialogOpen = false">
              Cancel
            </Button>
            <Button type="button" @click="handleBrowseForPC">
              Browse
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  </header>
</template>
