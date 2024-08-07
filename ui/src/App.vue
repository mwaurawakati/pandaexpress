<template>
  <notifications position="middle right" />
  <router-view />
</template>
<script setup>
import { provide } from "vue";
import { useDialog } from "./store/state";

const dialog = useDialog();

provide("dialog", dialog);
</script>
<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  background-color: white;
}

nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
body,
html {
  margin: 0;
  padding: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: white;
}
.active {
  background-color: #42b983;
}
</style>
