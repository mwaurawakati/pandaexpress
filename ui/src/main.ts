import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import "./index.css";
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

createApp(App)
  .component("font-awesome-icon", FontAwesomeIcon)
  .use(store)
  .use(router)
  .mount("#app");
