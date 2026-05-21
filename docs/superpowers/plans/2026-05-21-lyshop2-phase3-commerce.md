# lyshop 2.0 Phase 3: Commerce Plugins

**Goal:** Implement `marketing` (coupon/seckill/points), `wechat_pay`, `alipay`, `sms`, and `wechat_auth` plugins — complete with driver registrations, admin config pages, and frontend integration.

**Architecture:** Each plugin implements its respective driver interface and registers via `init()`. Payment drivers implement `core/driver/payment.Driver`. SMS implements `core/driver/sms.Driver`. OAuth implements `core/driver/oauth.Driver`.

**Spec:** `docs/superpowers/specs/2026-05-21-lyshop2-design.md` §6.3, §6.7–6.9

---

## File Map

```
server/plugins/
├── marketing/
│   ├── plugin.json / plugin.go
│   ├── model/  (coupon, coupon_user, activity, activity_product, points_log)
│   ├── service/ (coupon.go, activity.go, points.go)
│   └── api/    (front.go, admin.go)
├── wechat_pay/
│   ├── plugin.json / plugin.go  (registers payment.Driver "wechat")
│   ├── driver.go                (WechatPayDriver: CreateOrder/Refund/Notify)
│   └── config.go                (load AppID/MchID/Key from ConfigKV)
├── alipay/
│   ├── plugin.json / plugin.go  (registers payment.Driver "alipay")
│   ├── driver.go                (AlipayDriver)
│   └── config.go
├── sms/
│   ├── plugin.json / plugin.go  (registers sms.Driver)
│   ├── driver.go                (multi-provider: aliyun/tencent)
│   └── config.go
└── wechat_auth/
    ├── plugin.json / plugin.go  (registers oauth.Driver "wechat")
    ├── driver.go                (miniapp code2session + H5 OAuth2)
    └── config.go

admin/src/views/
├── marketing/  (CouponList.vue, ActivityList.vue)
└── system/     (PaymentConfig.vue, SmsConfig.vue)

app/pages/
└── marketing/  (coupon.vue, seckill.vue)
```
