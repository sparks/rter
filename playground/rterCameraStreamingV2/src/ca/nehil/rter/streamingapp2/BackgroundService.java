/**
 * 
 */
package ca.nehil.rter.streamingapp2;

import android.app.Service;
import android.content.Intent;
import android.os.IBinder;
import android.widget.Toast;
/**
 * @author nehiljain
 *
 */
public class BackgroundService extends Service {

	/**
	 * 
	 */
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
		Toast.makeText(this, "Bakcground  Service created", Toast.LENGTH_SHORT).show();
		Thread initBkgdThread = new Thread(new Runnable(){
			public void run(){
				
				for(int i=0; i<= 240;i++ ){
					Toast.makeText(BackgroundService.this, "Thread Running", Toast.LENGTH_SHORT).show();
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
		Toast.makeText(this, "Service destroyed", Toast.LENGTH_SHORT).show();
	}

}
