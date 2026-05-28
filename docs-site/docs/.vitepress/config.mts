import { defineConfig } from "vitepress";

export default defineConfig({
  base: "/lyshop/",
  lang: "zh-CN",
  title: "LYShop 零云商城",
  description: "LYShop 零云商城 - 开源插件化商城系统，基于 Go + Vue3 + uni-app，支持多端（PC/H5/小程序/App），插件化架构，AI生图，IM客服，Docker一键部署，免费开源电商解决方案",

  head: [
    ["link", { rel: "icon", href: "/lyshop/lyshop-mark.svg", type: "image/svg+xml" }],
    ["meta", { name: "keywords", content: "lyshop,零云商城,开源商城,Go商城,Vue3商城,uni-app商城,插件化商城,多端商城,小程序商城,H5商城,电商系统,开源电商,免费商城系统,商城源码,Go电商,微信商城,商城开源项目" }],
    ["meta", { name: "author", content: "ijry" }],
    ["meta", { property: "og:title", content: "LYShop 零云商城 - 开源插件化多端商城系统" }],
    ["meta", { property: "og:description", content: "基于 Go + Vue3 + uni-app 的开源插件化商城，支持PC/H5/小程序/App多端，AI生图，IM客服，Docker一键部署" }],
    ["meta", { property: "og:type", content: "website" }],
    ["meta", { property: "og:url", content: "https://ijry.github.io/lyshop/" }],
  ],

  sitemap: {
    hostname: "https://ijry.github.io/lyshop",
  },

  themeConfig: {
    logo: "/lyshop-mark.svg",
    nav: [
      { text: "首页", link: "/" },
      { text: "功能介绍", link: "/guide/features" },
      { text: "部署文档", link: "/deploy/" },
      { text: "接口文档", link: "/api/" },
      { text: "二次开发", link: "/dev/secondary-development" },
      {
        text: "在线演示",
        items: [
          { text: "PC 商城演示", link: "/web-demo/index.html", target: "_blank" },
          { text: "管理后台演示（admin/admin123）", link: "/admin-demo/index.html", target: "_blank" },
          { text: "H5 移动端演示（右下角浮窗）", link: "/" },
          { text: "商家 eapp 演示", link: "/eapp-demo/index.html#/pages/dashboard/index", target: "_blank" }
        ]
      },
      { text: "AI花絮", link: "/ai-notes/" }
    ],
    sidebar: {
      "/guide/": [
        {
          text: "产品指南",
          items: [
            { text: "功能介绍", link: "/guide/features" },
            { text: "商家移动端 eapp", link: "/guide/eapp-merchant" }
          ]
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
            { text: "后台概览接口", link: "/api/admin" },
            { text: "商品接口", link: "/api/product" },
            { text: "订单接口", link: "/api/order" },
            { text: "仓储接口", link: "/api/wms" },
            { text: "装修接口", link: "/api/decor" },
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
      ],
      "/ai-notes/": [
        {
          text: "AI花絮",
          items: [
            { text: "导读", link: "/ai-notes/" },
            { text: "能力总览", link: "/ai-notes/capabilities" },
            { text: "接口速查", link: "/ai-notes/api" },
            { text: "演进记录", link: "/ai-notes/evolution" }
          ]
        }
      ]
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/ijry/lyshop" }
    ],
    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026 LYShop"
    }
  }
});
