package model

/*
嵌套+外键  ProductImage，ProductSize，ProductSeo

字段ProductSku可以保证数据的幂等性。
SPU：是标准化产品单元，区分品种
SKU是库存量单位，区分单品
商品：特指与商家有关的商品，可对应多个SKU。
如流光蓝（三种颜色：流光蓝、霓光紫、霓光渐变色）+8G+128G（两种配置：8G+128G、6G+128G）。
即Oppo R17有一个SPU、6种SKU（3种颜色*2种配置）。
*/
type Product struct {
	ID                 int64          `gorm:"primary_key;not_null;auto_increment" json:"id"`
	ProductName        string         `json:"product_name"`
	ProductSku         string         `gorm:"unique_index:not_null" json:"product_sku"`
	ProductPrice       float64        `json:"product_price"`
	ProductDescription string         `json:"product_description"`
	ProductImage       []ProductImage `gorm:"ForeignKey:ImageProductID" json:"product_image"`
	ProductSize        []ProductSize  `gorm:"ForeignKey:SizeProductID" json:"product_size"`
	ProductSeo         ProductSeo     `gorm:"ForeignKey:SeoProductID" json:"product_seo"`
}
