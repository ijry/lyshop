package model

import (
	"github.com/ijry/lyshop/core/model"
)

// DistributionConfig 分销配置
type DistributionConfig struct {
	model.Base
	Enabled          bool    `gorm:"not null;default:true" json:"enabled"`                    // 是否启用
	Level            int     `gorm:"not null;default:2" json:"level"`                         // 分销层级（1-3）
	Level1Rate       float64 `gorm:"type:decimal(5,2);not null;default:10" json:"level1_rate"` // 一级佣金比例
	Level2Rate       float64 `gorm:"type:decimal(5,2);not null;default:5" json:"level2_rate"`  // 二级佣金比例
	Level3Rate       float64 `gorm:"type:decimal(5,2);not null;default:2" json:"level3_rate"`  // 三级佣金比例
	MinWithdraw      float64 `gorm:"type:decimal(10,2);not null;default:100" json:"min_withdraw"` // 最低提现金额
	WithdrawFeeRate  float64 `gorm:"type:decimal(5,2);not null;default:0" json:"withdraw_fee_rate"` // 提现手续费比例
	AutoApprove      bool    `gorm:"not null;default:false" json:"auto_approve"`              // 自动审核分销商
	RequireRealName  bool    `gorm:"not null;default:true" json:"require_real_name"`          // 要求实名认证
}

func (DistributionConfig) TableName() string {
	return "distribution_configs"
}

// Distributor 分销商
type Distributor struct {
	model.Base
	UserID          uint64  `gorm:"not null;uniqueIndex" json:"user_id"`
	ParentID        uint64  `gorm:"default:0;index" json:"parent_id"`                    // 上级分销商ID
	Level           int     `gorm:"not null;default:1" json:"level"`                     // 分销商等级
	Status          string  `gorm:"size:32;not null;default:'pending';index" json:"status"` // pending|active|disabled
	TotalEarnings   float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total_earnings"` // 累计收益
	AvailableAmount float64 `gorm:"type:decimal(10,2);not null;default:0" json:"available_amount"` // 可提现金额
	FrozenAmount    float64 `gorm:"type:decimal(10,2);not null;default:0" json:"frozen_amount"` // 冻结金额
	WithdrawnAmount float64 `gorm:"type:decimal(10,2);not null;default:0" json:"withdrawn_amount"` // 已提现金额
	TotalCustomers  int     `gorm:"not null;default:0" json:"total_customers"`           // 累计客户数
	TotalOrders     int     `gorm:"not null;default:0" json:"total_orders"`              // 累计订单数
	RealName        string  `gorm:"size:64" json:"real_name"`
	IDCard          string  `gorm:"size:32" json:"id_card"`
	Phone           string  `gorm:"size:32" json:"phone"`
}

func (Distributor) TableName() string {
	return "distributors"
}

// DistributionOrder 分销订单
type DistributionOrder struct {
	model.Base
	OrderID        uint64  `gorm:"not null;index" json:"order_id"`
	DistributorID  uint64  `gorm:"not null;index" json:"distributor_id"`
	Level          int     `gorm:"not null" json:"level"`                                // 分销层级（1/2/3）
	OrderAmount    float64 `gorm:"type:decimal(10,2);not null" json:"order_amount"`      // 订单金额
	CommissionRate float64 `gorm:"type:decimal(5,2);not null" json:"commission_rate"`    // 佣金比例
	Commission     float64 `gorm:"type:decimal(10,2);not null" json:"commission"`        // 佣金金额
	Status         string  `gorm:"size:32;not null;default:'pending';index" json:"status"` // pending|settled|cancelled
	SettledAt      *model.Time `json:"settled_at"`                                       // 结算时间
}

func (DistributionOrder) TableName() string {
	return "distribution_orders"
}

// DistributionWithdrawal 提现申请
type DistributionWithdrawal struct {
	model.Base
	DistributorID uint64      `gorm:"not null;index" json:"distributor_id"`
	Amount        float64     `gorm:"type:decimal(10,2);not null" json:"amount"`              // 提现金额
	Fee           float64     `gorm:"type:decimal(10,2);not null;default:0" json:"fee"`       // 手续费
	ActualAmount  float64     `gorm:"type:decimal(10,2);not null" json:"actual_amount"`       // 实际到账
	Status        string      `gorm:"size:32;not null;default:'pending';index" json:"status"` // pending|approved|rejected|completed
	RejectReason  string      `gorm:"size:255" json:"reject_reason"`
	BankName      string      `gorm:"size:64" json:"bank_name"`
	BankAccount   string      `gorm:"size:64" json:"bank_account"`
	AccountName   string      `gorm:"size:64" json:"account_name"`
	ProcessedAt   *model.Time `json:"processed_at"`
	CompletedAt   *model.Time `json:"completed_at"`
}

func (DistributionWithdrawal) TableName() string {
	return "distribution_withdrawals"
}
