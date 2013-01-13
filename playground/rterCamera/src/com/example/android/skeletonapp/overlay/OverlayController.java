package com.example.android.skeletonapp.overlay;

import java.util.Arrays;

import android.annotation.TargetApi;
import android.content.Context;
import android.hardware.Sensor;
import android.hardware.SensorEvent;
import android.hardware.SensorEventListener;
import android.hardware.SensorManager;
import android.os.Build;
import android.util.Log;

@TargetApi(Build.VERSION_CODES.GINGERBREAD)

/**
 * NORTH: 0 deg
 * EAST: +90 deg
 * WEST: -90 deg
 * SOUTH: +/- 180 deg
 * 
 * @author stepan
 *
 */
public class OverlayController implements SensorEventListener {
	protected float desiredOrientation;
	protected float currentOrientation;

	protected CameraGLSurfaceView mGLView;
	protected Context context;

	private static final String TAG = "OpenGL Overlay Controller";
	
	float[] aValues = new float[3];
	float[] mValues = new float[3];

	public OverlayController(Context context) {
		this.context = context;
		this.mGLView = new CameraGLSurfaceView(context);
	}

	/**
	 * @return the camera GLView
	 */
	public CameraGLSurfaceView getGLView() {
		return this.mGLView;
	}

	/**
	 * Set the desired absolute bearing
	 * 
	 * @param bearing
	 */
	public void setDesiredHeading(float heading) {
		desiredOrientation = heading;
	}

	/**
	 * Set the desired offset from the current bearing
	 * 
	 * @param offset
	 */
	public void setHeadingOffset(float offset) {
		desiredOrientation = currentOrientation + offset;
	}

	@Override
	public void onAccuracyChanged(Sensor sensor, int accuracy) {
		// TODO Auto-generated method stub

	}

	@Override
	public void onSensorChanged(SensorEvent event) {
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
		orientationValues[0] = (float) Math.toDegrees(orientationValues[0]);
		orientationValues[1] = (float) Math.toDegrees(orientationValues[1]);
		
		// this angle tells us the device orientation
		// between 90 and -90 is right side up (landscape); otherwise upside down
		orientationValues[2] = (float) Math.toDegrees(orientationValues[2]);
		
		Log.e(TAG, "x,y,z: " + Arrays.toString(orientationValues));
	}

}
