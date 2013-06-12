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

import org.apache.http.Header;
import org.apache.http.HttpResponse;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ByteArrayEntity;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.BasicHttpParams;
import org.apache.http.params.HttpConnectionParams;
import org.apache.http.params.HttpParams;
import org.json.JSONException;
import org.json.JSONObject;

import android.animation.Animator;
import android.animation.AnimatorListenerAdapter;
import android.annotation.TargetApi;
import android.app.Activity;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.os.Build;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.KeyEvent;
import android.view.Menu;
import android.view.View;
import android.view.inputmethod.EditorInfo;
import android.widget.EditText;
import android.widget.TextView;

/**
 * Activity which displays a login screen to the user, offering registration as
 * well.
 */
public class GetTokenActivity extends Activity {
	
//	private static final String SERVER_URL = "http://rter.cim.mcgill.ca";
	private static final String SERVER_URL = "http://132.206.74.145:8000";
	private static final String TAG = "GetTokenActivity";
	private String rterCreds=null;
	/**
	 * Keep track of the login task to ensure we can cancel it if requested.
	 */
	private HandshakeTask handshakeTask = null;

	
	private View mLoginStatusView;
	private TextView mLoginStatusMessageView;
	private SharedPreferences cookies;
	private SharedPreferences.Editor prefEditor;
	private String setRterCredentials;
	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);

		setContentView(R.layout.activity_token);
		
		cookies = getSharedPreferences("RterUserCreds", MODE_PRIVATE);
		prefEditor = cookies.edit();
		
		String setUsername = cookies.getString("Username", "not-set");
		String setPassword = cookies.getString("Password", "not-set");
		setRterCredentials = cookies.getString("RterCreds", "not-set");
		Log.d(TAG, "Prefs ==> Username:"+setUsername+" :: Password:" + setPassword +" :: rter cred:" + setRterCredentials);
		// Set up the login form.
		mLoginStatusView = findViewById(R.id.login_status);
		mLoginStatusMessageView = (TextView) findViewById(R.id.login_status_message);
		attemptHandshake();
	}

	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		super.onCreateOptionsMenu(menu);
		getMenuInflater().inflate(R.menu.login, menu);
		return true;
	}

	/**
	 * Attempts to sign in or register the account specified by the login form.
	 * If there are form errors (invalid email, missing fields, etc.), the
	 * errors are presented and no actual login attempt is made.
	 */
	public void attemptHandshake() {
		
			// Show a progress spinner, and kick off a background task to
			// perform the user login attempt.
			mLoginStatusMessageView.setText(R.string.login_progress_signing_in);
			showProgress(true);
			handshakeTask = new HandshakeTask();
			handshakeTask.execute();
		
	}

	/**
	 * Shows the progress UI and hides the login form.
	 */
	@TargetApi(Build.VERSION_CODES.HONEYCOMB_MR2)
	private void showProgress(final boolean show) {
		// On Honeycomb MR2 we have the ViewPropertyAnimator APIs, which allow
		// for very easy animations. If available, use these APIs to fade-in
		// the progress spinner.
		if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.HONEYCOMB_MR2) {
			int shortAnimTime = getResources().getInteger(
					android.R.integer.config_shortAnimTime);

			mLoginStatusView.setVisibility(View.VISIBLE);
			mLoginStatusView.animate().setDuration(shortAnimTime)
					.alpha(show ? 1 : 0)
					.setListener(new AnimatorListenerAdapter() {
						@Override
						public void onAnimationEnd(Animator animation) {
							mLoginStatusView.setVisibility(show ? View.VISIBLE
									: View.GONE);
						}
					});

			
		} else {
			// The ViewPropertyAnimator APIs are not available, so simply show
			// and hide the relevant UI components.
			mLoginStatusView.setVisibility(show ? View.VISIBLE : View.GONE);
			
		}
	}

	/**
	 * Represents an asynchronous login/registration task used to authenticate
	 * the user.
	 */
	public class HandshakeTask extends AsyncTask<Void, Void, Boolean> {
		private static final String TAG = "GetTokenActivity HandshakeTask";
		
		@Override
		protected Boolean doInBackground(Void... params) {
			// TODO: attempt authentication against a network service.
			
			JSONObject jsonObjSend = new JSONObject();
			
			
			
			Date date = new Date();
			SimpleDateFormat dateFormatUTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'");
			dateFormatUTC.setTimeZone(TimeZone.getTimeZone("UTC"));
			String formattedDate = dateFormatUTC.format(date);
			Log.i(TAG, "The Time stamp "+formattedDate);
	
			
			
			try {
				jsonObjSend.put("Type", "streaming-video-v1");
				jsonObjSend.put("Live", true);
				jsonObjSend.put("StartTime", formattedDate);
				jsonObjSend.put("HasGeo", true);
				jsonObjSend.put("HasHeading", true);
				// Output the JSON object we're sending to Logcat:
				Log.i(TAG, jsonObjSend.toString(2));
				Log.i(TAG,"Cookie being sent" + setRterCredentials);
				
				int TIMEOUT_MILLISEC = 10000;  // = 10 seconds
				URL url = new URL(SERVER_URL+"/1.0/items");
				HttpURLConnection httpcon = (HttpURLConnection) url.openConnection();
//				httpcon.setDoOutput(true);
//				httpcon.setRequestProperty("Content-Type", "application/json");
//				httpcon.setRequestProperty("Accept", "application/json");
				httpcon.setRequestProperty("Cookie", setRterCredentials );
				httpcon.setRequestMethod("POST");
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
	                
	                String itemID = jObject.getString("ID");
	                String uploadURI = jObject.getString("UploadURI");
	                JSONObject token = jObject.getJSONObject("Token");
	                String rter_resource = token.getString("rter_resource");
	                String rter_signature = token.getString("rter_signature");
	                String rter_valid_until = token.getString("rter_valid_until");
	                Log.i(TAG,"Response from connection rter_resource : " + rter_resource);
	                Log.i(TAG,"Response from connection rter_signature : " + rter_signature);
	                
	                prefEditor.putString("ID", itemID); 
	                prefEditor.putString("rter_resource", rter_resource);  
					prefEditor.putString("rter_signature", rter_signature); 
					prefEditor.putString("rter_valid_until", rter_valid_until); 
					prefEditor.commit();
	                
	                
	                
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
			
			// TODO: register the new account here.
			return true;
		}

		@Override
		protected void onPostExecute(final Boolean success) {
			handshakeTask = null;
			showProgress(false);

			if (success) {				
				Log.i(TAG, "Calling Intent to Streaming ACtivity");
				Intent intent = new Intent(GetTokenActivity.this, StreamingActivity.class);
		        startActivity(intent);
				
			} else {
				
			}
		}

		@Override
		protected void onCancelled() {
			handshakeTask = null;
			showProgress(false);
		}
	}
}
