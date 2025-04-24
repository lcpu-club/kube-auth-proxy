import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("@/views/HomeView.vue"),
    },
    {
      path: "/auth/token/:token",
      name: "auth_token",
      component: () => import("@/views/AuthTokenView.vue"),
    },
    {
      path: "/tokens",
      name: "token",
      component: () => import("@/views/TokenView.vue"),
    },
    {
      path: "/logout",
      name: "logout",
      component: () => import("@/views/LogoutView.vue"),
    },
    {
      path: "/jobs",
      name: "jobs",
      component: () => import("@/views/JobsView.vue"),
    },
    {
      path: "/pods",
      name: "pods",
      component: () => import("@/views/PodsView.vue"),
    },
    {
      path: "/webterm",
      name: "webterm",
      component: () => import("@/views/WebTermView.vue"),
    },
    {
      path: "/storage",
      name: "storage",
      component: () => import("@/views/StorageView.vue"),
    },
    {
      path: "/config/configmap",
      name: "configmap",
      component: () => import("@/views/ConfigMapView.vue"),
    },
    {
      path: "/config/secret",
      name: "secret",
      component: () => import("@/views/SecretView.vue"),
    },
    {
      path: "/config/sshauthkey",
      name: "sshauthkey",
      component: () => import("@/views/SSHAuthKeyView.vue"),
    },
    {
      path: "/config/sshkeypair",
      name: "sshkeypair",
      component: () => import("@/views/SSHKeyPairView.vue"),
    }
  ],
});

export default router;
