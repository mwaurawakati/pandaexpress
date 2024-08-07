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
          @click="sl.openDialog"
        >
          <span class="font-bold text-2xl">Language</span>
          <span>Set preffered language</span>
        </div>
        <div
          class="shadow-1xl h-16 m-2 rounded bg-slate-300 flex flex-col"
          @click="sa.openDialog"
        >
          <span class="font-bold text-2xl">Shipping Address</span>
          <span>Set shipping Address</span>
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
}
</script>
