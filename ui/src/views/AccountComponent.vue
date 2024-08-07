<template>
  <div class="container">
    <div style="display: flex; flex-grow: 1; flex-direction: column">
      <div
        class="h-16 border-b border-b-solid w-full text-4xl items-center justify-center flex inline-block bg-white shadow-2xl"
      >
        Account
      </div>
      <div style="flex-grow: 1" class="my-4">
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="navigateTo('/account_info')"
        >
          <span class="font-bold text-2xl">Account Info</span>
          <span>Click to see account information</span>
        </div>
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="navigateTo('/wallet')"
        >
          <span class="font-bold text-2xl">Wallet</span>
          <span>Click to see wallet info</span>
        </div>
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="navigateTo('/transactions')"
        >
          <span class="font-bold text-2xl">Transactions</span>
          <span>Click to see transactions</span>
        </div>
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="navigateTo('/orders')"
        >
          <span class="font-bold text-2xl">Orders</span>
          <span>Click to see past orders</span>
        </div>
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="navigateTo('/refferal')"
        >
          <span class="font-bold text-2xl">Refferal</span>
          <span>Click to see Refferal info</span>
        </div>
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="navigateTo('/support')"
        >
          <span class="font-bold text-2xl">Support</span>
          <span>Click to contact support</span>
        </div>
      </div>
    </div>
    <SetShippingAddress />
    <SetLanguage />
  </div>
  <NavComponent />
</template>
<script lang="ts">
import { Options, Vue } from "vue-class-component";
import NavComponent from "./NavComponent.vue";
import SetShippingAddress from "./SetShippingAddress.vue";
import { useSAStore, useLanguageStore } from "../store/dstore";
import SetLanguage from "./SetLanguage.vue";
import { Router } from "vue-router";

@Options({
  components: {
    NavComponent,
    SetShippingAddress,
    SetLanguage,
  },
})
export default class SettingsView extends Vue {
  sa = useSAStore();
  sl = useLanguageStore();
  router!: Router;

  created() {
    this.router = this.$router;
  }

  navigateTo(path: string) {
    this.router.push(path);
  }
}
</script>
