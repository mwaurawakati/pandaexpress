<template>
  <div class="container">
    <div>Account</div>
    <div>Balances</div>
    <div>Transactions</div>
    <div>Orders</div>
    <div>Refferals</div>
    <NavComponent />
  </div>
</template>
<script lang="ts">
import { Options, Vue } from "vue-class-component";
import NavComponent from "./NavComponent.vue";

@Options({
  components: {
    NavComponent,
  },
})
export default class AccountView extends Vue {}
</script>
