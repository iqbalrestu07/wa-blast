package flags

var AppName = "COBRA Core Service"
var AppVersion = "1.0.0"
var AppCommitHash = "N/A"

const (
	// Prefix for environment variables
	EnvPrefix = "COBRA"

	// Content Type Headers
	HeaderKeyContentType        = "Content-Type"
	HeaderKeyCOBRAAuthorization = "Authorization"
	HeaderKeyCOBRAAccessToken   = "X-COBRA-Access-Token"
	HeaderKeyCOBRATokenExpired  = "X-COBRA-Token-Expired"
	HeaderKeyCOBRASubject       = "X-COBRA-Subject"

	// Content Type Value
	ContentTypeJSON = "application/json; charset=utf-8"
	ContentTypeXML  = "application/xml; charset=utf-8"
	ContentTypeHTML = "text/html; charset=utf-8"

	// ACL
	ACLAuthenticatedAdmin     = ""
	ACLAuthenticatedUser      = "1"
	ACLAuthenticatedAnonymous = "2"
	ACLEveryone               = "3"

	//status
	UserStatusActive = 1
	UserStatusBanned = 0

	//Roles User
	RoleUser    = 1
	RoleSupport = 5
	RoleVip     = 7
	RoleVvip    = 8
	RoleMaster  = 9

	// Main Product
	RatingTahunan = "1"
	PaketRating   = "2"

	//Filter....
	FilterName        = "filter_product_name"
	FilterTahun       = "filter_tahun"
	FilterSubCategory = "filter_category"
	FilterStatus      = "filter_status"

	QueueAcceptOrder       = "queue.accept.order"
	QueueNotificationOrder = "queue.notification.order"
	QueueSyncStatusOrder   = "queue.notification.order"

	QueueAcceptTopup = "queue.proses.topup"

	//Flow Order
	OrderPlaced     = 1
	OrderAccept     = 2
	OrderOnProgress = 3
	OrderSuccess    = 10
	OrderCancel     = 11
	OrderReject     = 99

	//Type Deposit
	DepositTopup       = "topup"
	DepositFee         = "admin-fee"
	DepositOrder       = "order-item"
	DepositTransferIn  = "in-transfer"
	DepositTransferOut = "out-transfer"

	//Topup
	TopupPlaced = 1
	TopupAccept = 10
	TopupCancel = 11
	TopupReject = 99

	// Status
	StatusWaitingPayment = "Waiting Payment"
	StatusUserConfirmed  = "User Confirmed"
	StatusCompeled       = "Completed"
	StatusCanceled       = "Canceled"
	StatusRefund         = "Refund"
)
