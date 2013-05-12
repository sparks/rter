package ca.nehil.rter.streamingapp2.overlay;

import java.util.Arrays;

import ca.nehil.rter.streamingapp2.overlay.CameraGLRenderer.Indicate;
import ca.nehil.rter.streamingapp2.util.*;

import android.content.Context;
import android.hardware.GeomagneticField;
import android.hardware.Sensor;
import android.hardware.SensorEvent;
import android.hardware.SensorEventListener;
import android.hardware.SensorManager;
import android.location.Location;
import android.location.LocationListener;
import android.os.Bundle;
import android.util.Log;

/**
 * NORTH: 0 deg
 * EAST: +90 deg
 * WEST: -90 deg
 * SOUTH: +/- 180 deg
 * 
 * @author stepan
 *
 */
public class OverlayController implements SensorEventListener, LocationListener {
	protected float desiredOrientation;
	protected float currentOrientation;
	protected float deviceOrientation;
	protected boolean rightSideUp = true;
	protected boolean freeRoam = true;

	protected CameraGLSurfaceView mGLView;
	protected CameraGLRenderer mGLRenderer;
	protected Context context;

	private static final String TAG = "OpenGL Overlay Controller";

	float[] aValues = new float[3];
	float[] mValues = new float[3];
	float declination = 0;	//geo magnetic declinationf from true North

	// max orientation tolerance in degrees
	public float orientationTolerance = 10.0f;
	
	private MovingAverageCompass orientationFilter;

	public OverlayController(Context context) {
		this.context = context;
		this.mGLView = new CameraGLSurfaceView(context);
		this.mGLRenderer = this.mGLView.getGLRenderer();
		
		orientationFilter = new MovingAverageCompass(30);
	}

	/**
	 * @return the camera GLView
	 */
	public CameraGLSurfaceView getGLView() {
		return this.mGLView;
	}

	/**
	 * when set to 'true', no indicator arrows will be given, and frame will be
	 * blue
	 * 
	 * @param freeRoam
	 */
	public void letFreeRoam(boolean freeRoam) {
		this.freeRoam = freeRoam;
	}
	
	public float getCurrentOrientation() {
		return this.currentOrientation;
	}
	/**
	 * Set the desired absolute bearing Should be between +180 and -180, but
	 * will work otherwise
	 * 
	 * @param orientation
	 */
	public void setDesiredOrientation(float orientation) {
		orientation = fixAngle(orientation);
		desiredOrientation = orientation;
	}

	/**
	 * makes sure angle is between -180 and 180
	 * 
	 * @param angle
	 * @return fixed angle
	 */
	protected float fixAngle(float angle) {
		if (angle > 180.0f) {
			angle = -180.0f + angle % 180;
		} else if (angle < -180.0f) {
			angle = 180.0f - Math.abs(angle) % 180;
		}

		return angle;
	}

	/**
	 * Set the desired offset from the current bearing should be between +180
	 * and -180, but will work otherwise
	 * 
	 * @param offset
	 */
	public void setOrientationOffset(float offset) {
		this.setDesiredOrientation(currentOrientation + offset);
	}

	@Override
	public void onAccuracyChanged(Sensor sensor, int accuracy) {
		// TODO Auto-generated method stub

	}

	@Override
	public void onSensorChanged(SensorEvent event) {
		/**
		 * code adapted from here:
		 * http://stackoverflow.com/questions/8989103/sensormanager
		 * -getorientation-gives-very-unstable-results
		 */
		switch (event.sensor.getType()) {
		case Sensor.TYPE_ACCELEROMETER:
			System.arraycopy(event.values, 0, aValues, 0, 3);
			break;
		case Sensor.TYPE_MAGNETIC_FIELD:
			System.arraycopy(event.values, 0, mValues, 0, 3);
			break;
		}
		float[] R = new float[16];
		float[] orientationValues = new float[3];

		if (aValues == null || mValues == null)
			return;

		if (!SensorManager.getRotationMatrix(R, null, aValues, mValues))
			return;

		float[] outR = new float[16];
		SensorManager.remapCoordinateSystem(R, SensorManager.AXIS_Z,
				SensorManager.AXIS_MINUS_X, outR);

		SensorManager.getOrientation(outR, orientationValues);

		// this angle tells us the orientation
		this.orientationFilter.pushValue((float) Math.toDegrees(orientationValues[0]));
		this.currentOrientation = this.orientationFilter.getValue() +this.declination;

		// this is not used currently, 90 when phone facing the sky, -90 when
		// facing the ground
		// orientationValues[1] = (float) Math.toDegrees(orientationValues[1]);

		// this angle tells us the device orientation
		// between 90 and -90 is right side up (landscape); otherwise upside
		// down
		this.deviceOrientation = (float) Math.toDegrees(orientationValues[2]);

//		Log.d("orientation", "x: " + String.format("%.1f", Math.toDegrees(orientationValues[0]))
//				+ ", y: " + String.format("%.1f", Math.toDegrees(orientationValues[1]))
//				+ ", z: " + String.format("%.1f", Math.toDegrees(orientationValues[2])));
//		Log.d("reald orientation", "x: " + String.format("%.1f", this.currentOrientation));

		if (this.freeRoam) {
			this.mGLRenderer.indicateTurn(Indicate.FREE, 0.0f);

			return;
		}

		// check orientation of device
		if (deviceOrientation <= 90.0f && deviceOrientation >= -90.0f) {
			this.rightSideUp = true;
		} else
			this.rightSideUp = false;

		// graphics logic
		boolean rightArrow = true;
		float difference = fixAngle(desiredOrientation - currentOrientation);
		if (Math.abs(difference) > orientationTolerance) {
			
			if (difference > 0) {
				// turn right
				rightArrow = true;
			} else {
				// turn left
				rightArrow = false;
			}

			// flip arrow incase device is flipped
			if (!this.rightSideUp) {
				rightArrow = !rightArrow;
			}

			if (rightArrow) {
				this.mGLRenderer.indicateTurn(Indicate.RIGHT,
						Math.abs(difference) / 180.0f);
			} else {
				this.mGLRenderer.indicateTurn(Indicate.LEFT,
						Math.abs(difference) / 180.0f);
			}

		} else {
			this.mGLRenderer.indicateTurn(Indicate.NONE, 0.0f);
		}

	}

	@Override
	public void onLocationChanged(Location loc) {
		//calculate and store declination for compass offsetting to true north
        GeomagneticField gmf = new GeomagneticField(
            (float)loc.getLatitude(), (float)loc.getLongitude(), (float)loc.getAltitude(), System.currentTimeMillis());
        this.declination = gmf.getDeclination();
	}

	@Override
	public void onProviderDisabled(String provider) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void onProviderEnabled(String provider) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void onStatusChanged(String provider, int status, Bundle extras) {
		// TODO Auto-generated method stub
		
	}

}
