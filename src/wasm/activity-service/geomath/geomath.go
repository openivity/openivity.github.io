// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package geomath

import (
	"math"
)

// HaversineDistance returns distance meters between two coordinates calculated using Haversine formula.
// The Haversine formula can result in an error of up to 0.5% since the Earth is not even a sphere — it’s an oblate ellipsoid.
//
// ref: https://en.wikipedia.org/wiki/Haversine_formula
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const radius = 6371.0

	lat1, lon1 = degreesToRadians(lat1), degreesToRadians(lon1)
	lat2, lon2 = degreesToRadians(lat2), degreesToRadians(lon2)

	// Haversine formula
	a := math.Pow(math.Sin((lat2-lat1)/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin((lon2-lon1)/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the distance (in km)
	distance := radius * c

	return distance * 1000 // in meters
}

// VincentyDistance returns distance meters between two coordinates calculated using Vincenty formula.
// The Vincenty formula is accurate up to 0.5mm. However, Vincenty's formula is generally more complex
// and computationally intensive compared to the Haversine formula.
//
// ref: https://en.wikipedia.org/wiki/Vincenty%27s_formulae
func VincentyDistance(lat1, lon1, lat2, lon2 float64) float64 {
	a := 6378137.0         // WGS84 semi-major axis
	f := 1 / 298.257223563 // WGS84 flattening
	b := (1 - f) * a
	tolerance := 1e-12

	lat1, lon1 = degreesToRadians(lat1), degreesToRadians(lon1)
	lat2, lon2 = degreesToRadians(lat2), degreesToRadians(lon2)

	L := lon2 - lon1
	U1 := math.Atan((1 - f) * math.Tan(lat1))
	U2 := math.Atan((1 - f) * math.Tan(lat2))
	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)

	lambda := L
	var lambdaP float64
	var sinSigma, cosSigma, sigma, sinAlpha, cosSqAlpha, cos2SigmaM, C float64

	for {
		sinLambda := math.Sin(lambda)
		cosLambda := math.Cos(lambda)
		sinSigma = math.Sqrt((cosU2*sinLambda)*(cosU2*sinLambda) + (cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda))
		if sinSigma == 0 {
			return 0.0 // Co-incident points
		}
		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha = 1 - sinAlpha*sinAlpha
		cos2SigmaM = cosSigma - 2*sinU1*sinU2/cosSqAlpha
		C = f / 16 * cosSqAlpha * (4 + f*(4-3*cosSqAlpha))
		lambdaP = lambda
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))
		if math.Abs(lambda-lambdaP) <= tolerance {
			break
		}
	}

	uSq := cosSqAlpha * (a*a - b*b) / (b * b)
	A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
	B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))
	deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))

	distance := b * A * (sigma - deltaSigma)

	return distance
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
