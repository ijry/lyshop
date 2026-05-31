import request from './request'

// ---- Dashboard ----
export const getDashboard = (params?: any) => request.get('/dashboard', { params })

// ---- Product ----
export const getCategories = (params?: any) => request.get<any[]>('/categories', { params })
export const createCategory = (data: any) => request.post('/categories', data)
export const updateCategory = (id: number, data: any) => request.put(`/categories/${id}`, data)
export const deleteCategory = (id: number) => request.delete(`/categories/${id}`)

export const getProducts = (params?: any) => request.get('/products', { params })
export const getProduct = (id: number) => request.get(`/products/${id}`)
export const createProduct = (data: any) => request.post('/products', data)
export const updateProduct = (id: number, data: any) => request.put(`/products/${id}`, data)
export const deleteProduct = (id: number) => request.delete(`/products/${id}`)
export const uploadFile = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  return request.post('/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
export const getAiModels = () => request.get<any[]>('/ai/models')
export const generateAiImage = (data: any) => request.post('/ai/generate', data)
export const getAiTask = (id: number) => request.get(`/ai/tasks/${id}`)
export const getVipLevels = (params?: any) => request.get('/vip/levels', { params })
export const getVipSkuPrices = (params?: any) => request.get('/vip/sku-prices', { params })
export const createVipSkuPrice = (data: any) => request.post('/vip/sku-prices', data)
export const updateVipSkuPrice = (id: number, data: any) => request.put(`/vip/sku-prices/${id}`, data)
export const deleteVipSkuPrice = (id: number) => request.delete(`/vip/sku-prices/${id}`)

// ---- Order ----
export const getOrders = (params?: any) => request.get('/orders', { params })
export const getOrderDetail = (id: number) => request.get(`/orders/${id}`)
export const shipOrder = (id: number, data: any) =>
  request.put(`/orders/${id}/ship`, data)
export const syncShipment = (orderID: number, shipmentID: number) =>
  request.post(`/orders/${orderID}/shipments/${shipmentID}/sync`)
export const getShipmentTracks = (orderID: number, shipmentID: number) =>
  request.get(`/orders/${orderID}/shipments/${shipmentID}/tracks`)
export const getAfterSales = (params?: any) => request.get('/after-sales', { params })
export const getAfterSaleDetail = (id: number) => request.get(`/after-sales/${id}`)
export const auditAfterSale = (id: number, data: any) => request.post(`/after-sales/${id}/audit`, data)
export const receiveAfterSale = (id: number) => request.post(`/after-sales/${id}/receive`)
export const refundAfterSale = (id: number, data: any) => request.post(`/after-sales/${id}/refund`, data)
export const completeAfterSale = (id: number) => request.post(`/after-sales/${id}/complete`)
export const closeAfterSale = (id: number, data: any) => request.post(`/after-sales/${id}/close`, data)
export const getReviews = (params?: any) => request.get('/reviews', { params })
export const getReviewDetail = (id: number) => request.get(`/reviews/${id}`)
export const replyReview = (id: number, content: string) => request.post(`/reviews/${id}/reply`, { content })

// ---- Delivery ----
export const getDeliveryMode = () => request.get<{ mode: string }>('/delivery/mode')

// ---- WMS ----
export const getWarehouses = () => request.get('/wms/warehouses')
export const createWarehouse = (data: any) => request.post('/wms/warehouses', data)
export const getStocks = (params?: any) => request.get('/wms/stocks', { params })
