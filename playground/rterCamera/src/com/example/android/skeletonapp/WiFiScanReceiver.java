package com.example.android.skeletonapp;

import java.util.List;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.net.wifi.ScanResult;
import android.net.wifi.WifiManager;
import android.util.Log;
import android.widget.Toast;


public class WiFiScanReceiver extends BroadcastReceiver {
  private static final String TAG = "WiFiScanReceiver";
  
  CameraPreviewActivity wifiDemo;
  
  private String[][] wifiMap= {
			{"00:1f:45:f3:1e:11","lat","lng"}	
  };
  public WiFiScanReceiver(CameraPreviewActivity wifiDemo) {
    super();
    this.wifiDemo = wifiDemo;
  }

  @Override
  public void onReceive(Context c, Intent intent) {
    List<ScanResult> results = wifiDemo.myWifiManager.getScanResults();
    ScanResult bestSignal = null;
    for (ScanResult result : results) {
      if (bestSignal == null
          || WifiManager.compareSignalLevel(bestSignal.level, result.level) < 0)
        bestSignal = result;
    }

    String message = String.format("%s networks found. %s is the strongest. %s with access mac ",
        results.size(), bestSignal.SSID, bestSignal.BSSID);
    
    Log.d(TAG, message);

	for(int i = 0;i <= wifiMap.length-1; i++){
		if(wifiMap[i][0].matches(bestSignal.BSSID))
		{
			wifiDemo.internalLat = wifiMap[i][1] ;
			wifiDemo.internalLng = wifiMap[i][2] ;
			Log.d("WIFI", "WIFI: lat= "+wifiMap[i][1] +" and lng= "+wifiMap[i][2]);
		}
	}
	
    Toast.makeText(wifiDemo, message, Toast.LENGTH_LONG).show();
    Log.d(TAG, "onReceive() message: " + message);
  }

}