package am_test

import (
	"fmt"
	"github.com/kairos10/weqalign/api/am"
)

func ExampleGetPlateSolver() {
	am.AmShellCommands["solver"] = "./testres/solve.sh"

	solver := am.GetPlateSolver()
	
	img := am.ImageResource{FilePath: "testres/img1.jpg"}
	am.LoggerFunc = func(s string) { fmt.Println(s) }
	solver(&img)
	am.LoggerFunc = nil
	solver(&am.ImageResource{FilePath: "testres/img2.jpg"})
	// Output:
	// SOLVED: testres/img1
}

func ExamplePlateGetWcsInfo() {
	solver := am.GetPlateSolver()
	img := am.ImageResource{FilePath: "testres/img1.jpg"}
	solver(&img)
	am.PlateGetWcsInfo(&img)
	fmt.Println(img.NegParity)
	// Output: true
}

func ExamplePlateRD2XY() {
	solver := am.GetPlateSolver()
	img := am.ImageResource{FilePath: "testres/img1.jpg"}
	solver(&img)
	x, y := am.PlateRD2XY(&img, 0, 90) // get position for NCP[ra=0; dec=90°]
	fmt.Printf("x[%.1f] y[%.1f]", x, y)
	// Output: x[141.8] y[191.9]
}

func ExamplePlateGetStars() {
	solver := am.GetPlateSolver()
	img := am.ImageResource{FilePath: "testres/img1.jpg"}
	solver(&img)
	stars := am.PlateGetStars(&img)
	fmt.Println(len(stars), "stars")
	if len(stars) > 0 {
		fmt.Printf("1 @ x[%v]/y[%v]\n", stars[0].X, stars[0].Y)
	}
	// Output:
	// 15 stars
	// 1 @ x[147]/y[124]
}
