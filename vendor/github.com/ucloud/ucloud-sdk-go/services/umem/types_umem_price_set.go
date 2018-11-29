package umem

/*
UMemPriceSet - DescribeUMemPrice

this model is auto created by ucloud code generater for open api,
you can also see https://docs.ucloud.cn for detail.
*/
type UMemPriceSet struct {

	// Year， Month， Dynamic，Trial
	ChargeType string

	// 价格，单位: 元，保留小数点后两位有效数字
	Price float64
}
