// src/stores/userStore.ts
import { defineStore } from "pinia";
import { useNotification } from "@kyvg/vue3-notification";
import axios from "axios";

const notification = useNotification();

export const useUserStore = defineStore("user", {
  state: () => ({
    ID: 0,
    UserName: "",
    FirstName: "",
    LastName: "",
    PreferredLanguage: {
      Code: "",
      EnglishName: "",
      Name: "",
    },
    ShippingDetails: {
      Name: "",
      Continent: "",
      Country: "",
      City: "",
      Email: "",
      Phone: "",
      Street: "",
      AppartmentNumber: "",
      Addresses: [],
    },
    Referee: "",
    Referrals: [],
    ReferralCode: "refTEST",
  }),
  actions: {
    async fetchUser() {
      try {
        const user = window.Telegram.WebApp.initDataUnsafe.user;
        //TGUser = user;
        //const i = 6721747351;
        const response = await axios.get(`/api/v1/user?user_id=${user.id}`);
        //console.log(response.data);

        // Set state properties based on the response
        this.ID = response.data.ID;
        this.UserName = response.data.UserName;
        this.FirstName = response.data.FirstName;
        this.LastName = response.data.LastName;
        this.PreferredLanguage = response.data.PreferredLanguage;
        this.ShippingDetails = response.data.ShippingDetails;
        this.Referee = response.data.Referee;
        this.Referrals = response.data.Referrals ? response.data.Referrals : [];
        this.ReferralCode = response.data.ReferralCode;
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
