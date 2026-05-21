import { defineConfig } from "vitepress";

export default defineConfig({
  lang: "zh-CN",
  title: "LYShop 文档中心",
  description: "LYShop 项目官网与技术文档",

  themeConfig: {
    logo: "/logo.svg",
    nav: [
      { text: "首页", link: "/" },
      { text: "功能介绍", link: "/guide/features" },
      { text: "部署文档", link: "/deploy/" },
      { text: "接口文档", link: "/api/" },
      { text: "二次开发", link: "/dev/secondary-development" }
    ],
    sidebar: {
      "/guide/": [
        {
          text: "产品指南",
          items: [{ text: "功能介绍", link: "/guide/features" }]
        }
      ],
      "/deploy/": [
        {
          text: "部署文档",
          items: [{ text: "部署总览", link: "/deploy/" }]
        }
      ],
      "/api/": [
        {
          text: "接口总览",
          items: [
            { text: "接口首页", link: "/api/" },
            { text: "认证接口", link: "/api/auth" },
            { text: "商品接口", link: "/api/product" },
            { text: "订单接口", link: "/api/order" },
            { text: "营销接口", link: "/api/marketing" },
            { text: "IM 接口", link: "/api/im" },
            { text: "支付接口", link: "/api/payment" }
          ]
        }
      ],
      "/dev/": [
        {
          text: "二次开发",
          items: [{ text: "扩展指南", link: "/dev/secondary-development" }]
        }
      ]
    },
    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026 LYShop"
    }
  }
});
