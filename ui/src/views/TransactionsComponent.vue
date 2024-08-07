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
          Transactions
        </div>
      </div>
    </div>
    <div class="w-full mx-auto mt-8 p-4 bg-white shadow rounded-lg">
      <!--h2 class="text-2xl font-bold mb-4">Transactions</h2-->
      <div v-if="transactionsStore.Transactions.length > 0">
        <div style="overflow: scroll">
          <table class="min-w-full bg-white">
            <thead>
              <tr>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  Type
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  Amount
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  To
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  From
                </th>
                <th class="py-2 px-4 border-b-2 border-gray-200 bg-gray-100">
                  Transaction ID
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="transaction in transactionsStore.Transactions"
                :key="transaction.TransactionID"
              >
                <td class="py-2 px-4 border-b">
                  {{ transactionTypeMap[transaction.Type] }}
                </td>
                <td class="py-2 px-4 border-b">
                  {{ transaction.Amount / 1000000 }}
                </td>
                <td class="py-2 px-4 border-b">{{ transaction.To }}</td>
                <td class="py-2 px-4 border-b">{{ transaction.From }}</td>
                <td class="py-2 px-4 border-b">
                  {{ transaction.transaction_id }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex justify-between mt-4">
          <button
            @click="transactionsStore.fetchPrevTransactions()"
            :disabled="transactionsStore.PrevOffset > -1"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Previous
          </button>
          <button
            @click="transactionsStore.fetchTransactions()"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            :disabled="transactionsStore.NextOffset != 0"
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
import { useTransactionsStore } from "@/store/transactionsStore";
import { onMounted } from "vue";
const transactionTypeMap: { [key: number]: string } = {
  0: "Deposit",
  1: "Withdraw",
  2: "Purchase",
  3: "Shipping Fee",
  4: "Referral Earning",
};
const transactionsStore = useTransactionsStore();
const router = useRouter();

const navigateTo = (path: string) => {
  router.push(path);
};

onMounted(() => {
  transactionsStore.fetchTransactions(0);
});
</script>
