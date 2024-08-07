// src/stores/walletStore.ts
import { defineStore } from "pinia";
import { useNotification } from "@kyvg/vue3-notification";
import axios from "axios";

const notification = useNotification();

export const useWalletStore = defineStore("wallet", {
  state: () => ({
    UserID: 0,
    Balance: 0.0,
    Address: "",
  }),
  actions: {
    async fetchWallet() {
      try {
        const user = window.Telegram.WebApp.initDataUnsafe.user;
        //TGUser = user;
        //const i = 6721747351;
        const response = await axios.get(`/api/v1/wallet?user_id=${user.id}`);
        //console.log(response.data);
        // Set state properties based on the response
        this.UserID = response.data.UserID;
        this.Balance = response.data.Balance / 1000000;
        this.Address = response.data.Address;
      } catch (error) {
        notification.notify({
          title: "Error",
          text: "Error fetching user",
          type: "error",
        });
      }
    },
  },
});
