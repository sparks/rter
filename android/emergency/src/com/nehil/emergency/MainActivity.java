package com.nehil.emergency;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.Calendar;

import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ByteArrayEntity;
import org.apache.http.impl.client.BasicResponseHandler;
import org.apache.http.impl.client.DefaultHttpClient;
import org.json.JSONException;
import org.json.JSONObject;

import android.app.Activity;
import android.content.Context;
import android.location.Criteria;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.os.Bundle;
import android.os.Handler;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

public class MainActivity extends Activity implements LocationListener {

	private Handler m_handler1;

	private LocationManager locationManager;
	private Button button;
	private String provider;
	private TextView latituteField;
	private TextView longitudeField;
	@Override
	  public void onCreate(Bundle savedInstanceState) {
	    super.onCreate(savedInstanceState);
	    setContentView(R.layout.activity_main);
	    m_handler1 = new Handler();
	    
	    latituteField = (TextView) findViewById(R.id.TextView02);
	    longitudeField = (TextView) findViewById(R.id.TextView04);

	    // Get the location manager
	    locationManager = (LocationManager) getSystemService(Context.LOCATION_SERVICE);
	    // Define the criteria how to select the locatioin provider -> use
	    // default
	    Criteria criteria = new Criteria();
	    provider = locationManager.getBestProvider(criteria, false);
	    Location location = locationManager.getLastKnownLocation(provider);

	    // Initialize the location fields
	    if (location != null) {
	      System.out.println("Provider " + provider + " has been selected.");
	      onLocationChanged(location);
	    } else {
	      latituteField.setText("Location not available");
	      longitudeField.setText("Location not available");
	    }
	    button = (Button) findViewById(R.id.button1);
	    
	  }
	
	
	
	
	
	
	
	
	  /* Request updates at startup */
	  @Override
	  protected void onResume() {
	    super.onResume();
	    locationManager.requestLocationUpdates(provider, 400, 1, this);
	  }

	  /* Remove the locationlistener updates when Activity is paused */
	  @Override
	  protected void onPause() {
	    super.onPause();
	    locationManager.removeUpdates(this);
	  }

	 
	  public void onLocationChanged(Location location) {
	    float lati = (float) (location.getLatitude());
	    float longi = (float) (location.getLongitude());
	    latituteField.setText(String.valueOf(lati));
	    longitudeField.setText(String.valueOf(longi));
	    try {
			this.sendServerData(lati, longi);
		} catch (ClientProtocolException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (JSONException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	  }

	  
	  public void onStatusChanged(String provider, int status, Bundle extras) {
	    // TODO Auto-generated method stub

	  }

	  
	  public void onProviderEnabled(String provider) {
	    Toast.makeText(this, "Enabled new provider " + provider,
	        Toast.LENGTH_SHORT).show();

	  }

	 
	  public void onProviderDisabled(String provider) {
	    Toast.makeText(this, "Disabled provider " + provider,
	        Toast.LENGTH_SHORT).show();
	  }
	  
	  
	  
	  
	  
//	  public void sendServerData(final float lati,final float longi) throws JSONException, ClientProtocolException, IOException {
//		  Log.d("Nehil","Came inside sendServerData");
//		  Runnable runnable = new Runnable() {
//		      
//		    public void run() {
//		    	  Log.d("Nehil","Came inside run");
//		    DefaultHttpClient httpClient = new DefaultHttpClient();
//		    ResponseHandler <String> resonseHandler = new BasicResponseHandler();
//		    HttpPost postMethod = new HttpPost("http://e-caffeine.net/nehil_sandbox/emer/post.php");
//
//		    
//		    Calendar ci = Calendar.getInstance();
//
//		    String CiDateTime = "" + ci.get(Calendar.YEAR) + "-" + 
//		         (ci.get(Calendar.MONTH) + 1) + "-" +
//		         ci.get(Calendar.DAY_OF_MONTH) + " " +
//		         ci.get(Calendar.HOUR) + ":" +
//		         ci.get(Calendar.MINUTE) +  ":" +
//		         ci.get(Calendar.SECOND);
//		    
//		   
//		    JSONObject json = new JSONObject();
//		    try {
//				json.put("latitude",lati);
//			} catch (JSONException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//		    try {
//				json.put("longitude",longi);
//			} catch (JSONException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//		    try {
//				json.put("time",CiDateTime);
//			} catch (JSONException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//		    Log.d("Nehil","json :"+ json.toString());
//		    
//		    postMethod.setHeader( "Content-Type", "application/json" );
//     
//		    try {
//				postMethod.setEntity(new ByteArrayEntity(json.toString().getBytes("UTF8")));
//			} catch (UnsupportedEncodingException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//		    String response;
//			try {
//				response = httpClient.execute(postMethod,resonseHandler);
//			} catch (ClientProtocolException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			} catch (IOException e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//		    
//		    
//		      }
//		    };
//		    }

}
