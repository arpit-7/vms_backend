package models



type CustomMap struct {
	
	ID          uint    `gorm:"primarykey" json:"id"`
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Type        string `gorm:"column:type;type:varchar(50);not null" json:"type"` // image/raster/vector
	ImageData   string `gorm:"column:image_data;type:longtext" json:"imageData,omitempty"`
	ImageWidth  int    `gorm:"column:image_width" json:"imageWidth"`
	ImageHeight int    `gorm:"column:image_height" json:"imageHeight"`
	Available   bool   `gorm:"column:available;default:true" json:"available"`
	TileURL     string `gorm:"column:tile_url;type:text" json:"tileUrlPattern,omitempty"`
	StyleURL    string `gorm:"column:style_url;type:text" json:"styleUrl,omitempty"`
	BoundsJSON  string `gorm:"column:bounds;type:json" json:"bounds,omitempty"`
}

type CameraPosition struct {
	
	ID          uint   `gorm:"primarykey" json:"id"`
	CustomMapID uint   `gorm:"column:custom_map_id;not null;index" json:"customMapId"`
	CameraID    string `gorm:"column:camera_id;type:varchar(255);not null" json:"cameraId"`
	CameraName  string `gorm:"column:camera_name;type:varchar(255)" json:"cameraName"`
	X           int    `gorm:"column:x" json:"x"`
	Y           int    `gorm:"column:y" json:"y"`
	Bearing     int    `gorm:"column:bearing" json:"bearing"`
	FOV         int    `gorm:"column:fov" json:"fov"`
	Range       int    `gorm:"column:range" json:"range"`
	
	CustomMap CustomMap `gorm:"foreignKey:CustomMapID;references:ID" json:"-"`
}