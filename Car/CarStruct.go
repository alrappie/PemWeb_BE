package Car

type Car struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"uniqueIndex:model;size:50" json:"nama"`
	Year   string `gorm:"uniqueIndex:model;size:4" json:"tahun"`
	Engine string `gorm:"size:10" json:"mesin"`
	Price  int    `gorm:"size:20" json:"harga"`
	Stock  int    `gorm:"size:3" json:"jumlah"`
}
