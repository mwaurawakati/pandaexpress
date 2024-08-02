<template>
  <div class="container">
    <div style="display: flex; flex-grow: 1; flex-direction: column">
      <div
        class="h-16 border-b border-b-solid w-full text-4xl items-center justify-center flex inline-block bg-white shadow-2xl"
      >
        Settings
      </div>
      <div style="flex-grow: 1" class="my-4">
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="open = true"
        >
          <span class="font-bold text-2xl">Language</span>
          <span>Set preffered language</span>
        </div>
        <div class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col">
          <span class="font-bold text-2xl">Shipping Address</span>
          <span>Set shipping Address</span>
        </div>
      </div>
    </div>
    <SetShippingAddress
      :open="isDialogOpen"
      @update:open="isDialogOpen = $event"
    />
  </div>
  <NavComponent />
</template>
<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { PhotoIcon, UserCircleIcon } from "@heroicons/vue/24/solid";
import NavComponent from "./NavComponent.vue";
import TestC from "./TestC.vue";
import SetShippingAddress from "./SetShippingAddress.vue";
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  TransitionChild,
  TransitionRoot,
} from "@headlessui/vue";
@Options({
  components: {
    NavComponent,
    TestC,
    SetShippingAddress,
  },
})
export default class SettingsView extends Vue {
  open = false;
}
</script>
