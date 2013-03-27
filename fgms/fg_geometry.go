

package fgms


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

func (me *Point3D) Length () int {
	return 0 //(sqrt ((m_X * m_X) + (m_Y * m_Y) + (m_Z * m_Z)));
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


