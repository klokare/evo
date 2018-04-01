package boxes

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"

	"github.com/klokare/evo"
)

// Maximum distance possible on the board. Board width is always 2.0 independent of resolution
const maxDistance = 2.82842712474619  // greatest delta in 1 dimension is 2 (-1 to 1) so sqrt(2^2 + 2^2) = sqrt(8) = 2.82842712474619
const maxFitness = maxDistance * 75.0 // max distance times number of cases
const npr = 9                         // number of boxes per row in output image

// Evaluator evaluates the phenome using the visual discrimination experiment described in http://eplex.cs.ucf.edu/papers/stanley_alife09.pdf
type Evaluator struct {
	data       *mat.Dense
	small      []int // location of small's top lef
	big        []int // location of big's top left
	targets    []int // centre of big
	resolution int
	empty      *image.NRGBA
	filename   string
}

// NewEvaluator creates a new boxes evaluator with a resolution of res.
func NewEvaluator(res int) (e *Evaluator) {

	// Create the new evaluator
	e = &Evaluator{
		resolution: res,
	}

	// Create the cases
	e.makeTrials()

	// Create the empty image
	e.makeImage()
	return
}

// SetOutput tells the evaluator to record the decisions to an image. This is not concurrent safe and
// should only be used when processing a single phenome. See the callback in the boxes example.
func (e *Evaluator) SetOutput(filename string) { e.filename = filename }

// Evaluate the phenome against 75 trials
func (e Evaluator) Evaluate(p evo.Phenome) (r evo.Result, err error) {

	// Recording output
	var img *image.NRGBA
	var recording bool
	if e.filename != "" {
		recording = true
		img = copyImage(e.empty)
	}

	// Evaluate the trials
	var outputs evo.Matrix
	if outputs, err = p.Activate(e.data); err != nil {
		return
	}

	// Calculate the errors
	c := e.resolution * e.resolution
	sum := 0.0
	for i := 0; i < 75; i++ {

		// Note the value at the centre of the large box.
		t := e.targets[i]
		v := outputs.At(i, t)

		// Determine the maximum value
		mv := 0.0
		tmp := 0.0
		for j := 0; j < c; j++ {
			tmp = outputs.At(i, j)
			if mv < tmp {
				mv = tmp
			}
		}

		// Find the greatest distance to do a square with the maximum value. This penalises the
		// phenome for assigning maximum value to more than one location.
		d := 0.0
		for j := 0; j < c; j++ {
			if outputs.At(i, j) >= v {
				tmp = distance(e.resolution, t, j)
				if d < tmp {
					d = tmp
				}
				if recording {
					x, y := unindex(e.resolution, j)
					bx := i % npr
					by := (i - bx) / npr
					left := bx * (e.resolution + 1)
					top := by * (e.resolution + 1)
					img.SetNRGBA(left+x+1, top+y+1, color.NRGBA{R: 255, G: 0, B: 0, A: 128})
				}
			}
		}
		sum += d * d // distance squared
	}

	// Compute the root mean squared distance and then the fitness
	rmsd := math.Sqrt(sum / 75.0)
	r = evo.Result{
		ID:      p.ID,
		Fitness: (maxDistance - rmsd) / maxDistance * 100.0, // rescale so that perfect sore is 100.0 and worst score is 0.0
		Solved:  sum == 0.0,
	}

	// Complete recording
	if recording {
		var f *os.File
		if f, err = os.Create(e.filename); err != nil {
			return
		}
		defer f.Close()

		if err = png.Encode(f, img); err != nil {
			return
		}
	}
	return
}

