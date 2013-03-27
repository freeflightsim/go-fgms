

package fgms

import (
	"math"
)

// from .hxx
const SG_180 = 180.0
const SG_PI = 3.1415926535
const SG_RADIANS_TO_DEGREES = (SG_180/SG_PI)
const SG_DEGREES_TO_RADIANS = (SG_PI/SG_180)
const SG_FEET_TO_METER    = 0.3048

//from c.xx
/**
 *  High-precision versions of the above produced with an arbitrary
 * precision calculator (the compiler might lose a few bits in the FPU
 * operations).  These are specified to 81 bits of mantissa, which is
 * higher than any FPU known to me:
 */
const SQUASH  = 0.9966471893352525192801545;
const STRETCH = 1.0033640898209764189003079;
const POLRAD  = 6356752.3142451794975639668;

// Radians To Nautical Miles 
const SG_RAD_TO_NM  = 3437.7467707849392526

// Nautical Miles in a Meter 
const SG_NM_TO_METER  = 1852.0000

// Meters to Feet 
const SG_METER_TO_FEET  = 3.28083989501312335958

// PI2 
const SGD_PI_2    = 1.57079632679489661923


const ( X = 0 
		Y 
		Z 
)
const ( Lat = 0 
		Lon 
		Alt
)


type Point3D struct {
	x float64 
	y float64 
	z float64 
}
func (me *Point3D) Set(x, y, z float64){
	me.x = x
	me.y = y
	me.z = z
}
func (me *Point3D) Clear(){
	me.x = 0
	me.y = 0
	me.z = 0
}

func (me *Point3D) Length () float64 {
	//return (sqrt ((m_X * m_X) + (m_Y * m_Y) + (m_Z * m_Z)));
	return math.Sqrt( (me.x * me.x) + (me.y * me.y) + (me.z * me.z) )
}




func Point3DSubract(p1, p2 Point3D) Point3D{
	
	return Point3D{x: p1.x - p2.x, y: p1.y - p2.y, z: p1.z - p2.z}  
	
}

//////////////////////////////////////////////////////////////////////
/**
 * @brief Calculate distance of clients
 */
func Distance ( P1, P2 Point3D) float32 {
	
	//P = P1 - P2
	var P Point3D
	P = Point3DSubract( P1, P2)
	
	//return (float)(P.length() / SG_NM_TO_METER);
	return float32(P.Length() / SG_NM_TO_METER)
} // Distance ( const Point3D & P1, const Point3D & P2 )



//-------------------------------------------------------------------

//#define _EQURAD     (6378137.0)
const _EQURAD = 6378137.0

//#define E2 fabs(1 - SQUASH*SQUASH)
var e2 float64  = math.Abs( 1 - SQUASH * SQUASH )

//static double ra2 = 1/(_EQURAD*_EQURAD);
var ra2 float64 =  1 / (_EQURAD *_EQURAD)

//static double e2 = E2;
//static double e4 = E2*E2;
var e4 float64 = e2 * e2




/* 
 Convert a cartexian XYZ coordinate to a geodetic lat/lon/alt.
   This function is a copy of what's in SimGear,
  simgear/math/SGGeodesy.cxx and fgms http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_geometry.cxx#line407
*/

func SG_CartToGeod ( CartPoint Point3D ) Point3D {

	// according to
	// H. Vermeille,
	// Direct transformation from geocentric to geodetic cordinates,
	// Journal of Geodesy (2002) 76:451-454
	//double x = CartPoint[X];
	//double y = CartPoint[Y];
	//double z = CartPoint[Z];
	x := CartPoint.x
	y := CartPoint.y
	z := CartPoint.z
	
	//double XXpYY = x*x+y*y;
	var XXpYY float64 = (x * x) + (y * y)
	
	//double sqrtXXpYY = sqrt(XXpYY);
	var  sqrtXXpYY float64 = math.Sqrt(XXpYY)
	
	//double p = XXpYY*ra2;
	var p float64 = XXpYY * ra2
	
	//double q = z*z*(1-e2)*ra2;
	var q float64 = z*z*(1-e2) * ra2
	
	//double r = 1/6.0*(p+q-e4);
	var r float64 = 1 / 6.0 * (p + q - e4)
	
	//double s = e4*p*q/(4*r*r*r);
	var s float64 = e4 * p * q / (4 * r * r * r)
	
	//double t = pow(1+s+sqrt(s*(2+s)), 1/3.0);
	var t float64 = math.Pow(1 + s + math.Sqrt(s * (2 + s)), 1 / 3.0)
	
	//double u = r*(1+t+1/t);
	var u float64 = r * (1 + t + 1 / t)
	
	//double v = sqrt(u*u+e4*q);
	var v float64 = math.Sqrt(u * u + e4 * q)
	
	//double w = e2*(u+v-q)/(2*v);
	var w float64 = e2 * (u + v - q) / (2 * v)
	
	//double k = sqrt(u+v+w*w)-w;
	var k float64 = math.Sqrt(u + v + w * w) -w
	
	//double D = k*sqrtXXpYY/(k+e2);
	var D float64 = k * sqrtXXpYY / (k + e2)
	
	var GeodPoint Point3D
	//GeodPoint[Lon] = (2*atan2(y, x+sqrtXXpYY)) * SG_RADIANS_TO_DEGREES;
	GeodPoint.y = (2 * math.Atan2(y, x + sqrtXXpYY)) * SG_RADIANS_TO_DEGREES
	
	//double sqrtDDpZZ = sqrt(D*D+z*z);
	var sqrtDDpZZ float64 = math.Sqrt( D * D + z * z)
	//GeodPoint[Lat] = (2*atan2(z, D+sqrtDDpZZ)) * SG_RADIANS_TO_DEGREES;
	GeodPoint.x = (2* math.Atan2(z, D + sqrtDDpZZ)) * SG_RADIANS_TO_DEGREES
	
	//GeodPoint[Alt] = ((k+e2-1)*sqrtDDpZZ/k) * SG_METER_TO_FEET;
	GeodPoint.z = ((k + e2 - 1) * sqrtDDpZZ / k) * SG_METER_TO_FEET
	return GeodPoint
} // sgCartToGeod()

