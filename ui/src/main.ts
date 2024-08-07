import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import "./index.css";
import { createPinia } from "pinia";
/* import the fontawesome core */
import { library } from "@fortawesome/fontawesome-svg-core";

/* import font awesome icon component */
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

/* import specific icons */
import {
  faUserSecret,
  faCog,
  faCartPlus,
  faShoppingCart,
  faStore,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
library.add(faUserSecret, faCog, faCartPlus, faShoppingCart, faStore, faUser);
import Notifications from "@kyvg/vue3-notification";
const app = createApp(App)
  .component("font-awesome-icon", FontAwesomeIcon)
  .use(store)
  .use(router);
const pinia = createPinia();
app.use(pinia);
app.use(Notifications);
app.mount("#app");
