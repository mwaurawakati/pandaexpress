<template>
  <div class="container">
    cart
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
export default class CartView extends Vue {}
</script>
