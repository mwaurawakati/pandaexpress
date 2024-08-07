<template>
  <div class="container">
    <div
      class="flex sm:items-start m-4 border-b-solid border-b border-black pb-4"
    >
      <div
        class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10"
        @click="navigateTo('/account')"
      >
        <ArrowLeftIcon class="h-8 w-8 text-blue-600" aria-hidden="true" />
      </div>
      <div class="mx-auto mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
        <div as="h1" class="text-2xl font-semibold leading-6 text-gray-900">
          Wallet Info
        </div>
      </div>
    </div>
    <!---->
    <div class="w-full mx-auto p-4 flex flex-col items-baseline">
      <div v-if="walletStore.Address" class="flex flex-col items-baseline">
        <p>
          <span class="font-semibold">Address:</span>
          {{ walletStore.Address }}
        </p>
        <p>
          <span class="font-semibold">Balance:</span>
          {{ walletStore.Balance }}
        </p>
      </div>
      <div v-else>
        <p>Loading...</p>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from "vue-router";
import { ArrowLeftIcon } from "@heroicons/vue/24/outline";
import { onMounted } from "vue";
import { useWalletStore } from "@/store/walletStore";
const walletStore = useWalletStore();
onMounted(async () => {
  await walletStore.fetchWallet();
});
const router = useRouter();

const navigateTo = (path: string) => {
  router.push(path);
};
</script>
