# lyshop 2.0 Phase 2: Product / Order / WMS Plugins

**Goal:** Implement the `product`, `order`, and `wms` plugins — models, migrations, services, admin APIs, frontend APIs, admin Vue3 pages, and uni-app pages.

**Architecture:** Each plugin lives in `server/plugins/<name>/`, implements `core/plugin.Plugin`, self-registers via `init()`, and is blank-imported in `main.go`. GORM AutoMigrate handles schema creation.

**Tech Stack:** Go 1.22 · Gin · GORM · Redis (cart) · Vue3 + TailwindCSS (admin pages) · uni-app + uview-plus 3.x (frontend pages)

**Spec:** `docs/superpowers/specs/2026-05-21-lyshop2-design.md` §6.1–6.5

---

## File Map

```
server/plugins/
├── product/
│   ├── plugin.json
│   ├── plugin.go          # implements Plugin interface, AutoMigrate
│   ├── model/
│   │   ├── category.go
│   │   ├── product.go
│   │   ├── sku.go
│   │   └── image.go
│   ├── service/
│   │   ├── category.go
│   │   └── product.go
│   └── api/
│       ├── front.go       # GET /api/v1/products, /products/:id, /categories
│       └── admin.go       # CRUD /admin/api/products, /categories
│
├── order/
│   ├── plugin.json
│   ├── plugin.go
│   ├── model/
│   │   ├── address.go
│   │   ├── cart.go        # GORM model (MySQL fallback)
│   │   ├── order.go
│   │   ├── order_item.go
│   │   ├── order_payment.go
│   │   └── order_refund.go
│   ├── service/
│   │   ├── cart.go        # Redis-based cart
│   │   └── order.go       # create, pay, ship, status transitions
│   └── api/
│       ├── front.go       # cart, create order, pay, my orders
│       └── admin.go       # order list, ship, refund
│
└── wms/
    ├── plugin.json
    ├── plugin.go
    ├── model/
    │   ├── warehouse.go
    │   ├── stock.go
    │   ├── inbound.go
    │   ├── outbound.go
    │   └── stock_log.go
    ├── service/
    │   └── stock.go       # inbound, outbound, adjust, query
    └── api/
        └── admin.go       # warehouses, stocks, inbound, outbound

admin/src/
└── views/
    ├── product/
    │   ├── ProductList.vue
    │   ├── ProductForm.vue
    │   └── CategoryList.vue
    ├── order/
    │   └── OrderList.vue
    └── wms/
        └── StockList.vue

app/pages/
├── product/
│   ├── list.vue
│   └── detail.vue
├── cart/
│   └── index.vue
└── order/
    ├── confirm.vue
    └── list.vue
```

---

## Task P1: product plugin — skeleton + models + migration

- [ ] Create `server/plugins/product/plugin.json`
- [ ] Create `server/plugins/product/model/category.go`
- [ ] Create `server/plugins/product/model/product.go`
- [ ] Create `server/plugins/product/model/sku.go`
- [ ] Create `server/plugins/product/model/image.go`
- [ ] Create `server/plugins/product/plugin.go`
- [ ] Blank-import in `server/main.go`
- [ ] `go build ./...` passes
- [ ] Commit

## Task P2: product service + APIs

- [ ] Create `server/plugins/product/service/category.go`
- [ ] Create `server/plugins/product/service/product.go`
- [ ] Create `server/plugins/product/api/front.go`
- [ ] Create `server/plugins/product/api/admin.go`
- [ ] `go build ./...` passes
- [ ] Commit

## Task P3: order plugin — skeleton + models

- [ ] Create `server/plugins/order/plugin.json`
- [ ] Create order models (address, cart, order, items, payment, refund)
- [ ] Create `server/plugins/order/plugin.go`
- [ ] Blank-import in `server/main.go`
- [ ] Commit

## Task P4: order service (cart + order)

- [ ] Create `server/plugins/order/service/cart.go`
- [ ] Create `server/plugins/order/service/order.go`
- [ ] Commit

## Task P5: order APIs

- [ ] Create `server/plugins/order/api/front.go`
- [ ] Create `server/plugins/order/api/admin.go`
- [ ] `go build ./...` passes
- [ ] Commit

## Task P6: wms plugin

- [ ] Create `server/plugins/wms/plugin.json`
- [ ] Create wms models
- [ ] Create `server/plugins/wms/plugin.go`
- [ ] Create `server/plugins/wms/service/stock.go`
- [ ] Create `server/plugins/wms/api/admin.go`
- [ ] Blank-import in `server/main.go`
- [ ] `go build ./...` passes
- [ ] Commit

## Task P7: admin Vue3 pages

- [ ] Add routes to `admin/src/router/index.ts`
- [ ] Create `ProductList.vue`, `ProductForm.vue`, `CategoryList.vue`
- [ ] Create `OrderList.vue`
- [ ] Create `StockList.vue`
- [ ] `npm run build` passes
- [ ] Commit

## Task P8: uni-app pages

- [ ] Update `app/pages.json`
- [ ] Create product list + detail pages
- [ ] Create cart page
- [ ] Create order confirm + list pages
- [ ] Commit + push
