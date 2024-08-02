<template>
  <div class="nav">
    <div class="flex-item">
      <button :class="{ active: isActive('/') }" @click="navigateTo('/')">
        <span><font-awesome-icon :icon="['fas', 'store']" /></span>
        <span>Shop</span>
      </button>
    </div>
    <div class="flex-item">
      <button
        :class="{ active: isActive('/cart') }"
        @click="navigateTo('/cart')"
      >
        <span><font-awesome-icon :icon="['fas', 'shopping-cart']" /></span>
        <span>Cart</span>
      </button>
    </div>
    <div class="flex-item">
      <button
        :class="{ active: isActive('/account') }"
        @click="navigateTo('/account')"
      >
        <span><font-awesome-icon :icon="['fas', 'user']" /></span>
        <span>Account</span>
      </button>
    </div>
    <div class="flex-item">
      <button
        :class="{ active: isActive('/settings') }"
        @click="navigateTo('/settings')"
      >
        <span><font-awesome-icon :icon="['fas', 'cog']" /></span>
        <span>Settings</span>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    isActive(path) {
      return this.$route.path === path;
    },
    navigateTo(path) {
      this.$router.push(path);
    },
  },
};
</script>

<style scoped>
.nav {
  width: 100%;
  height: 50px;
  bottom: 5px;
  right: 10px;
  position: fixed;
  border-top: black 1px solid;
  display: flex;
  background-color: whitesmoke;
}

.flex-item {
  flex: 1;
  /* Each item takes equal space */
  text-align: center;
  /* Center text inside items (optional) */
  border: 1px solid #000;
  /* Optional: visual boundary for items */
  box-sizing: border-box;
  /* Include padding and border in width/height */
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-item button {
  width: 100%;
  height: 100%;
  border-radius: 0px;
  display: flex;
  flex-direction: column;
}

.flex-item button {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
