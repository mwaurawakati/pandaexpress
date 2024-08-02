<template>
  <div class="container">
    <!--div class="nav">
      <div class="flex-item">
        <button class="active">
          <span><font-awesome-icon :icon="['fas', 'store']" /></span>
          <span>Shop</span>
        </button>
      </div>
      <div class="flex-item">
        <button>
          <span><font-awesome-icon :icon="['fas', 'shopping-cart']" /></span>
          <span>Cart</span>
        </button>
      </div>
      <div class="flex-item">
        <button>
          <span><font-awesome-icon :icon="['fas', 'user']" /></span>
          <span>Account</span>
        </button>
      </div>
      <div class="flex-item">
        <button href="/settings">
          <span><font-awesome-icon :icon="['fas', 'cog']" /></span>
          <span>Settings</span>
        </button>
      </div>
    </div-->
    <NavComponent />
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import HelloWorld from "@/components/HelloWorld.vue"; // @ is an alias to /src
import NavComponent from "./NavComponent.vue";

@Options({
  components: {
    HelloWorld,
    NavComponent,
  },
})
export default class HomeView extends Vue {}
</script>

<style>
.container {
  /*width: 100vw;
  height: 100vh;
  top: 0;
  bottom: 0;
  margin: 5;
  position: fixed;*/
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}
</style>
