import request from './request'

// ---- Product ----
export const getCategories = () => request.get('/categories')
export const createCategory = (data: any) => request.post('/categories', data)
export const deleteCategory = (id: number) => request.delete(`/categories/${id}`)

export const getProducts = (params?: any) => request.get('/products', { params })
export const getProduct = (id: number) => request.get(`/products/${id}`)
export const createProduct = (data: any) => request.post('/products', data)
export const updateProduct = (id: number, data: any) => request.put(`/products/${id}`, data)
export const deleteProduct = (id: number) => request.delete(`/products/${id}`)

// ---- Order ----
export const getOrders = (params?: any) => request.get('/orders', { params })
export const shipOrder = (id: number, trackingNo: string) =>
  request.put(`/orders/${id}/ship`, { tracking_no: trackingNo })

// ---- WMS ----
export const getWarehouses = () => request.get('/wms/warehouses')
export const createWarehouse = (data: any) => request.post('/wms/warehouses', data)
export const getStocks = (params?: any) => request.get('/wms/stocks', { params })
export const doInbound = (data: any) => request.post('/wms/inbound', data)
export const doOutbound = (data: any) => request.post('/wms/outbound', data)
