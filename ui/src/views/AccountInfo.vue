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
          Account Info
        </div>
      </div>
    </div>
    <!------------------------>
    <div class="w-full mx-auto p-4 flex flex-col items-baseline">
      <div
        class="bg-gray-300 shadow-md rounded-lg p-6 w-full flex flex-col items-baseline"
      >
        <h2 class="text-2xl font-bold mb-4">User Information</h2>
        <div class="mb-4">
          <span class="font-semibold">ID:</span> {{ userStore.ID }}
        </div>
        <div class="mb-4">
          <span class="font-semibold">Telegram Username:</span>
          {{ userStore.UserName }}
        </div>
        <div class="mb-4">
          <span class="font-semibold">Telegram Name:</span>
          {{ userStore.FirstName }}
          {{ userStore.LastName }}
        </div>
        <div class="mb-4">
          <span class="font-semibold">Preferred Language:</span>
          {{ userStore.PreferredLanguage.EnglishName }}
        </div>
        <div class="mb-4">
          <span class="font-semibold">Referee:</span> {{ userStore.Referre }}
        </div>
        <div class="mb-4">
          <span class="font-semibold">Referrals:</span>
          {{ userStore.Referrals.length }}
        </div>
        <div class="flex items-center">
          <span class="font-semibold">Referral Code:</span>
          <span class="ml-2">{{ userStore.ReferralCode }}</span>
          <button
            @click="copyToClipboard(userStore.ReferralCode)"
            class="ml-4 text-white px-2 py-1 rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-300"
          >
            <ClipboardIcon class="h-6 w-6 text-blue-600" />
          </button>
        </div>
        <div class="mb-4">
          <h3 class="text-xl font-semibold mb-2">Shipping Details:</h3>
          <p class="text-left">Name : {{ userStore.ShippingDetails.Name }}</p>
          <p class="text-left">Email: {{ userStore.ShippingDetails.Email }}</p>
          <p class="text-left">Phone: {{ userStore.ShippingDetails.Phone }}</p>
          <p class="text-left">Email: {{ userStore.ShippingDetails.Email }}</p>
          <p class="text-left">
            Continent: {{ userStore.ShippingDetails.Continent }}
          </p>
          <p class="text-left">
            Country: {{ userStore.ShippingDetails.Country }}
          </p>
          <p class="text-left">City: {{ userStore.ShippingDetails.City }}</p>
          <p
            class="text-left"
            v-if="userStore.ShippingDetails.Addresses.length > 0"
          >
            Address 1: {{ userStore.ShippingDetails.Addresses[0] }}
          </p>
          <p
            class="text-left"
            v-if="userStore.ShippingDetails.Addresses.length > 1"
          >
            Address 2: {{ userStore.ShippingDetails.Addresses[1] }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { useRouter } from "vue-router";
import { ArrowLeftIcon, ClipboardIcon } from "@heroicons/vue/24/outline";
import { onMounted } from "vue";
import { useNotification } from "@kyvg/vue3-notification";
const notification = useNotification();
import { useUserStore } from "@/store/userStore";
const userStore = useUserStore();
onMounted(async () => {
  await userStore.fetchUser();
});
const router = useRouter();
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text);
    notification.notify({
      title: "Success",
      text: "Referral code copied to clipboard!",
      type: "success",
    });
  } catch (error) {
    notification.notify({
      title: "Error",
      text: "Failed to copy referral code.",
      type: "error",
    });
  }
};
const navigateTo = (path) => {
  router.push(path);
};
</script>

<style scoped>
.max-w-4xl {
  max-width: 64rem;
}

.bg-white {
  background-color: white;
}

.shadow-md {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.rounded-lg {
  border-radius: 0.5rem;
}

.p-4 {
  padding: 1rem;
}

.p-6 {
  padding: 1.5rem;
}

.mb-4 {
  margin-bottom: 1rem;
}

.font-bold {
  font-weight: bold;
}

.text-2xl {
  font-size: 1.5rem;
}

.text-xl {
  font-size: 1.25rem;
}

.font-semibold {
  font-weight: 600;
}

.list-disc {
  list-style-type: disc;
}

.list-inside {
  list-style-position: inside;
}
</style>