// makeTrials creates the data for evaluation, returning the matrix an a slice with the centre of
// the larger box's index in each box's row in the matrix
func (e *Evaluator) makeTrials() {

	// Determine the unit size and value
	res := e.resolution
	u := res / 11            // unit size based on a smallest grid of 11
	w := 3*u - 1             // additional width/height of the larger box
	v := 11.0 / float64(res) // value of pixel scaled for resolution where res of 11 = 1.0

	// Initialise the matrix and targets for 75 trials
	e.data = mat.NewDense(75, res*res, nil)
	e.targets = make([]int, 75)
	e.small = make([]int, 75)
	e.big = make([]int, 75)

	// Iterate the cases
	rng := evo.NewRandom()
	pc := rng.Perm(res * res)
	for s := 0; s < 25; s++ {

		// Identify the centre of the small box for this case
		xs, ys := unindex(res, pc[s])

		// Iterate the positioning of the larger box
		for l := 0; l < 3; l++ {

			// Note the top left position via offset from small box
			xl, yl := xs, ys
			switch l {
			case 0:
				xl = (xl + 5) % res // 5 to the right
			case 1:
				yl = (yl + 5) % res // 5 down
			case 2:
				xl = (xl + 5) % res // 5 to the right
				yl = (yl + 5) % res // 5 down
			}

			// Wrap large box around if it goes off the "board"
			if (xl + w) > (res - 1) {
				xl = (xl + w) - (res - 1) - 1
			}
			if (yl + w) > (res - 1) {
				yl = (yl + w) - (res - 1) - 1
			}

			// Add the small box to the data
			row := make([]float64, res*res)
			row[index(res, xs, ys)] = v

			// Add the large box to the data
			for x := xl; x < xl+w; x++ {
				for y := yl; y < yl+w; y++ {
					row[index(res, x, y)] = v
				}
			}

			// Update the evaluator
			e.targets[s*3+l] = index(res, xl+w/2, yl+w/2)
			e.small[s*3+l] = index(res, xs, ys)
			e.big[s*3+l] = index(res, xl, yl)
			e.data.SetRow(s*3+l, row)
		}
	}
	return
}

// return the index of the point in the 1-dimensional array which is row ordered
func index(res, x, y int) int {
	return x*res + y
}

func unindex(res, i int) (x, y int) {
	y = i % res
	x = (i - y) / res
	return
}

func coordinates(res, x, y int) (float64, float64) {
	n := float64(res - 1)
	return float64(x)/n*2.0 - 1.0, float64(y)/n*2.0 - 1.0
}

func distance(res, a, b int) float64 {
	x, y := unindex(res, a)
	ax, ay := coordinates(res, x, y)
	x, y = unindex(res, b)
	bx, by := coordinates(res, x, y)
	return math.Sqrt((ax-bx)*(ax-bx) + (ay-by)*(ay-by))
}

func (e *Evaluator) makeImage() {

	// Determine the unit size and value
	res := e.resolution
	u := res / 11 // unit size based on a smallest grid of 11

	// Create the image
	w := e.resolution + 1 // Width of image is box+1 so we can have a border. left and top will not have border line
	e.empty = image.NewNRGBA(image.Rect(0, 0, npr*w+1, npr*w+1))
	for i := 0; i < 75; i++ {

		// Determine box left and top
		bx := i % npr
		by := (i - bx) / npr
		left := bx * w
		top := by * w

		// Draw borders
		for xx := left; xx <= left+w; xx++ {
			e.empty.SetNRGBA(xx, top, color.NRGBA{R: 0, G: 0, B: 0, A: 128})
			e.empty.SetNRGBA(xx, top+w, color.NRGBA{R: 0, G: 0, B: 0, A: 128})
		}
		for yy := top; yy <= top+w; yy++ {
			e.empty.SetNRGBA(left, yy, color.NRGBA{R: 0, G: 0, B: 0, A: 128})
			e.empty.SetNRGBA(left+w, yy, color.NRGBA{R: 0, G: 0, B: 0, A: 128})
		}

		// Draw the small square
		sx, sy := unindex(e.resolution, e.small[i])
		for xx := left + sx + 1; xx < left+sx+u+1; xx++ {
			for yy := top + sy + 1; yy < top+sy+u+1; yy++ {
				e.empty.Set(xx, yy, color.NRGBA{R: 0, G: 0, B: 255, A: 128})
			}
		}

		// Draw the big square
		lx, ly := unindex(e.resolution, e.big[i])
		for xx := left + lx + 1; xx < left+lx+u*3+1; xx++ {
			for yy := top + ly + 1; yy < top+ly+u*3+1; yy++ {
				e.empty.Set(xx, yy, color.NRGBA{R: 0, G: 0, B: 255, A: 128})
			}
		}
	}

}

func copyImage(a *image.NRGBA) (b *image.NRGBA) {
	r := a.Bounds()
	b = image.NewNRGBA(r)
	for x := 0; x < r.Max.X; x++ {
		for y := 0; y < r.Max.Y; y++ {
			b.Set(x, y, a.At(x, y))
		}
	}
	return
}
