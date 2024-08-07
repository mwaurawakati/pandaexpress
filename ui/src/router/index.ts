import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import HomeView from "../views/HomeView.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "home",
    component: HomeView,
  },
  {
    path: "/about",
    name: "about",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/AboutView.vue"),
  },
  {
    path: "/cart",
    name: "cart",
    component: () =>
      import(/* webpackChunkName: "cart" */ "../views/CartComponent.vue"),
  },
  {
    path: "/account",
    name: "account",
    component: () =>
      import(/* webpackChunkName: "account" */ "../views/AccountComponent.vue"),
  },
  {
    path: "/orders",
    name: "orders",
    component: () =>
      import(/* webpackChunkName: "orders" */ "../views/OrdersComponent.vue"),
  },
  {
    path: "/transactions",
    name: "transactions",
    component: () =>
      import(
        /* webpackChunkName: "transactions" */ "../views/TransactionsComponent.vue"
      ),
  },
  {
    path: "/refferal",
    name: "refferal",
    component: () =>
      import(
        /* webpackChunkName: "refferal" */ "../views/RefferalComponent.vue"
      ),
  },
  {
    path: "/support",
    name: "support",
    component: () =>
      import(/* webpackChunkName: "support" */ "../views/SupportComponent.vue"),
  },
  {
    path: "/account_info",
    name: "account_info",
    component: () =>
      import(/* webpackChunkName: "account_info" */ "../views/AccountInfo.vue"),
  },
  {
    path: "/wallet",
    name: "wallet",
    component: () =>
      import(/* webpackChunkName: "wallet" */ "../views/WalletInfo.vue"),
  },
  {
    path: "/settings",
    name: "settings",
    component: () =>
      import(/* webpackChunkName: "settings" */ "../views/SettingsPage.vue"),
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("@/views/NotFound.vue"),
  },
  {
    path: "/orders/:id",
    name: "order-details",
    component: () =>
      import(
        /* webpackChunkName: "order-details" */ "../views/OrderDetails.vue"
      ),
    props: true, // This allows route params to be passed as props
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
