package extract

// this is basically a clone of go-colorweave

import (
	"github.com/aaronland/go-colours/closest"
	"github.com/nfnt/resize"
	"image"
	"log"
	"sort"
	"sync"
	"time"
)

func Extract(im image.Image) ([]string, error) {

	im = resize.Resize(100, 0, im, resize.Bilinear)

	bounds := im.Bounds()

	// wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)

	lookup := make(map[string]int)

	for i := 0; i <= bounds.Max.X; i++ {

		for j := 0; j <= bounds.Max.Y; j++ {

			pixel := im.At(i, j)
			red, green, blue, _ := pixel.RGBA()
			rgb := []int{int(red / 255), int(green / 255), int(blue / 255)}

			c := closest.Closest(rgb, "foo")

			log.Println(i, j, c)

			t1 := time.Now()

			mu.Lock()

			count, ok := lookup[c]

			if ok {
				count += 1
			} else {
				count = 1
			}

			lookup[c] = count
			mu.Unlock()

			t2 := time.Since(t1)

			log.Println(i, j, t2)

		}
	}

	keys := make([]int, 0, len(lookup))

	for _, val := range lookup {
		keys = append(keys, val)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	colours := make([]string, 0)

	for _, c := range reverse_map(lookup) {
		colours = append(colours, c)
	}

	return colours, nil
}

func reverse_map(m map[string]int) map[int]string {
	n := make(map[int]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}
