package controllers

import "hublabs/coupon-api/models"

const DefaultMaxResultCount = 10

type SearchInput struct {
	Q      string               `json:"q" query:"q"`
	Ids    string               `json:"ids" query:"ids"`
	Fields models.FieldTypeList `json:"fields" query:"fields"`
	Filter string               `json:"filter" query:"filter" swagger:"enum(unused|used|expired|valid|invalid|available|unavailable|wait)"`
	Offset int                  `json:"offset",query:"skipCount"`
	Limit  int                  `json:"limit" query:"maxResultCount"`
	Order  []string             `json:"order",query:"order"`
	Sortby []string             `json:"sortby" query:"sortby"`
	Status string               `json:"status" query:"status"`
}

type FieldInput struct {
	Fields models.FieldTypeList `query:"fields"`
}

type CustSearchInput struct {
	Q        string   `json:"q" query:"q"`
	Nos      string   `json:"nos"`
	CustId   string   `json:"custId"`
	SaleType string   `json:"saleType"`
	Filter   string   `json:"filter" query:"filter" swagger:"enum(unused|used|expired|valid|invalid|available|unavailable|wait)"`
	Offset   int      `json:"offset",query:"skipCount"`
	Limit    int      `json:"limit" query:"maxResultCount"`
	Order    []string `json:"order",query:"order"`
	Sortby   []string `json:"sortby" query:"sortby"`
}

type SendInput struct {
	CustId   string          `json:"custId"`
	BrandId  int64           `json:"brandId"`
	Birthday string          `json:"birthday"`
	SendType models.SendType `json:"sendType`
}
