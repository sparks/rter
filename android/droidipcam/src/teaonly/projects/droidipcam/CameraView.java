package teaonly.projects.droidipcam;

import android.content.Context;
import android.hardware.Camera;
import android.media.AudioManager;
import android.media.CamcorderProfile;
import android.media.MediaRecorder;
import android.util.AttributeSet;
import android.util.Log;
import android.view.MotionEvent;
import android.view.SurfaceHolder;
import android.view.SurfaceView;
import android.view.View;
import java.io.IOException;
import java.io.File;
import java.io.FileOutputStream;
import java.io.FileDescriptor;
import java.lang.System;
import java.lang.Thread;
import java.nio.ByteBuffer;
import java.nio.IntBuffer;

public class CameraView extends View implements SurfaceHolder.Callback, View.OnTouchListener{

    private AudioManager mAudioManager = null; 
    private Camera myCamera = null;
    private MediaRecorder myMediaRecorder = null;
    private SurfaceHolder myCamSHolder;
    private SurfaceView	myCameraSView;

    public CameraView(Context c, AttributeSet attr){
        super(c, attr);
        
        mAudioManager = (AudioManager)c.getSystemService(Context.AUDIO_SERVICE);
        mAudioManager.setStreamMute(AudioManager.STREAM_SYSTEM, true); 
    }

    public void SetupCamera(SurfaceView sv){    	
    	myCameraSView = sv;
    	myCamSHolder = myCameraSView.getHolder();
    	myCamSHolder.addCallback(this);
    	myCamSHolder.setType(SurfaceHolder.SURFACE_TYPE_PUSH_BUFFERS);
        
        myCamera = Camera.open();
        /*
        Camera.Parameters p = myCamera.getParameters();
        myCamera.setParameters(p);
        */

        setOnTouchListener(this);
    }
    
    public void PrepareMedia(int wid, int hei) {
        myMediaRecorder =  new MediaRecorder();
        myCamera.stopPreview();
        myCamera.unlock();
        
        myMediaRecorder.setCamera(myCamera);
        myMediaRecorder.setAudioSource(MediaRecorder.AudioSource.CAMCORDER);
        myMediaRecorder.setVideoSource(MediaRecorder.VideoSource.CAMERA);
	    
        CamcorderProfile targetProfile = CamcorderProfile.get(CamcorderProfile.QUALITY_LOW);
        targetProfile.videoFrameWidth = wid;
        targetProfile.videoFrameHeight = hei;
        targetProfile.videoFrameRate = 25;
        targetProfile.videoBitRate = 512*1024;
        targetProfile.videoCodec = MediaRecorder.VideoEncoder.H264;
        targetProfile.audioCodec = MediaRecorder.AudioEncoder.AMR_NB;
        targetProfile.fileFormat = MediaRecorder.OutputFormat.MPEG_4;
        myMediaRecorder.setProfile(targetProfile);
    }
   
    private boolean realyStart() {
        
        myMediaRecorder.setPreviewDisplay(myCamSHolder.getSurface());
        try {
        	myMediaRecorder.prepare();
	    } catch (IllegalStateException e) {
	        releaseMediaRecorder();	
	        Log.d("TEAONLY", "JAVA:  camera prepare illegal error");
            return false;
	    } catch (IOException e) {
	        releaseMediaRecorder();	    
	        Log.d("TEAONLY", "JAVA:  camera prepare io error");
            return false;
	    }
	    
        try {
            myMediaRecorder.start();
        } catch( Exception e) {
            releaseMediaRecorder();
	        Log.d("TEAONLY", "JAVA:  camera start error");
            return false;
        }

        return true;
    }

    public boolean StartStreaming(FileDescriptor targetFd) {
        myMediaRecorder.setOutputFile(targetFd);
        myMediaRecorder.setMaxDuration(9600000); 	// Set max duration 4 hours
        //myMediaRecorder.setMaxFileSize(1600000000); // Set max file size 16G
        myMediaRecorder.setOnInfoListener(streamingEventHandler);
        return realyStart();
    }

    public boolean StartRecording(String targetFile) {
        myMediaRecorder.setOutputFile(targetFile);
                
        return realyStart();
    }
    
    public void StopMedia() {
        myMediaRecorder.stop();
        releaseMediaRecorder();        
    }

    private void releaseMediaRecorder(){
        if (myMediaRecorder != null) {
        	myMediaRecorder.reset();   // clear recorder configuration
        	myMediaRecorder.release(); // release the recorder object
        	myMediaRecorder = null;
            myCamera.lock();           // lock camera for later use
            myCamera.startPreview();
        }
        myMediaRecorder = null;
    }

     
    private MediaRecorder.OnInfoListener streamingEventHandler = new MediaRecorder.OnInfoListener() {
        @Override
        public void onInfo(MediaRecorder mr, int what, int extra) {
            Log.d("TEAONLY", "MediaRecorder event = " + what);    
        }
    };

    @Override
    public void surfaceChanged(SurfaceHolder sh, int format, int w, int h){
    	if ( myCamera != null && myMediaRecorder == null) {
            myCamera.stopPreview();
            try {
                myCamera.setPreviewDisplay(sh);
            } catch ( Exception ex) {
                ex.printStackTrace(); 
            }
            myCamera.startPreview();
        }
    }
    
	@Override
    public void surfaceCreated(SurfaceHolder sh){
    }
    
	@Override
    public void surfaceDestroyed(SurfaceHolder sh){
    }

    @Override
    public boolean onTouch(View v, MotionEvent evt) {
        return true;        
    }

}
