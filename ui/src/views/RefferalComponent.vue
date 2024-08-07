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
          Referrals
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from "vue-router";
import { ArrowLeftIcon } from "@heroicons/vue/24/outline";
const router = useRouter();

const navigateTo = (path: string) => {
  router.push(path);
};
</script>
