package teaonly.projects.droidipcam;

import teaonly.projects.droidipcam.R;

import java.io.*;

import android.app.Activity;
import android.content.Context;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper; 
import android.util.Log;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.SurfaceView;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.Window;
import android.view.WindowManager;
import android.widget.Button;
import android.widget.TextView;
import android.widget.RadioButton;
import android.widget.RadioGroup;
import android.widget.Toast;

public class MainActivity extends Activity {
	private static final int MENU_EXIT = 0xCC882201;
  
    StreamingLoop cameraLoop;
    StreamingLoop httpLoop;

    NativeAgent nativeAgt;
    CameraView myCamView;
    StreamingServer strServer;
    
    TextView myMessage;
    Button btnStart;
    RadioGroup resRadios;

    boolean inServer = false;
    boolean inStreaming = false;
    int targetWid = 320;
    int targetHei = 240;

    final String checkingFile = "/sdcard/ipcam/myvideo.mp4";
    final String resourceDirectory = "/sdcard/ipcam";

    // memory object for encoder
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);        


        requestWindowFeature(Window.FEATURE_NO_TITLE);
        Window win = getWindow();
		win.addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);		
        win.setFlags(WindowManager.LayoutParams.FLAG_FULLSCREEN,
				WindowManager.LayoutParams.FLAG_FULLSCREEN); 
        
        setContentView(R.layout.main);
        
        //nehils addition
        LocationManager lm = (LocationManager) getSystemService(Context.LOCATION_SERVICE);
        LocationListener ll = new mylocationlistener();
        lm.requestLocationUpdates(LocationManager.GPS_PROVIDER, 0, 0, ll);
    }
    
    private class mylocationlistener implements LocationListener {
    	
        @Override
        public void onLocationChanged(Location location) {
            if (location != null) {
            	int lat = (int) (location.getLatitude());
                int lng = (int) (location.getLongitude());
                Log.d("LOCATION CHANGED", location.getLatitude() + "");
                Log.d("LOCATION CHANGED", location.getLongitude() + "");
                Toast.makeText(MainActivity.this,
                		location.getLatitude() + "" + location.getLongitude(),
                		Toast.LENGTH_LONG).show();
            }
        }
        @Override
        public void onProviderDisabled(String provider) {
        }
        @Override
        public void onProviderEnabled(String provider) {
        	Toast.makeText(this, "Enabled new provider " + provider,
        	        Toast.LENGTH_SHORT).show();
        	        
        }
        @Override
        public void onStatusChanged(String provider, int status, Bundle extras) {
        }
		
    }
    
	@Override
	public boolean onCreateOptionsMenu(Menu m){
    	m.add(0, MENU_EXIT, 0, "Exit");
    	return true;
    }
	
	@Override
    public boolean onOptionsItemSelected(MenuItem i){
    	switch(i.getItemId()){
		    case MENU_EXIT:
                finish();
                System.exit(0);
		    	return true;	    	
		    default:
		    	return false;
		}
    }

    @Override
    public void onDestroy(){
    	super.onDestroy();
    }

    @Override
    public void onStart(){
    	super.onStart();
        setup();
    }

    @Override
    public void onResume(){
    	super.onResume();
    }

    @Override
    public void onPause(){    	
    	super.onPause();
        finish();
        System.exit(0);
    }
    
    private void clearResource() {
        String[] str ={"rm", "-r", resourceDirectory};

        try { 
            Process ps = Runtime.getRuntime().exec(str);
            try {
                ps.waitFor();
            } catch (InterruptedException e) {
                e.printStackTrace();
            } 
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private void buildResource() {
        String[] str ={"mkdir", resourceDirectory};

        try { 
            Process ps = Runtime.getRuntime().exec(str);
            try {
                ps.waitFor();
            } catch (InterruptedException e) {
                e.printStackTrace();
            } 
        
            copyResourceFile(R.raw.index, resourceDirectory + "/index.html"  );
            copyResourceFile(R.raw.style, resourceDirectory + "/style.css"  );
            copyResourceFile(R.raw.player, resourceDirectory + "/player.js"  );
            copyResourceFile(R.raw.player_object, resourceDirectory + "/player_object.swf"  );
            copyResourceFile(R.raw.player_controler, resourceDirectory + "/player_controler.swf"  ); 
        }
        catch (IOException e) {
            e.printStackTrace();
        }
    }

    private void setup() {
        clearResource();
        buildResource(); 

        NativeAgent.LoadLibraries();
        nativeAgt = new NativeAgent();
        cameraLoop = new StreamingLoop("teaonly.projects");
        httpLoop = new StreamingLoop("teaonly.http");

    	myCamView = (CameraView)findViewById(R.id.surface_overlay);
        SurfaceView sv = (SurfaceView)findViewById(R.id.surface_camera);
    	myCamView.SetupCamera(sv);
       
        myMessage = (TextView)findViewById(R.id.label_msg);

        btnStart = (Button)findViewById(R.id.btn_start);
        btnStart.setOnClickListener(startAction);
        btnStart.setEnabled(true);
		
		/*
        Button btnTest = (Button)findViewById(R.id.btn_test);
        btnTest.setOnClickListener(testAction);
        btnTest.setVisibility(View.INVISIBLE);
		*/

        RadioButton rb = (RadioButton)findViewById(R.id.res_low);
        rb.setOnClickListener(low_res_listener);
        rb = (RadioButton)findViewById(R.id.res_medium);
        rb.setOnClickListener(medium_res_listener);
        rb = (RadioButton)findViewById(R.id.res_high);
        rb.setOnClickListener(high_res_listener);

        resRadios = (RadioGroup)findViewById(R.id.resolution);

        View  v = (View)findViewById(R.id.layout_setup);
        v.setVisibility(View.VISIBLE);
    }
    
    private void startServer() {
        inServer = true;
        btnStart.setText( getString(R.string.action_stop) );
        btnStart.setEnabled(true);    
        NetInfoAdapter.Update(this);
        myMessage.setText( getString(R.string.msg_prepare_ok) + " http://" + NetInfoAdapter.getInfo("IP")  + ":8080" );

        try {
            strServer = new StreamingServer(8080, resourceDirectory); 
            strServer.setOnRequestListen(streamingRequest);
        } catch( IOException e ) {
            e.printStackTrace();
            showToast(this, "Can't start http server..");
        }
    }

    private void stopServer() {
       inServer = false;
       btnStart.setText( getString(R.string.action_start) );
       btnStart.setEnabled(true);    
       myMessage.setText( getString(R.string.msg_idle));
       if ( strServer != null) {
            strServer.stop();
            strServer = null;
       }
    }

    private boolean startStreaming() {
        if ( inStreaming == true)
            return false;
        
        cameraLoop.InitLoop();
        httpLoop.InitLoop();
        nativeAgt.NativeStartStreamingMedia(cameraLoop.getReceiverFileDescriptor() , httpLoop.getSenderFileDescriptor());

        myCamView.PrepareMedia(targetWid, targetHei);
        boolean ret = myCamView.StartStreaming(cameraLoop.getSenderFileDescriptor());
        if ( ret == false) {
            return false;
        } 
        
        new Handler(Looper.getMainLooper()).post(new Runnable() { 
            public void run() { 
                showToast(MainActivity.this, getString(R.string.msg_streaming));
                btnStart.setEnabled(false);
            } 
        });

        inStreaming = true;
        return true;
    }

    private void stopStreaming() {
        if ( inStreaming == false)
            return;
        inStreaming = false;

        myCamView.StopMedia(); 
        httpLoop.ReleaseLoop();
        cameraLoop.ReleaseLoop();
        
        nativeAgt.NativeStopStreamingMedia();
        new Handler(Looper.getMainLooper()).post(new Runnable() { 
            public void run() { 
                btnStart.setEnabled(true);
            } 
        });
    }

    private void doAction() {
         if ( inServer == false) {
            myCamView.PrepareMedia(targetWid, targetHei);
            boolean ret = myCamView.StartRecording(checkingFile);
           if ( ret ) {
               btnStart.setEnabled(false);
               resRadios.setEnabled(false);
               myMessage.setText( getString(R.string.msg_prepare_waiting));
               new Handler().postDelayed(new Runnable() { 
                    public void run() { 
                        myCamView.StopMedia();
                        if ( NativeAgent.NativeCheckMedia(targetWid, targetHei, checkingFile) ) {
                            startServer();    
                        } else {
                            btnStart.setEnabled(true);
                            resRadios.setEnabled(false);
                            showToast(MainActivity.this, getString(R.string.msg_prepare_error2));
                        }
                    } 
                }, 2000); // 2 seconds to release 
            } else {
                showToast(this, getString(R.string.msg_prepare_error1));
            }
        } else {
            stopServer();
        }
    
    }

    private void copyResourceFile(int rid, String targetFile) throws IOException {
        InputStream fin = ((Context)this).getResources().openRawResource(rid);
        FileOutputStream fos = new FileOutputStream(targetFile);  

        int     length;
        byte[] buffer = new byte[1024*32]; 
        while( (length = fin.read(buffer)) != -1){
            fos.write(buffer,0,length);
        }
        fin.close();
        fos.close();
    }

    private void showToast(Context context, String message) { 
        // create the view
        LayoutInflater vi = (LayoutInflater)context.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
        View view = vi.inflate(R.layout.message_toast, null);

        // set the text in the view
        TextView tv = (TextView)view.findViewById(R.id.message);
        tv.setText(message);

        // show the toast
        Toast toast = new Toast(context);
        toast.setView(view);
        toast.setDuration(Toast.LENGTH_SHORT);
        toast.show();
    }   

    StreamingServer.OnRequestListen streamingRequest = new StreamingServer.OnRequestListen() {
        @Override
        public InputStream onRequest() {
            Log.d("TEAONLY", "Request live streaming...");
            if ( startStreaming() == false)
                return null;
            try {
                InputStream ins = httpLoop.getInputStream(); 
                return ins;
            } catch (IOException e) {
                e.printStackTrace();
                Log.d("TEAONLY", "call httpLoop.getInputStream() error");
                stopStreaming();              
            } 
            Log.d("TEAONLY", "Return a null response to request");
            return null;
        }
        
        @Override 
        public void requestDone() {
            Log.d("TEAONLY", "Request live streaming is Done!");
            stopStreaming();     
        }
    };

    private OnClickListener startAction = new OnClickListener() {
        @Override
        public void onClick(View v) {
            doAction();
        }
    };

    private OnClickListener testAction = new OnClickListener() {
        @Override
        public void onClick(View v) {
            // do some debug testing here.
        }
    };
    
    private OnClickListener low_res_listener = new OnClickListener() {
        @Override
        public void onClick(View v) {
            targetWid = 320;
            targetHei = 240;
        }
    };
    private OnClickListener medium_res_listener = new OnClickListener() {
        @Override
        public void onClick(View v) {
            targetWid = 640;
            targetHei = 480;
        }
    };
    private OnClickListener high_res_listener = new OnClickListener() {
        @Override
        public void onClick(View v) {
            targetWid = 1280;
            targetHei = 720;
        }
    };
}
