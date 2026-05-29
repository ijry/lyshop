package i18n

var en = map[string]string{
	// Auth / Middleware
	"auth.loginRequired":      "Please log in",
	"auth.tokenInvalid":       "Token invalid or expired",
	"auth.noPermission":       "No permission",
	"auth.noPermissionDetail": "No permission: %s",
	"auth.codeExpired":        "Verification code expired",
	"auth.codeInvalid":        "Invalid verification code",
	"auth.accountDisabled":    "Account disabled",
	"auth.userNotFound":       "User not found",
	"auth.phoneMismatch":      "Phone number does not match account",
	"auth.accountDeleted":     "Account has been deleted",
	"auth.credentialsInvalid": "Invalid username or password",
	"auth.phoneRequired":      "Phone number is required",
	"auth.sendFailed":         "Send failed: %s",

	// Order Status Labels
	"order.status.pending":    "Pending Payment",
	"order.status.shipped":    "Processing",
	"order.status.delivering": "Shipping",
	"order.status.completed":  "Completed",
	"order.status.afterSale":  "After-Sale",
	"order.status.canceled":   "Canceled",

	// After Sale Case Type Labels
	"afterSale.type.return":   "Return",
	"afterSale.type.exchange": "Exchange",

	// After Sale Status Labels
	"afterSale.status.applied":       "Applied",
	"afterSale.status.waitReturn":    "Awaiting Return",
	"afterSale.status.returning":     "Returning",
	"afterSale.status.received":      "Warehouse Received",
	"afterSale.status.refundPending": "Refund Pending",
	"afterSale.status.refunded":      "Refunded",
	"afterSale.status.reshipPending": "Reship Pending",
	"afterSale.status.reshipped":     "Reshipped",
	"afterSale.status.completed":     "Completed",
	"afterSale.status.rejected":      "Rejected",
	"afterSale.status.closed":        "Closed",

	// Delivery Type Labels
	"delivery.type.express": "Express",
	"delivery.type.local":   "Local Delivery",

	// Shipment Status Labels
	"shipment.status.pending":   "Pending Pickup",
	"shipment.status.shipped":   "Shipped",
	"shipment.status.inTransit": "In Transit",
	"shipment.status.signed":    "Delivered",
	"shipment.status.exception": "Exception",

	// Shipment Biz Type Labels
	"shipment.bizType.initial": "Initial",
	"shipment.bizType.reship":  "Reship",
	"shipment.bizType.return":  "Return",

	// Shipment Direction Labels
	"shipment.direction.outbound": "Outbound",
	"shipment.direction.inbound":  "Return",

	// After Sale Action Labels
	"afterSale.action.apply":    "Applied",
	"afterSale.action.review":   "Reviewed",
	"afterSale.action.return":   "Return Shipped",
	"afterSale.action.receive":  "Received",
	"afterSale.action.refund":   "Refunded",
	"afterSale.action.reship":   "Reshipped",
	"afterSale.action.close":    "Closed",
	"afterSale.action.complete": "Completed",

	// Order Error Messages
	"order.err.addressNotFound":     "Delivery address not found",
	"order.err.cartEmpty":           "Cart is empty",
	"order.err.noValidProducts":     "No valid products selected",
	"order.err.priceCalcFailed":     "Price calculation failed",
	"order.err.insufficientStock":   "Insufficient stock",
	"order.err.unsupportedDelivery": "Unsupported delivery type",
	"order.err.notFound":            "Order not found",
	"order.err.cannotShip":          "Cannot ship in current status",
	"order.err.reshipNeedCase":      "Reship requires an after-sale case",
	"order.err.reshipCaseMismatch":  "After-sale case does not match order",
	"order.err.reshipStatusInvalid": "After-sale status does not allow reship",
	"order.err.cannotPay":           "Order not found or cannot pay in current status",
	"order.err.cannotCancel":        "Order not found or cannot cancel in current status",

	// After Sale Error Messages
	"afterSale.err.invalidParams":       "Invalid parameters",
	"afterSale.err.reasonRequired":      "Please provide a reason",
	"afterSale.err.selectProducts":      "Please select products",
	"afterSale.err.orderNotFound":       "Order not found",
	"afterSale.err.orderNotPaid":        "Order not paid, after-sale not available",
	"afterSale.err.selectValidProducts": "Please select valid products",
	"afterSale.err.itemNotFound":        "Order item not found",
	"afterSale.err.activeCase":          "Active after-sale case exists",
	"afterSale.err.qtyExceeded":         "Quantity exceeds purchased amount",
	"afterSale.err.cannotReview":        "Cannot review in current status",
	"afterSale.err.cannotReturn":        "Cannot submit return shipping in current status",
	"afterSale.err.cannotReceive":       "Cannot confirm receipt in current status",
	"afterSale.err.cannotRefund":        "Cannot process refund in current status",
	"afterSale.err.invalidRefundAmount": "Invalid refund amount",
	"afterSale.err.cannotComplete":      "Cannot complete in current status",
	"afterSale.err.cannotClose":         "Cannot close in current status",
	"afterSale.err.closeReasonRequired": "Please provide a close reason",
	"afterSale.err.notFound":            "After-sale case not found",

	// After Sale Log Messages
	"afterSale.log.applied":       "After-sale application submitted",
	"afterSale.log.reviewed":      "After-sale reviewed",
	"afterSale.log.returnShipped": "Return shipment submitted",
	"afterSale.log.received":      "Warehouse receipt confirmed",
	"afterSale.log.refunded":      "Refund processed",
	"afterSale.log.completed":     "After-sale completed",
	"afterSale.log.closed":        "After-sale closed",
	"afterSale.log.reship":        "After-sale reship",

	// Delivery Validation
	"delivery.err.expressRequired":    "Please select an express company",
	"delivery.err.trackingRequired":   "Tracking number is required",
	"delivery.err.riderNameRequired":  "Rider name is required",
	"delivery.err.riderPhoneRequired": "Rider phone is required",

	// Upload
	"upload.err.fileRequired": "Please select a file",
}
