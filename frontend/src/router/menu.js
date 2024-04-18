export default [
  {
    path: "/ipfs",
    name: "ipfs",
    component: () => import("@/views/Home.vue"),
    meta: {
      title: "基于纠删码的IPFS存储优化",
    }
  },
];
