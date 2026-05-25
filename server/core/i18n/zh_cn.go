package i18n

var zhCN = map[string]string{
	// Auth / Middleware
	"auth.loginRequired":       "请先登录",
	"auth.tokenInvalid":        "Token无效或已过期",
	"auth.noPermission":        "无权限",
	"auth.noPermissionDetail":  "无权限: %s",
	"auth.codeExpired":         "验证码已过期",
	"auth.codeInvalid":         "验证码错误",
	"auth.accountDisabled":     "账号已被禁用",
	"auth.userNotFound":        "用户不存在",
	"auth.phoneMismatch":       "手机号与账号不匹配",
	"auth.accountDeleted":      "已注销用户",
	"auth.credentialsInvalid":  "用户名或密码错误",
	"auth.phoneRequired":       "手机号不能为空",
	"auth.sendFailed":          "发送失败: %s",

	// Order Status Labels
	"order.status.pending":    "待付款",
	"order.status.shipped":    "待发货",
	"order.status.delivering": "待收货",
	"order.status.completed":  "已完成",
	"order.status.afterSale":  "售后",

	// After Sale Case Type Labels
	"afterSale.type.return":   "退货",
	"afterSale.type.exchange": "换货",

	// After Sale Status Labels
	"afterSale.status.applied":       "已申请",
	"afterSale.status.waitReturn":    "待用户回寄",
	"afterSale.status.returning":     "用户回寄中",
	"afterSale.status.received":      "仓库已收货",
	"afterSale.status.refundPending": "待退款",
	"afterSale.status.refunded":      "已退款",
	"afterSale.status.reshipPending": "待补发",
	"afterSale.status.reshipped":     "已补发",
	"afterSale.status.completed":     "已完结",
	"afterSale.status.rejected":      "已拒绝",
	"afterSale.status.closed":        "已关闭",

	// Delivery Type Labels
	"delivery.type.express": "快递配送",
	"delivery.type.local":   "同城配送",

	// Shipment Status Labels
	"shipment.status.pending":   "待揽收",
	"shipment.status.shipped":   "已发货",
	"shipment.status.inTransit": "运输中",
	"shipment.status.signed":    "已签收",
	"shipment.status.exception": "物流异常",

	// Shipment Biz Type Labels
	"shipment.bizType.initial": "首发",
	"shipment.bizType.reship":  "补发",
	"shipment.bizType.return":  "回寄",

	// Shipment Direction Labels
	"shipment.direction.outbound": "寄出",
	"shipment.direction.inbound":  "回寄",

	// After Sale Action Labels
	"afterSale.action.apply":   "提交申请",
	"afterSale.action.review":  "审核",
	"afterSale.action.return":  "回寄物流",
	"afterSale.action.receive": "确认收货",
	"afterSale.action.refund":  "退款",
	"afterSale.action.reship":  "补发",
	"afterSale.action.close":   "关闭",
	"afterSale.action.complete":"完结",

	// Order Error Messages
	"order.err.addressNotFound":   "收货地址不存在",
	"order.err.cartEmpty":         "购物车为空",
	"order.err.noValidProducts":   "未选择有效商品",
	"order.err.priceCalcFailed":   "价格计算失败",
	"order.err.insufficientStock": "商品库存不足",
	"order.err.unsupportedDelivery":"不支持的发货类型",
	"order.err.notFound":          "订单不存在",
	"order.err.cannotShip":        "当前状态不可发货",
	"order.err.reshipNeedCase":    "补发需关联售后单",
	"order.err.reshipCaseMismatch":"补发售后单与订单不匹配",
	"order.err.reshipStatusInvalid":"当前售后状态不可补发",
	"order.err.cannotPay":         "订单不存在或当前状态不可支付",

	// After Sale Error Messages
	"afterSale.err.invalidParams":     "参数错误",
	"afterSale.err.reasonRequired":    "请填写售后原因",
	"afterSale.err.selectProducts":    "请选择售后商品",
	"afterSale.err.orderNotFound":     "订单不存在",
	"afterSale.err.orderNotPaid":      "订单未支付，暂不可售后",
	"afterSale.err.selectValidProducts":"请选择有效售后商品",
	"afterSale.err.itemNotFound":      "订单商品不存在",
	"afterSale.err.activeCase":        "已有进行中售后",
	"afterSale.err.qtyExceeded":       "售后数量超过已购数量",
	"afterSale.err.cannotReview":      "当前状态不可审核",
	"afterSale.err.cannotReturn":      "当前状态不可提交回寄物流",
	"afterSale.err.cannotReceive":     "当前状态不可确认收货",
	"afterSale.err.cannotRefund":      "当前状态不可登记退款",
	"afterSale.err.invalidRefundAmount":"退款金额无效",
	"afterSale.err.cannotComplete":    "当前状态不可完结",
	"afterSale.err.cannotClose":       "当前状态不可关闭",
	"afterSale.err.closeReasonRequired":"请填写关闭原因",
	"afterSale.err.notFound":          "售后单不存在",

	// After Sale Log Messages
	"afterSale.log.applied":       "提交售后申请",
	"afterSale.log.reviewed":      "售后审核",
	"afterSale.log.returnShipped": "用户提交回寄物流",
	"afterSale.log.received":      "仓库收货确认",
	"afterSale.log.refunded":      "退款登记",
	"afterSale.log.completed":     "售后完结",
	"afterSale.log.closed":        "关闭售后",
	"afterSale.log.reship":        "售后补发",

	// Delivery Validation
	"delivery.err.expressRequired":    "请选择快递公司",
	"delivery.err.trackingRequired":   "请填写快递单号",
	"delivery.err.riderNameRequired":  "请填写骑手名称",
	"delivery.err.riderPhoneRequired": "请填写骑手电话",

	// Upload
	"upload.err.fileRequired": "请选择文件",
}
