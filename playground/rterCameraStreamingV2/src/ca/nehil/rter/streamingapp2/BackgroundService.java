/**
 * 
 */
package ca.nehil.rter.streamingapp2;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.UnsupportedEncodingException;
import java.net.HttpURLConnection;
import java.net.URL;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.TimeZone;

import org.apache.http.client.ClientProtocolException;
import org.json.JSONException;
import org.json.JSONObject;

import ca.nehil.rter.streamingapp2.StreamingActivity.NotificationRunnable;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.location.Criteria;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.os.Bundle;
import android.os.Handler;
import android.os.IBinder;
import android.util.Log;
import android.view.Gravity;
import android.widget.Toast;
/**
 * @author nehiljain
 *
 */
public class BackgroundService extends Service implements 
LocationListener  {

	/**
	 * 
	 */
	
	private int PutLocationTimer = 15000; /* Updating the User location, heading and orientation every 4 secs. */
	private static final String TAG = "Background Service";
	private static final String SERVER_URL = "http://rter.cim.mcgill.ca";
//	private static final String SERVER_URL = "http://132.206.74.145:8000";
	
	private SharedPreferences cookies;
	private SharedPreferences.Editor prefEditor;
	private String setRterResource;
	private String setRterCredentials;
	private String setUsername;
	private LocationManager locationManager;
	private String provider;
	
	private float lati;
	private float longi;
	
	
	public BackgroundService() {
		// TODO Auto-generated constructor stub
	}

	/* (non-Javadoc)
	 * @see android.app.Service#onBind(android.content.Intent)
	 */
	@Override
	public IBinder onBind(Intent arg0) {
		// TODO Auto-generated method stub
		return null;
	}
	
	
	@Override
	public void onCreate(){
		super.onCreate();
		Toast.makeText(this, TAG + " created", Toast.LENGTH_SHORT).show();
		
		locationManager = (LocationManager) getSystemService(Context.LOCATION_SERVICE);
		if(!locationManager.isProviderEnabled(LocationManager.GPS_PROVIDER));
		{
			Log.e(TAG, "GPS not available");
		}
		
		Criteria criteria = new Criteria();
		provider = locationManager.getBestProvider(criteria, false);
		
		Log.d(TAG, "Requesting location");
		
		locationManager.requestLocationUpdates(provider, 2, 1, this);
		if(provider != null){
			Location location = locationManager.getLastKnownLocation(provider);
			// Initialize the location fields
			if (location != null) {
				System.out.println("Provider " + provider
						+ " has been selected. and location "+ location);
				onLocationChanged(location);
			} else {
				Toast toast = Toast.makeText(this, "Location not available", Toast.LENGTH_LONG);
				toast.setGravity(Gravity.TOP, 0, 0);
				toast.show();
				lati = (float) (45.505958f);
				longi = (float)(-73.576254f);
				Log.d(TAG, "Location not available");
			}
		}
		
		
		cookies = getSharedPreferences("RterUserCreds", MODE_PRIVATE);
		prefEditor = cookies.edit();
		setUsername = cookies.getString("Username", "not-set");
		setRterCredentials = cookies.getString("RterCreds", "not-set");
		
		Log.d(TAG, "Prefs ==> rter_resource:"+setRterCredentials+" :: Username:" + setUsername );
		
		
		Thread initBkgdThread = new Thread(new Runnable(){
			public void run(){
				
				for(int i=0; i<= 240;i++ ){
					Log.d("Service", "Running inside thread");
					try {
						Thread.sleep(2000);
					} catch (InterruptedException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					}
				}
				
				Toast.makeText(BackgroundService.this, "Thread Running", Toast.LENGTH_SHORT).show();
			}
		});
		initBkgdThread.start();
		
	}
	public void onDestroy(){
		super.onDestroy();
		locationManager.removeUpdates(this);
		Toast.makeText(this, 
			TAG +	" destroyed", Toast.LENGTH_SHORT).show();
	}
	
	
	
	class PuSensorsFeed extends Thread {
		class PutSensorsFeed extends Thread {
		    private Handler handler = null;
		    private NotificationRunnable runnable = null;
		    
		    public PutSensorsFeed(Handler handler, NotificationRunnable runnable) {
		        this.handler = handler;
		        this.runnable = runnable;
		        this.handler.post(this.runnable);
		    }
		    
		    /**
		    * Show UI notification.
		    * @param message
		    */
		    private void showMessage(String message) {
		        this.runnable.setMessage(message);
		        this.handler.post(this.runnable);
		    }
		    
		    
		    private void postHeading(){
		    	JSONObject jsonObjSend = new JSONObject();								
				
				try {
					
					float lat = lati;
					float lng = longi;
					float heading = 0;
					
					jsonObjSend.put("Username", setUsername);
					jsonObjSend.put("Lat", lat );
					jsonObjSend.put("Lng", lng );
					jsonObjSend.put("Heading", heading);
					
					// Output the JSON object we're sending to Logcat:
					Log.i(TAG,"postHeading()::Body of update location service feed json = "+ jsonObjSend.toString(2));				
					
					int TIMEOUT_MILLISEC = 1000;  // = 1 seconds
					Log.i(TAG,"postHeading()Put Request being sent" +SERVER_URL+"/1.0/users/"+setUsername);
					URL url = new URL(SERVER_URL+"/1.0/users/"+setUsername);
					HttpURLConnection httpcon = (HttpURLConnection) url.openConnection();
					if(setRterCredentials.equalsIgnoreCase("not-set")){
						Log.d(TAG, "Process flow is wrong");
					}else{
						httpcon.setRequestProperty("Cookie", setRterCredentials );
					}
					
					httpcon.setRequestProperty("Content-Type", "application/json");
					httpcon.setRequestMethod("PUT");
					httpcon.setConnectTimeout(TIMEOUT_MILLISEC);
					httpcon.setReadTimeout(TIMEOUT_MILLISEC);
					httpcon.connect();
					byte[] outputBytes = jsonObjSend.toString().getBytes("UTF-8");
					OutputStream os = httpcon.getOutputStream();
					os.write(outputBytes);

					os.close();
					
					int status = httpcon.getResponseCode();
					Log.i(TAG,"Status of response " + status);
					switch (status) {
		            case 200:
		          
		               Log.i(TAG,"PUT sensor Feed response = successful");              
		               BufferedReader br = new BufferedReader(new InputStreamReader(httpcon.getInputStream()));
		               StringBuilder sb = new StringBuilder();
		               String line;
		               while ((line = br.readLine()) != null) {
		                    sb.append(line+"\n");
		               }
		               String result = sb.toString();
		               br.close();
		                
		               JSONObject jObject = new JSONObject(result);
		               Log.i(TAG,"Response from connection " + jObject.toString(2));
		                    
					}
					
				} catch (JSONException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				} catch (UnsupportedEncodingException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				} catch (ClientProtocolException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				} catch (IOException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				}
		    }
		    
		 
		    @Override
		    public void run() {
		    	Log.d(TAG," Update heading and location thread started");    	
		    	while(true) {
		    	    long millis = System.currentTimeMillis();
		    	    this.postHeading();
		    	    
		    	    	    	    
		    	    try {
						Thread.sleep((PutLocationTimer - millis % 1000));
					} catch (InterruptedException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					}
		    	}	
					
		    }
		}
	}
	@Override
	public void onLocationChanged(Location location) {
		// TODO Auto-generated method stub
		lati =  (float) (location.getLatitude());
		longi =  (float) (location.getLongitude());
		Log.d(TAG, "Location Changed with lat"+lati+" and lng"+longi);
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



	
