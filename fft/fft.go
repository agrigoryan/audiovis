package fft

import (
	"errors"
	"math"
	"math/cmplx"
)

func Fft(input []float64, output []complex128) error {
	if len(input) > len(output) {
		return errors.New("input length is greater than the output length")
	}
	ditFft(input, output, len(input), 1)
	// dft(input, output)
	return nil
}

func ditFft(input []float64, output []complex128, n int, s int) {
	if n == 1 {
		output[0] = complex(input[0], 0)
		return
	}
	ditFft(input, output, n/2, 2*s)
	ditFft(input[s:], output[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * output[k+n/2]
		output[k], output[k+n/2] = output[k]+tf, output[k]-tf
	}
}

func dft(input []float64, output []complex128) {
	n := len(input)
	for f := 0; f < n; f++ {
		period := float64(f) / float64(n)
		output[f] = 0
		for i := 0; i < n; i++ {
			angle := 2.0 * math.Pi * period * float64(i)
			output[i] += complex(input[i]*math.Sin(angle), input[i]*math.Cos(angle))
			// euler's formula e^(i*x) = cos(x) + i*sin(x)
			// output[i] += cmplx.Exp(complex(0, 2.0*math.Pi*period*float64(i)))
		}
	}
}
