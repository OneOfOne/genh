package genh

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"go.oneofone.dev/genh/gsets"
)

type cloneStruct struct {
	Y map[string]any

	Ptr       *int
	PtrPtr    **int
	PtrPtrPtr ***int
	NilPtr    *int

	SA []any
	S  string
	X  []int
	A  [5]uint64
	x  int
	C  cloner
	C2 *cloner
	C3 cloner0

	SS simpleStruct
}

func (c cloneStruct) XV() int {
	return c.x
}

type (
	ifaceBug interface {
		XV() int
	}
	bugSlice []ifaceBug
	cloner   struct {
		A int
	}
)

func (c *cloner) Clone() *cloner {
	return c
}

type cloner0 struct {
	A      int
	cloned bool
}

func (c cloner0) Clone() cloner0 {
	c.cloned = true
	return c
}

type simpleStruct struct {
	A int
	B int64
	C string
	D bool
}

func TestBug01(t *testing.T) {
	s := bugSlice{cloneStruct{x: 42}, &cloneStruct{x: 42}}
	c := Clone(s, true)
	t.Log(c[0].XV(), c[1].XV())
}

func TestClone(t *testing.T) {
	n := 42
	pn := &n
	ppn := &pn
	var nilMap map[string]any
	src := &cloneStruct{
		S: "string",
		X: []int{1, 2, 3, 6, 8, 9},
		Y: map[string]any{
			"x": 1, "y": 2.2,
			"z": []int{1, 2, 3, 6, 8, 9},
		},
		Ptr:       pn,
		PtrPtr:    ppn,
		PtrPtrPtr: &ppn,
		A:         [5]uint64{1 << 2, 1 << 4, 1 << 6, 1 << 8, 1 << 10},
		SA:        []any{1, 2.2, "string", []int{1, 2, 3, 6, 8, 9}, nilMap},
		x:         n,

		C:  cloner{A: 420},
		C2: &cloner{A: 420},
		C3: cloner0{420, false},

		SS: simpleStruct{1, 2, "3", true},
	}

	dst := Clone(src, true)

	if dst == src {
		t.Fatal("cp == s")
	}

	if dst.Ptr == src.Ptr {
		t.Fatal("cp.Ptr == s.Ptr")
	}

	if dst.PtrPtr == src.PtrPtr {
		t.Fatal("cp.PtrPtr == s.PtrPtr")
	}

	if dst.PtrPtrPtr == src.PtrPtrPtr {
		t.Fatal("cp.PtrPtrPtr == s.PtrPtrPtr")
	}

	if src.x != dst.x {
		t.Fatal("src.x != dst.x", src.x, dst.x)
	}

	if !dst.C3.cloned {
		t.Fatal("!dst.C3.cloned")
	}

	dst.C3.cloned = false // so the next check passes

	if !reflect.DeepEqual(src, dst) {
		j1, _ := json.Marshal(src)
		j2, _ := json.Marshal(dst)
		t.Fatalf("!reflect.DeepEqual(src, dst)\nsrc: %s\n----\ndst: %s", j1, j2)
	}

	sj, _ := json.Marshal(src)
	dj, _ := json.Marshal(dst)
	if !bytes.Equal(sj, dj) {
		t.Fatalf("!bytes.Equal(src, dst):\nsrc: %s\ndst: %s", sj, dj)
	}

	dst = Clone(src, false)
	if dst.x == src.x {
		t.Fatal("src.x == dst.x", src.x, dst.x)
	}
	t.Logf("%s", sj)

	if dst.Y["z"].([]int)[0] = 42; src.Y["z"].([]int)[0] != 1 {
		t.Fatal("src.y == dst.y", src.Y, dst.Y)
	}
}

var cloneSink *cloneStruct

func BenchmarkClone(b *testing.B) {
	bp := &BrandProduct{
		ReviewLinks: &BrandProductReviewLink{
			AppleStore: "",
			GooglePlay: "",
			Instagram:  "",
			Leafly:     "",
			Weedmaps:   "",
		},
		Mappings: map[string][]string{
			"a": {"b", "c", "d"},
			"b": {"b", "c", "d"},
			"c": {"b", "c", "d"},
		},
		AltNames:        gsets.Of("a", "b", "c", "d"),
		Batches:         make([]*BrandProductBatch, 1024),
		RelatedProducts: make([]*BrandProductRelated, 1024),
	}
	for i := 0; i < 1024; i++ {
		bp.Batches[i] = &BrandProductBatch{
			Collectibles: make([]*Collectible, 1024),
		}
		bp.RelatedProducts[i] = &BrandProductRelated{}
	}
	_ = bp
	b.RunParallel(func(p *testing.PB) {
		var cloneSink BrandProduct
		for p.Next() {
			j, _ := json.Marshal(bp)
			json.Unmarshal(j, &cloneSink)
			// if Clone(bp, true) == nil {
			// 	b.Fatal("nil")
			// }
		}
	})
}

