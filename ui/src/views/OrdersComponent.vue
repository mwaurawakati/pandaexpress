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
          Orders
        </div>
      </div>
    </div>
    <div class="w-full mx-auto mt-8 p-4 bg-white shadow rounded-lg">
      <!--h2 class="text-2xl font-bold mb-4">Transactions</h2-->
      <div v-if="ordersStore.Orders.length > 0">
        <div style="overflow: scroll">
          <table class="min-w-full bg-white">
            <thead>
              <tr>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  Order ID
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  Total Amount
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  Shipping Fee
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  ...
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in ordersStore.Orders" :key="order.ID">
                <td class="py-2 px-4 border-b">
                  {{ order.ID }}
                </td>
                <td class="py-2 px-4 border-b">
                  {{ order.TotalPrice }}
                </td>
                <td class="py-2 px-4 border-b">{{ order.ShippingFee }}</td>
                <td class="py-2 px-4 border-b">
                  <button @click="navigateTo(`/orders/${order.id}`)">
                    View Order
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex justify-between mt-4">
          <button
            @click="ordersStore.fetchPrevTransactions()"
            :disabled="ordersStore.PrevOffset > -1"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Previous
          </button>
          <button
            @click="ordersStore.fetchTransactions()"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            :disabled="ordersStore.NextOffset != 0"
          >
            Next
          </button>
        </div>
      </div>
      <div v-else>
        <p>No transactions found.</p>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from "vue-router";
import { ArrowLeftIcon } from "@heroicons/vue/24/outline";
import { useOrdersStore } from "@/store/ordersStore";
import { onMounted } from "vue";

const ordersStore = useOrdersStore();
const router = useRouter();

const navigateTo = (path: string) => {
  router.push(path);
};

onMounted(() => {
  ordersStore.fetchOrders();
});
</script>
