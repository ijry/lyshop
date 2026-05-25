import type { MockPreset } from './types'
import categories from '../data/categories.json'
import products from '../data/products.json'
import productDetail from '../data/product-detail.json'
import indexDecor from '../data/index-decor.json'
import seckills from '../data/seckills.json'
import groupBuy from '../data/group-buy.json'
import bargain from '../data/bargain.json'
import recommend from '../data/recommend.json'
import cart from '../data/cart.json'
import orders from '../data/orders.json'

export const mall: MockPreset = {
  key: 'mall',
  name: '综合商城',
  categories: categories as any,
  products: products as any,
  productDetail: productDetail as any,
  indexDecor: indexDecor as any,
  seckills: seckills as any,
  groupBuy: groupBuy as any,
  bargain: bargain as any,
  recommend: recommend as any,
  cart: cart as any,
  orders: orders as any,
}
