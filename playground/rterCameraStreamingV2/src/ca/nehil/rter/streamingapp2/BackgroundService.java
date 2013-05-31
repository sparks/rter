/**
 * 
 */
package ca.nehil.rter.streamingapp2;

import java.io.IOException;
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
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Handler;
import android.os.IBinder;
import android.util.Log;
import android.widget.Toast;
/**
 * @author nehiljain
 *
 */
public class BackgroundService extends Service {

	/**
	 * 
	 */
	
	private static final String TAG = "Background Service";
	private static final String SERVER_URL = "http://rter.cim.mcgill.ca";
//	private static final String SERVER_URL = "http://132.206.74.145:8000";
	private SharedPreferences cookies;
	private SharedPreferences.Editor prefEditor;
	private String setRterResource;
	private String setRterCredentials;
	private String setItemID;
	private String setRterSignature;
	private String setRterValidUntil;
	
	
	
	
	
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
		Toast.makeText(this, 
			TAG +	" destroyed", Toast.LENGTH_SHORT).show();
	}
	
	
	
	class PuSensorsFeed extends Thread {
	    private Handler handler = null;
	    private NotificationRunnable runnable = null;
	    
	    public PuSensorsFeed(Handler handler, NotificationRunnable runnable) {
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
	    
	    @Override
	    public void run() {
	    	showMessage("Closing feed thread started");
	    	
	    	
	    	JSONObject jsonObjSend = new JSONObject();
			
			
			
			Date date = new Date();
			SimpleDateFormat dateFormatUTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'");
			dateFormatUTC.setTimeZone(TimeZone.getTimeZone("UTC"));
			String formattedDate = dateFormatUTC.format(date);
			Log.i(TAG, "The Stop Timestamp "+formattedDate);
	
			try {
				
				jsonObjSend.put("Live", false);
				jsonObjSend.put("StoptTime", formattedDate);
				
				// Output the JSON object we're sending to Logcat:
				Log.i(TAG,"Body of closefeed json = "+ jsonObjSend.toString(2));
				
				
				int TIMEOUT_MILLISEC = 10000;  // = 10 seconds
				URL url = new URL(SERVER_URL+"/1.0/items/"+setItemID);
				HttpURLConnection httpcon = (HttpURLConnection) url.openConnection();
//				httpcon.setDoOutput(true);
//				httpcon.setRequestProperty("Content-Type", "application/json");
//				httpcon.setRequestProperty("Accept", "application/json");
				httpcon.setRequestProperty("Cookie", setRterCredentials );
				Log.i(TAG,"Cookie being sent" + setRterCredentials);
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
	            case 201:
	               Log.i(TAG,"PUT sensor Feed response = successful");              
	                
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
	}

	
	
	
	
	

}
