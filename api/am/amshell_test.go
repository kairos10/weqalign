package am_test

import (
	"fmt"
	"github.com/kairos10/weqalign/api/am"
)

func ExampleGetPlateSolver() {
	solver := am.GetPlateSolver()

	ok := solver(&am.ImagePlate{FilePath: "testres/img2.jpg"})
	fmt.Println(ok)

	am.AmShellCommands["solver"] = "./testres/solve.sh"
	img := am.ImagePlate{FilePath: "testres/img1.jpg"}
	am.LoggerFunc = func(s string) { fmt.Println(s) }
	solver(&img)
	am.LoggerFunc = nil
	// Output:
	// true
	// SOLVED: testres/img1
}

func ExamplePlateGetWcsInfo() {
	solver := am.GetPlateSolver()
	img := am.ImagePlate{FilePath: "testres/img1.jpg"}
	solver(&img)
	am.PlateGetWcsInfo(&img)
	fmt.Printf("parity[%v] RA[%.8f] DEC[%.8f]\n", img.NegParity, img.RACenter, img.DECCenter)
	// Output: 
	// parity[true] RA[153.09812355] DEC[83.58660028]
}

func ExamplePlateRD2XY() {
	solver := am.GetPlateSolver()
	img := am.ImagePlate{FilePath: "testres/img1.jpg"}
	solver(&img)
	x, y := am.PlateRD2XY(&img, 0, 90) // get position for NCP[ra=0; dec=90°]
	fmt.Printf("x[%.1f] y[%.1f]", x, y)
	// Output: x[141.8] y[191.9]
}

func ExamplePlateGetStars() {
	solver := am.GetPlateSolver()
	img := am.ImagePlate{FilePath: "testres/img1.jpg"}
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
