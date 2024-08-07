// src/stores/walletStore.ts
import { defineStore } from "pinia";
import { useNotification } from "@kyvg/vue3-notification";
import axios from "axios";

const notification = useNotification();

export const useOrdersStore = defineStore("orders", {
  state: () => ({
    NextOffset: 0,
    Orders: [],
    PrevOffset: -1,
  }),
  actions: {
    async fetchOrders() {
      try {
        const user = window.Telegram.WebApp.initDataUnsafe.user;
        //TGUser = user;
        //const i = 6721747351;
        const response = await axios.get(
          `/api/v1/orders?user_id=${user.id}&offset=${this.NextOffset}&num=50`
        );
        // Set state properties based on the response
        this.PrevOffset = this.NextOffset;
        this.NextOffset = response.data.NextOffset;
        this.Transactions = response.data.Transactions
          ? response.data.Transactions
          : [];
      } catch (error) {
        notification.notify({
          title: "Error",
          text: "Error fetching orders",
          type: "error",
        });
      }
    },
    async fetchPrevOrders() {
      try {
        //const user = window.Telegram.WebApp.initDataUnsafe.user;
        //TGUser = user;
        const i = 6721747351;
        const response = await axios.get(
          `/api/v1/transactions?user_id=${i}&offset=${this.NextOffset}&num=50`
        );
        // Set state properties based on the response
        this.PrevOffset = this.NextOffset;
        this.NextOffset = response.data.NextOffset;
        this.Transactions = response.data.Transactions;
      } catch (error) {
        notification.notify({
          title: "Error",
          text: "Error fetching orders",
          type: "error",
        });
      }
    },
  },
});
