export interface MockPreset {
  key: string
  name: string
  categories: any[]
  products: { list: any[]; total: number; page: number; size: number }
  productDetail: any
  indexDecor: { components: any[] }
  seckills: any
  groupBuy: any
  bargain: any
  recommend: any[]
  cart: any[]
  orders: { list: any[]; total: number; page: number; size: number }
}
