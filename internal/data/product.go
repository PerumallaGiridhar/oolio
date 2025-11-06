package data

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}

type Product struct {
	ID       string  `json:"id"`
	Image    Image   `json:"image"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

// private in-memory list (simulate DB)
var products = []Product{
	{
		ID: "1",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
		},
		Name:     "Waffle with Berries",
		Category: "Waffle",
		Price:    6.5,
	},
	{
		ID: "2",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-creme-brulee-desktop.jpg",
		},
		Name:     "Vanilla Bean Crème Brûlée",
		Category: "Crème Brûlée",
		Price:    7.0,
	},
	{
		ID: "3",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-macaron-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-macaron-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-macaron-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-macaron-desktop.jpg",
		},
		Name:     "Macaron Mix of Five",
		Category: "Macaron",
		Price:    8.0,
	},
	{
		ID: "4",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-tiramisu-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-tiramisu-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-tiramisu-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-tiramisu-desktop.jpg",
		},
		Name:     "Classic Tiramisu",
		Category: "Tiramisu",
		Price:    5.5,
	},
	{
		ID: "5",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-baklava-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-baklava-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-baklava-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-baklava-desktop.jpg",
		},
		Name:     "Pistachio Baklava",
		Category: "Baklava",
		Price:    4.0,
	},
	{
		ID: "6",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-meringue-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-meringue-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-meringue-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-meringue-desktop.jpg",
		},
		Name:     "Lemon Meringue Pie",
		Category: "Pie",
		Price:    5.0,
	},
	{
		ID: "7",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-cake-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-cake-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-cake-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-cake-desktop.jpg",
		},
		Name:     "Red Velvet Cake",
		Category: "Cake",
		Price:    4.5,
	},
	{
		ID: "8",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-brownie-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-brownie-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-brownie-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-brownie-desktop.jpg",
		},
		Name:     "Salted Caramel Brownie",
		Category: "Brownie",
		Price:    4.5,
	},
	{
		ID: "9",
		Image: Image{
			Thumbnail: "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-thumbnail.jpg",
			Mobile:    "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-mobile.jpg",
			Tablet:    "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-tablet.jpg",
			Desktop:   "https://orderfoodonline.deno.dev/public/images/image-panna-cotta-desktop.jpg",
		},
		Name:     "Vanilla Panna Cotta",
		Category: "Panna Cotta",
		Price:    6.5,
	},
}

// exported functions — can later call DB instead
func GetAllProducts() []Product {
	return products
}

func GetProductByID(id string) (Product, bool) {
	for _, p := range products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}
