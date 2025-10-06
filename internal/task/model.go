package task

type Task struct {
	ID     uint32 `gorm:"primaryKey"`
	Title  string
	UserID uint32 `gorm:"index"`
}
