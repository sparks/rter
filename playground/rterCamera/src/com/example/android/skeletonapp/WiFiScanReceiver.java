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
			{"00:1F:45:E3:E2:99","43.64419","-79.38684"},
			{"20:B3:99:94:3B:08","43.64394","-79.38672"},
			{"20:B3:99:A1:4D:88","43.64378","-79.38619"},
			{"20:B3:99:A1:49:D8","43.64364","-79.38701"}
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
    Toast.makeText(wifiDemo, message, Toast.LENGTH_LONG).show();
    Log.d(TAG, message);
    String WifiID = bestSignal.BSSID.toUpperCase();
	for(int i = 0;i <= wifiMap.length-1; i++){
		if(wifiMap[i][0].matches(WifiID) && bestSignal.SSID.matches("Canada 3.0"))
		{	
			wifiDemo.internalLat = wifiMap[i][1] ;
			wifiDemo.internalLng = wifiMap[i][2] ;
			wifiDemo.changeLocation();
			Toast.makeText(wifiDemo, message, Toast.LENGTH_LONG).show();
			Toast.makeText(wifiDemo, message, Toast.LENGTH_LONG).show();
			
		    Log.d(TAG, "onReceive() message: " + message);
			Log.d("WIFI", "WIFI: lat= "+wifiMap[i][1] +" and lng= "+wifiMap[i][2]);
		}else{
			String internalLat="43.643886";
			String internalLng="-79.386885";
		}
	}
	
    
  }

}