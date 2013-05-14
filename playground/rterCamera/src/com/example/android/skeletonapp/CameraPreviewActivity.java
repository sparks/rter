/*
 * Copyright (C) 2007 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.example.android.skeletonapp;

import java.util.Random;



import com.example.android.skeletonapp.overlay.*;


import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.hardware.Camera;
import android.hardware.Camera.CameraInfo;
import android.hardware.Camera.PictureCallback;
import android.hardware.Sensor;
import android.hardware.SensorManager;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.os.PowerManager;

import android.util.Log;
import android.view.Gravity;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.Window;
import android.view.WindowManager;
import android.widget.FrameLayout;
import android.widget.Toast;
import android.location.Criteria;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.provider.Settings;

// ----------------------------------------------------------------------

public class CameraPreviewActivity extends Activity implements OnClickListener,
		LocationListener {
	private Preview mPreview;
	private FrameLayout mFrame; // need this to merge camera preview and openGL
								// view
	private CameraGLSurfaceView mGLView;
	private OverlayController overlay;
	SensorManager mSensorManager;
	Sensor mAcc, mMag;
	Camera mCamera;
	int numberOfCameras;
	int cameraCurrentlyLocked;
	PictureCallback mPicture = null;
	// The first rear facing camera
	int defaultCameraId;
	static boolean isFPS = false;

	static float justtesting;

	String selected_uid; // passed from other activity right now

	private LocationManager locationManager;
	private String provider;
	
	
	String internalLat=null;
	String internalLng=null;
	FrameInfo frameInfo;

	// to prevent sleeping
	PowerManager pm;
	PowerManager.WakeLock wl;

	private static final String TAG = "CameraPreview Activity";
	protected static final String MEDIA_TYPE_IMAGE = null;
	
	

	WifiManager myWifiManager;
	WifiInfo myWifiInfo; 
	BroadcastReceiver receiver;
	
	
		
	
	@SuppressLint("ParserError")
	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		Log.e(TAG, "onCreate");
		// Hide the window title.
		requestWindowFeature(Window.FEATURE_NO_TITLE);
		getWindow().addFlags(WindowManager.LayoutParams.FLAG_FULLSCREEN);
		String AndroidId = Settings.Secure.getString(getContentResolver(),
		         Settings.Secure.ANDROID_ID);
		frameInfo = new FrameInfo();

		// openGL overlay
		overlay = new OverlayController(this);

		// orientation
		mSensorManager = (SensorManager) getSystemService(SENSOR_SERVICE);
		mAcc = mSensorManager.getDefaultSensor(Sensor.TYPE_ACCELEROMETER);
		mMag = mSensorManager.getDefaultSensor(Sensor.TYPE_MAGNETIC_FIELD);

		// the frame layout will contain the camera preview and the gl view
		mFrame = new FrameLayout(this);

		// Create a RelativeLayout container that will hold a SurfaceView,
		// and set it as the content of our activity.
		mPreview = new Preview(this);

		// openGLview
		mGLView = overlay.getGLView();

		
		
		
		Log.e(TAG, "Fileoutput in phone id " + AndroidId);
		// uid = convertStringToByteArray(sUID);

		String androidIds[] = {
			"1e7f033bfc7b3625fa07c9a3b6b54d2c81eeff98",
			"fe7f033bfc7b3625fa06c9a3b6b54b2c81eeff98",
			"b6200c5cc15cfbddde2874c40952a7aa25a869dd",
			"852decd1fbc083cf6853e46feebb08622d653602",
			"e1830fcefc3f47647ffa08350348d7e34b142b0b",
			"48ad32292ff86b4148e0f754c2b9b55efad32d1e",
			"acb519f53a55d9dea06efbcc804eda79d305282e",
			"ze7f033bfc7b3625fa06c5a316b54b2c81eeff98",
			"t6200c5cc15cfbddde2875c41952a7aa25a869dd",
			"952decd1fbc083cf6853e56f1ebb08622d653602",
			"y1830fcefc3f47647ffa05351348d7e34b142b0b",
			"x8ad32292ff86b4148e0f55412b9b55efad32d1e",
			"qcb519f53a55d9dea06ef5cc104eda79d305282e"
		};
		
		int rnd = new Random().nextInt(androidIds.length);
		
	    
		
		selected_uid = AndroidId; //AndroidId;
		frameInfo.uid = selected_uid.getBytes();
		//Log.e(TAG, "selected_uid in phone id" + selected_uid);
		Log.e(TAG, "selected_uid in phone id " + new String(frameInfo.uid));
		// add the two views to the frame
		mFrame.addView(mPreview);
		mFrame.addView(mGLView);

		mPreview.setOnClickListener(this);

		setContentView(mFrame);

		// Find the total number of cameras available
		numberOfCameras = Camera.getNumberOfCameras();

		// Find the ID of the default camera
		CameraInfo cameraInfo = new CameraInfo();
		for (int i = 0; i < numberOfCameras; i++) {
			Camera.getCameraInfo(i, cameraInfo);
			if (cameraInfo.facing == CameraInfo.CAMERA_FACING_BACK) {
				defaultCameraId = i;
			}
		}
		
		Log.e("mac", "WIFI Requesting location");
		this.wifiLocalization();
		
		// Get the location manager
		locationManager = (LocationManager) getSystemService(Context.LOCATION_SERVICE);
		// Define the criteria how to select the location provider -> use
		// default
		
		if ( !locationManager.isProviderEnabled( LocationManager.GPS_PROVIDER ) ) {
			Log.e(TAG, "GPS not available");
	    }	
		Criteria criteria = new Criteria();
		provider = locationManager.getBestProvider(criteria, false);
		Log.e(TAG, "Requesting location");
		locationManager.requestLocationUpdates(provider, 0, 1, this);
		// register the overlay control for location updates as well, so we get the geomagnetic field
		locationManager.requestLocationUpdates(provider, 0, 1000, overlay);
		if (provider != null) {
			Location location = locationManager.getLastKnownLocation(provider);
			// Initialize the location fields
			if (location != null) {
				System.out.println("Provider " + provider
						+ " has been selected. and location "+ location);
				onLocationChanged(location);
			} else {
				Log.d(TAG, "Location not available");
			}
		}

		// power manager
		pm = (PowerManager) getSystemService(Context.POWER_SERVICE);
		wl = pm.newWakeLock(PowerManager.SCREEN_BRIGHT_WAKE_LOCK, TAG);

		// test, set desired orienation to north
		overlay.letFreeRoam(false);
		overlay.setDesiredOrientation(0.0f);
		CharSequence text = "Tap to start..";
		int duration = Toast.LENGTH_SHORT;

		Toast toast = Toast.makeText(this, text, duration);
		toast.setGravity(Gravity.TOP|Gravity.RIGHT, 0, 0);
		toast.show();
		
		
		
	}
	
	private void wifiLocalization() {
		// TODO Auto-generated method stub
	
		
		myWifiManager = (WifiManager)getSystemService(Context.WIFI_SERVICE);
		myWifiInfo = myWifiManager.getConnectionInfo();
		Log.d("mac address", "WIFI ="+myWifiInfo.getBSSID());
		Log.d("mac address", "WIFI ="+myWifiInfo.getSSID());
		Log.d("mac address", "WIFI ="+myWifiInfo.getMacAddress());
		// Register Broadcast Receiver
				if (receiver == null)
					receiver = new WiFiScanReceiver(this);

				registerReceiver(receiver, new IntentFilter(
						WifiManager.SCAN_RESULTS_AVAILABLE_ACTION));
		
	}

	@Override
	protected void onResume() {
		super.onResume();
		locationManager.requestLocationUpdates(provider, 0, 1, this);
		// register the overlay control for location updates as well, so we get the geomagnetic field
		locationManager.requestLocationUpdates(provider, 0, 1000, overlay);
		// Open the default i.e. the first rear facing camera.
		mCamera = Camera.open();
		cameraCurrentlyLocked = defaultCameraId;
		mPreview.setCamera(mCamera);
		
		// sensors
		mSensorManager.registerListener(overlay, mAcc,
				SensorManager.SENSOR_DELAY_NORMAL);
		mSensorManager.registerListener(overlay, mMag,
				SensorManager.SENSOR_DELAY_NORMAL);

		// acquire wake lock to make sure camera preview remains on and bright
		wl.acquire();
	}

	@Override
	protected void onPause() {
		super.onPause();
		Log.e(TAG, "onPause");
		locationManager.removeUpdates(this);
		locationManager.removeUpdates(overlay);
		
		// stop sensor updates
		mSensorManager.unregisterListener(overlay);
		
		// end photo thread
		isFPS = false;
//		if(photoThread.isAlive()) {
//			photoThread.stopPhotos();
//			try {
//				photoThread.join();
//			} catch (InterruptedException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//		}
		// Because the Camera object is a shared resource, it's very
		// important to release it when the activity is paused.
		if (mCamera != null) {
			mPreview.setCamera(null);
			mCamera.release();
			mCamera = null;
			mPreview.inPreview = false;
		}

		mSensorManager.unregisterListener(overlay);

		// release wake lock to allow phone to sleep
		wl.release();
	}

	public void onClick(View v) {
		// TODO Auto-generated method stub
		Log.e(TAG, "onClick");
		isFPS = !isFPS;
		Log.e(TAG, "onClick changes isFPS : " + isFPS);
		if (isFPS) {
			//photoThread.start();
			CharSequence text = "Starting Photo Stream ..";
			int duration = Toast.LENGTH_SHORT;

			Toast toast = Toast.makeText(this, text, duration);
			toast.setGravity(Gravity.TOP|Gravity.RIGHT, 0, 0);
			toast.show();
			mCamera.takePicture(null, null, photoCallback);
			Log.d(TAG, "starting picture thread");
			mPreview.inPreview = false;
		}

	}

	Camera.PictureCallback photoCallback = new Camera.PictureCallback() {
		public void onPictureTaken(byte[] data, Camera camera) {
			final Camera tmpCamera = camera;
			final byte[] tmpData = data;
			(new Thread(new Runnable() {
				public void run() {
					Log.e(TAG, "Inside Picture Callback");
					runOnUiThread(new Runnable() {
		                 public void run() {

		                     Toast toast = Toast.makeText(CameraPreviewActivity.this,"Streaming..",Toast.LENGTH_LONG);
		                     toast.setGravity(Gravity.TOP|Gravity.RIGHT, 0, 0);
		                     toast.show();
		                }
		            });
					Looper.prepare();
					
					
					// get orientation
					frameInfo.orientation = convertStringToByteArray(""
							+ overlay.getCurrentOrientation());
					new SavePhotoTask(overlay).execute(tmpData, frameInfo.uid, frameInfo.lat, frameInfo.lon,
							frameInfo.orientation);
					if (isFPS) {
						tmpCamera.startPreview();
						mPreview.inPreview = true;

					}

					long start_time = System.currentTimeMillis();
					while (System.currentTimeMillis() < start_time + 5000 && isFPS) {
						Thread.yield();
					}

					if (isFPS) {
						
						Log.d(TAG, "Picture taken");
						mCamera.takePicture(null, null, photoCallback);
					}
				}
			})).start();

		}

	};

	public static byte[] convertStringToByteArray(String s) {

		byte[] theByteArray = s.getBytes();

		Log.e(TAG, "length of byte array" + theByteArray.length);
		return theByteArray;

	}

	protected void onSaveInstanceState(Bundle outState) {
		super.onSaveInstanceState(outState);
	}

	@Override
	public void onLocationChanged(Location location) {
		// TODO Auto-generated method stub
		Log.d(TAG, "Location Changed");
//		String lati = "" + (location.getLatitude());
//		String longi = "" + (location.getLongitude());
		String lati = internalLat;
		String longi = internalLng;
		frameInfo.lat = convertStringToByteArray(lati);
		frameInfo.lon = convertStringToByteArray(longi);
		
		

	}

	@Override
	public void onProviderDisabled(String arg0) {
		// TODO Auto-generated method stub

	}

	@Override
	public void onProviderEnabled(String arg0) {
		// TODO Auto-generated method stub

	}

	@Override
	public void onStatusChanged(String arg0, int arg1, Bundle arg2) {
		// TODO Auto-generated method stub
	}
}

class FrameInfo {
	public byte[] uid;
	public byte[] lat;
	public byte[] lon;
	public byte[] orientation;
}

