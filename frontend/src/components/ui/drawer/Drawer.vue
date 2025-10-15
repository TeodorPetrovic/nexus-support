<script lang="ts" setup>
import type { DrawerRootEmits, DrawerRootProps } from "vaul-vue"
import { DrawerRoot } from "vaul-vue"

const props = defineProps<DrawerRootProps>()

const emits = defineEmits<DrawerRootEmits>()
</script>

<template>
  <DrawerRoot
    data-slot="drawer"
    :should-scale-background="true"
    v-bind="props as any"
    @update:open="(val) => emits('update:open', val)"
  >
    <slot />
  </DrawerRoot>
</template>
