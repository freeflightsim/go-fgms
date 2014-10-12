

package message


/* T_PositionMsg - Position Message

	Original Source: http://gitorious.org/fgms/fgms-0-x/blobs/master/src/flightgear/MultiPlayer/mpmessages.hxx#line78
	Note:
		all the important values are float32
		with the exception of position which is float64
		This caused a clash with Point3D which needs to be either 32 or 64
		- For now the 32's are converted to 64's
*/
type PositionMsg struct{

	/// Name of the aircraft model
	// - char Model[MAX_MODEL_NAME_LEN];
	ModelBytes [MAX_MODEL_NAME_LEN]byte

	// Time when this packet was generated
	// - xdr_data2_t time;
	Time uint64

	/// Time offset for network lag ?
	// - xdr_data2_t lag;
	Lag uint64

	// Position wrt the earth centered frame
	// - xdr_data2_t position[3];
	Position [3]float64


	// Orientation wrt the earth centered frame, stored in the angle axis
	// representation where the angle is coded into the axis length
	// - xdr_data_t orientation[3];
	Orientation [3]float32 //uint32

	// Linear velocity wrt the earth centered frame measured in the earth centered frame
	// - xdr_data_t linearVel[3];
	LinearVel [3]float32 //uint32

	// Angular velocity wrt the earth centered frame measured in the earth centered frame
	// - xdr_data_t angularVel[3];
	AngularVel [3]float32 // uint32

	// Linear acceleration wrt the earth centered frame measured in the earth centered frame
	// - xdr_data_t linearAccel[3];
	LinearAccel [3]float32 // uint32

	// Angular acceleration wrt the earth centered frame measured in the earth centered frame
	// - xdr_data_t angularAccel[3];
	AngularAccel [3]float32 //uint32
}

// Returns the Model as a string

func (me *PositionMsg) Model() string{
	return string(me.ModelBytes[:])

}