type BrandProduct struct {
	ReviewLinks *BrandProductReviewLink `json:"reviewLinks,omitempty"`
	Mappings    map[string][]string     `json:"sourceMapping,omitempty"`
	UserID      string                  `json:"userID,omitempty"`
	Name        string                  `json:"name,omitempty"`
	// AltNames are used to map other products to this one during brands-analysis, e.g., when different retailers have different names for this product
	AltNames          gsets.Strings          `json:"altNames,omitempty"`
	Brand             string                 `json:"brand,omitempty"`
	Category          string                 `json:"category,omitempty"`
	SubCategory       string                 `json:"subCategory,omitempty"`
	Sku               string                 `json:"sku,omitempty"`
	ProductShopURL    string                 `json:"productShopURL,omitempty"`
	ProductImageURL   string                 `json:"productImageURL,omitempty"`
	ID                string                 `json:"id,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Platform          string                 `json:"platform,omitempty"`
	Batches           []*BrandProductBatch   `json:"batches,omitempty"`
	RelatedProducts   []*BrandProductRelated `json:"relatedProducts,omitempty"`
	UpdatedAt         int64                  `json:"updated_at,omitempty"`
	AvgRetailPrice    float64                `json:"avgRetailPrice,omitempty"`
	AvgWholesalePrice float64                `json:"avgWholesalePrice,omitempty"`
	CreatedAt         int64                  `json:"created_at,omitempty"`
	IsReviewed        bool                   `json:"isReviewed,omitempty"`
	HideBatchInfo     bool                   `json:"hideBatchInfo,omitempty"`
	Archived          bool                   `json:"archived,omitempty"`
}

type BrandProductBatch struct {
	ID              string         `json:"id,omitempty"`
	ProdBatchNum    string         `json:"prodBatchNum,omitempty"`
	CannabinoidUnit string         `json:"cannabinoidUnit,omitempty"`
	RedirectURL     string         `json:"redirectURL,omitempty"`
	Collectibles    []*Collectible `json:"collectibles,omitempty"`
	Quantity        int            `json:"quantity,omitempty"`
	BatchDate       int64          `json:"batchDate,omitempty"`
	ThcPercent      float64        `json:"thcPercent,omitempty"`
	ThcaPercent     float64        `json:"thcaPercent,omitempty"`
	CbdaPercent     float64        `json:"cbdaPercent,omitempty"`
	CbcPercent      float64        `json:"cbcPercent,omitempty"`
	CbePercent      float64        `json:"cbePercent,omitempty"`
	CbgPercent      float64        `json:"cbgPercent,omitempty"`
	CbnPercent      float64        `json:"cbnPercent,omitempty"`
	Delta8Percent   float64        `json:"delta8Percent,omitempty"`
	TotalTHC        float64        `json:"totalTHC,omitempty"`
	TotalCanna      float64        `json:"totalCanna,omitempty"`
	CbdPercent      float64        `json:"cbdPercent,omitempty"`
	ShouldRedirect  bool           `json:"shouldRedirect,omitempty"`
}

type Collectible struct {
	ID         string  `json:"id"`
	SrcID      string  `json:"srcID,omitempty"`
	URL        string  `json:"url,omitempty"`
	ProductID  string  `json:"productID"`
	BatchID    string  `json:"batchID"`
	QR         []byte  `json:"qr,omitempty"`
	Rating     float64 `json:"rating,omitempty"`
	RedeemedAt int64   `json:"redeemedAt,omitempty"`
	Redeemed   bool    `json:"redeemed,omitempty"`
}

type BrandTemplate struct {
	UserID              string                 `json:"userID"`
	CollectibleRedirect string                 `json:"collectibleRedirect"`
	RelatedProducts     []*BrandProductRelated `json:"relatedProducts"`
	Accrual             float64                `json:"accrual,omitempty"`
	Created             int64                  `json:"created,omitempty" ts:"date,null"`
	Updated             int64                  `json:"updated,omitempty" ts:"date,null"`
}

type BrandProductRelated struct {
	Name     string `json:"name,omitempty"`
	ImageURL string `json:"imageURL,omitempty"`
	Link     string `json:"link,omitempty"`
}

type BrandProductReviewLink struct {
	AppleStore string `json:"appleStore,omitempty"`
	GooglePlay string `json:"googlePlay,omitempty"`
	Instagram  string `json:"instagram,omitempty"`
	Leafly     string `json:"leafly,omitempty"`
	Weedmaps   string `json:"weedmaps,omitempty"`
}
