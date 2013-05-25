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

package ca.nehil.rter.streamingapp2;

import java.io.BufferedReader;
import java.io.File;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintStream;
import java.net.Socket;
import java.util.Date;
import java.util.Random;


import ca.nehil.rter.streamingapp2.overlay.*;


import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.SharedPreferences;
import android.hardware.Camera;
import android.hardware.Camera.CameraInfo;
import android.hardware.Camera.PictureCallback;
import android.hardware.Sensor;
import android.hardware.SensorManager;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
import android.os.Bundle;
import android.os.Environment;
import android.os.Handler;
import android.os.Looper;
import android.os.ParcelFileDescriptor;
import android.os.PowerManager;
import android.provider.Settings.Secure;

import android.util.Log;
import android.view.Gravity;
import android.view.Menu;
import android.view.MenuItem;
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
import android.media.MediaRecorder;
import android.provider.Settings;
import android.net.wifi.WifiConfiguration;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;

// ----------------------------------------------------------------------

public class StreamingActivity extends Activity implements 
		LocationListener {
	
	public static final int SERVERPORT = 1200;
	public static String SERVERIP="192.168.30.8";
	Socket clientSocket;
	private SharedPreferences cookies;
	private SharedPreferences.Editor prefEditor;
	private String setRterResource;
	private String setRterSignature;
	private Handler handler = new Handler();
	ParcelFileDescriptor pfd=null;
	
	private Preview mPreview;
	public MediaRecorder mrec = new MediaRecorder();
	private FrameLayout mFrame; // need this to merge camera preview and openGL
								// view
	private CameraGLSurfaceView mGLView;
	private OverlayController overlay;
	SensorManager mSensorManager;
	Sensor mAcc, mMag;
	Camera mCamera;
	int numberOfCameras;
	int cameraCurrentlyLocked;
	//PictureCallback mPicture = null;
	// The first rear facing camera
	int defaultCameraId;
	static boolean isFPS = false;

	static float justtesting;

	private String AndroidId;
	private String selected_uid; // passed from other activity right now

	
	WifiManager myWifiManager;
	WifiInfo myWifiInfo; 
	BroadcastReceiver receiver;
	private LocationManager locationManager;
	private String provider;
	
	FrameInfo frameInfo;

	// to prevent sleeping
	PowerManager pm;
	PowerManager.WakeLock wl;

	private static final String TAG = "Streaming Activity";
	//protected static final String MEDIA_TYPE_IMAGE = null;
	
	private String[][] wifiMap= {
			{"00:1f:45:f3:1e:11","lat","lng"}	
	};
	
	
	public class SendVideoThread implements Runnable{
	    public void run(){
	        // From Server.java
	        try {
	            if(SERVERIP!=null){
	                handler.post(new Runnable() {
	                    @Override
	                    public void run() {
	                        Log.d(TAG, "Listening on IP: " + SERVERIP);
	                    }
	                });
//	                String host="132.206.74.113";
//	                int port = 9999;
	                Log.d(TAG,"Client will attempt connecting to server at host=" + SERVERIP + " port=" + SERVERPORT + ".");
	                clientSocket = new Socket(SERVERIP,SERVERPORT);
	                pfd  = ParcelFileDescriptor.fromSocket(clientSocket);
	             // ok, got a connection.  Let's use java.io.* niceties to read and write from the connection.
        			
	                //to send the android id uncomment below 3 line
        			PrintStream myOutput = new PrintStream(clientSocket.getOutputStream());	
                
//        			myOutput.print(AndroidId+";");
        			
        			// see if the server writes something back.
        			
	                
	            }
	        } catch (Exception e){
	            
	            e.printStackTrace();
	            System.out.println("Whoops, something bad happened!  I'm outta here.");
	        }
	        // End from server.java
	    }
	}
	
	
	
	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		// Inflate the menu; this adds items to the action bar if it is present.
		menu.add(0, 0, 0, "Start");
		return true;
	}
	
	@Override
    public boolean onOptionsItemSelected(MenuItem item)
    {
        if(item.getTitle().equals("Start"))
        {
            try {
                long starttime =  System.currentTimeMillis();
            	startRecording();
                item.setTitle("Stop");
                
                
                

            } catch (Exception e) {

                String message = e.getMessage();
                Log.e(TAG, "Problem " + message);
                mrec.release();
            }

        }
        else if(item.getTitle().equals("Stop"))
        {
            Log.e(TAG,"Alert Stop mrec");
        	mrec.stop();
            mrec.release();
            mrec = null;
            try {
				clientSocket.close();
			} catch (IOException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
            item.setTitle("Start");
        }

        return super.onOptionsItemSelected(item);
    }
	
	@Override
    protected void onDestroy() {
		if(mrec!=null)
        {
            mrec.stop();
            mrec.release();
            mCamera.release();
            mCamera.lock();
            try {
				clientSocket.close();
			} catch (IOException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
        }
		stopRecording();
        super.onDestroy();
    }
	
	protected void startRecording() throws IOException
    {
        if(mCamera==null)
         mCamera = Camera.open();
        
        String filename;
        String root = (Environment.getExternalStorageDirectory()).toString();
 		File rootDir = new File(Environment.getExternalStorageDirectory()
 				+ File.separator + "rter" + File.separator);
 		rootDir.mkdirs();
        Date date=new Date();
        filename="/rec"+date.toString().replace(" ", "_").replace(":", "_")+".ts";
         
        //create empty file it must use
        File file=new File(rootDir,filename);
         
        mrec = new MediaRecorder(); 

        mCamera.lock();
        mCamera.unlock();

        
        // Please maintain sequence of following code. 

        // If you change sequence it will not work
        mrec.setCamera(mCamera);    
        mrec.setVideoSource(MediaRecorder.VideoSource.CAMERA);
        //mrec.setAudioSource(MediaRecorder.AudioSource.MIC);     
//        mrec.setOutputFormat(MediaRecorder.OutputFormat.MPEG_4);
        mrec.setOutputFormat(8);
        mrec.setVideoEncoder(MediaRecorder.VideoEncoder.H264);
        //mrec.setAudioEncoder(MediaRecorder.AudioEncoder.DEFAULT);
//        mrec.setOutputFile(rootDir+filename);
        mrec.setOutputFile(pfd.getFileDescriptor());
        mrec.setVideoEncodingBitRate(600000);
        //mrec.setAudioEncodingBitRate(44100);
        mrec.setVideoFrameRate(15);
        mrec.setMaxDuration(-1);
        mrec.setPreviewDisplay(mPreview.mHolder.getSurface());
        
        mrec.prepare();
        mrec.start();
    }

    protected void stopRecording() {

        if(mrec!=null)
        {
            mrec.stop();
            mrec.release();
            mCamera.release();
            mCamera.lock();
        }
    }

    private void releaseMediaRecorder() {

        if (mrec != null) {
            mrec.reset(); // clear recorder configuration
            mrec.release(); // release the recorder object
        }
    }

    private void releaseCamera() {
        if (mCamera != null) {
            mCamera.release(); // release the camera for other applications
            mCamera = null;
        }

    }
	
	@SuppressLint("ParserError")
	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		Log.e(TAG, "onCreate");
		// Hide the window title.
//		requestWindowFeature(Window.FEATURE_NO_TITLE);
//		getWindow().addFlags(WindowManager.LayoutParams.FLAG_FULLSCREEN);
		AndroidId = Settings.Secure.getString(getContentResolver(),
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

		// add the two views to the frame
		mFrame.addView(mPreview);
		mFrame.addView(mGLView);

//		mPreview.setOnClickListener(this);

		setContentView(mFrame);

		// Find the total number of cameras available
		numberOfCameras = Camera.getNumberOfCameras();

		// Find the ID of the default camera
		CameraInfo cameraInfo = new CameraInfo();
		for (int i = 0; i < numberOfCameras; i++) {
			Camera.getCameraInfo(i, cameraInfo);
			Log.d(TAG, "Camera Id is " + i);
			if (cameraInfo.facing == CameraInfo.CAMERA_FACING_BACK) {
				Log.d(TAG, "Back facing camera chosen");
				defaultCameraId = i;
				Log.d(TAG, "defaultcamera ID :"+defaultCameraId );
				break;
			}
//			else if (cameraInfo.facing == CameraInfo.CAMERA_FACING_FRONT) {
//				Log.d(TAG, "Front facing camera chosen");
//				defaultCameraId = i;
//				Log.d(TAG, "defaultcamera ID :"+defaultCameraId );
//			}
		}
		cookies = getSharedPreferences("RterUserCreds", MODE_PRIVATE);
		prefEditor = cookies.edit();
		
		setRterResource = cookies.getString("rter_resource", "not-set");
		setRterSignature = cookies.getString("rter_signature", "not-set");
		
		
		Log.d(TAG, "Prefs ==> rter_resource:"+setRterResource+" :: rter_signature:" + setRterSignature );
		
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
	
		
		// Run new thread to handle socket communications
	    Thread sendVideo = new Thread(new SendVideoThread());
	    sendVideo.start();
	}
	
	@Override
	public void onStop() {
		unregisterReceiver(receiver);
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
		
		
		
		
		for(int i = 0;i <= wifiMap.length-1; i++){
			if(wifiMap[i][0].matches(myWifiInfo.getBSSID()))
			{
				Log.d("WIFI", "WIFI: lat= "+wifiMap[i][1] +" and lng= "+wifiMap[i][2]);
			}
		}
		
		
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
		//isFPS = false;
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

//	public void onClick(View v) {
//		// TODO Auto-generated method stub
//		Log.e(TAG, "onClick");
//		isFPS = !isFPS;
//		Log.e(TAG, "onClick changes isFPS : " + isFPS);
//		if (isFPS) {
//			//photoThread.start();
//			CharSequence text = "Starting Photo Stream ..";
//			int duration = Toast.LENGTH_SHORT;
//
//			Toast toast = Toast.makeText(this, text, duration);
//			toast.setGravity(Gravity.TOP|Gravity.RIGHT, 0, 0);
//			toast.show();
////			mCamera.takePicture(null, null, photoCallback);
//			Log.d(TAG, "starting picture thread");
//			mPreview.inPreview = false;
//		}
//
//	}

//	Camera.PictureCallback photoCallback = new Camera.PictureCallback() {
//		public void onPictureTaken(byte[] data, Camera camera) {
//			final Camera tmpCamera = camera;
//			final byte[] tmpData = data;
//			(new Thread(new Runnable() {
//				public void run() {
//					Log.e(TAG, "Inside Picture Callback");
//					runOnUiThread(new Runnable() {
//		                 public void run() {
//
//		                     Toast toast = Toast.makeText(StreamingActivity.this,"Streaming..",Toast.LENGTH_LONG);
//		                     toast.setGravity(Gravity.TOP|Gravity.RIGHT, 0, 0);
//		                     toast.show();
//		                }
//		            });
//					Looper.prepare();
//					
//					
//					// get orientation
//					frameInfo.orientation = convertStringToByteArray(""
//							+ overlay.getCurrentOrientation());
//					new SavePhotoTask(overlay).execute(tmpData, frameInfo.uid, frameInfo.lat, frameInfo.lon,
//							frameInfo.orientation);
//					if (isFPS) {
//						tmpCamera.startPreview();
//						mPreview.inPreview = true;
//
//					}
//
//					long start_time = System.currentTimeMillis();
//					while (System.currentTimeMillis() < start_time + 5000 && isFPS) {
//						Thread.yield();
//					}
//
//					if (isFPS) {
//						
//						Log.d(TAG, "Picture taken");
//						mCamera.takePicture(null, null, photoCallback);
//					}
//				}
//			})).start();
//
//		}
//
//	};

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
		String lati = "" + (location.getLatitude());
		String longi = "" + (location.getLongitude());
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

