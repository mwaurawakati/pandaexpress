<template>
  <!--nav>
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </nav-->
  <router-view />
</template>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
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
}
.active {
  background-color: #42b983;
}
</style>
