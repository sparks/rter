package ca.nehil.rter.streamingapp2;

import java.io.File;
import java.io.FileOutputStream;
import java.text.SimpleDateFormat;
import java.util.Date;

import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.HttpVersion;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntity;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.entity.mime.content.StringBody;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.BasicHttpParams;
import org.apache.http.params.CoreProtocolPNames;
import org.apache.http.params.HttpParams;
import org.apache.http.util.EntityUtils;

import ca.nehil.rter.streamingapp2.overlay.OverlayController;

import android.annotation.SuppressLint;
import android.os.AsyncTask;
import android.os.Environment;
import android.provider.Settings;
import android.util.Log;
import android.provider.Settings;

class SavePhotoTask extends AsyncTask<byte[], String, String> {

	private DefaultHttpClient mHttpClient;
	private File photo;

	OverlayController overlay;

	public SavePhotoTask(OverlayController overlay) {
		this.overlay = overlay;
	}

	@Override
	protected void onPostExecute(String result) {
		// delete the uploaded picture to free memory
		Log.d("SavePhotoTask", "Upload Complete");
		if (photo.exists()) {
			photo.delete();
		}
		Log.d("SavePhotoTask", "Photo deleted");

	}

	@SuppressLint("ParserError")
	@Override
	protected String doInBackground(byte[]... a) {
		String lat ="", lon = "";
		Log.e("SavePhotoTask", "uid" + new String(a[1]));
		String uid = new String(a[1]);
		if(a[2] != null && a[3] != null){
		lat = new String(a[2]);
		lon = new String(a[3]);
		}
		else{
			Log.e("SavePhotoTask", "ERROR::No Location!!");
		}
		
		String orient = new String(a[4]);

		Log.d("SavePhotoTask", "phone id " + uid + " lat: " + lat + " lon: "
				+ lon + " orient: " + orient);
		HttpParams params = new BasicHttpParams();
		params.setParameter(CoreProtocolPNames.PROTOCOL_VERSION,
				HttpVersion.HTTP_1_1);
		mHttpClient = new DefaultHttpClient(params);

		// get the address of SD card and check if it exists and makes a folder
		// rter
		String root = (Environment.getExternalStorageDirectory()).toString();
		File rootDir = new File(Environment.getExternalStorageDirectory()
				+ File.separator + "rter" + File.separator);
		rootDir.mkdirs();
		if (root != null) {
			Log.d("SavePhotoTask", "The address of the external storage is "
					+ root);
			// save in SD Card

			String timeStamp = new SimpleDateFormat("_yyyy_MM_dd_hh_mm_ss_SSS")
					.format(new Date());

			photo = new File(rootDir, "Scenephoto" + timeStamp + ".jpg");
			Log.d("SavePhotoTask", "Saving pic" + "Scenephoto" + timeStamp
					+ ".jpg");
			if (photo.exists()) {
				photo.delete();
			}
		} else {
			Log.e("SavePhotoTask",
					"SD card or anyother exernal storage doesnt esxist.");
		}
		try {
			FileOutputStream fos = new FileOutputStream(photo.getPath());

			fos.write(a[0]);
			fos.close();

			HttpPost httppost = new HttpPost(
					"http://rter.cim.mcgill.ca:80/multiup");

			MultipartEntity multipartEntity = new MultipartEntity(
					HttpMultipartMode.BROWSER_COMPATIBLE);
			multipartEntity.addPart("title", new StringBody("rTER"));
			multipartEntity.addPart("phone_id", new StringBody(uid));
			multipartEntity.addPart("lat", new StringBody(lat));
			multipartEntity.addPart("lng", new StringBody(lon));
			multipartEntity.addPart("heading", new StringBody(orient));
			multipartEntity.addPart("image", new FileBody(photo));
			httppost.setEntity(multipartEntity);
			Log.e("SavePhotoTask", "Upload executed");
			HttpResponse response = mHttpClient.execute(httppost);

			HttpEntity r_entity = response.getEntity();
			String responseString = EntityUtils.toString(r_entity);
			Log.d("Htttp Response", responseString);

			overlay.setDesiredOrientation(Float.parseFloat(responseString));

		} catch (java.io.IOException e) {
			Log.e("PictureDemo", "Exception in photoCallback", e);
		} catch (Exception e) {
			Log.e("ServerError", e.getLocalizedMessage(), e);
		}

		return (null);
	}

}
