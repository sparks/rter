package teaonly.projects.droidipcam;

import java.net.InetAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;

import android.content.Context;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
import android.util.Log;
import android.telephony.TelephonyManager;
import android.telephony.*;


public class NetInfoAdapter {
	
    private static Map<String,String> infoMap = new HashMap<String, String>();
	private static Map<Integer,String> phoneType = new HashMap<Integer, String>();
	private static Map<Integer,String> networkType = new HashMap<Integer, String>();

    static {
 		// Initialise some mappings
    	phoneType.put(0,"None");
    	phoneType.put(1,"GSM");
    	phoneType.put(2,"CDMA");
    	
    	networkType.put(0,"Unknown");
    	networkType.put(1,"GPRS");
    	networkType.put(2,"EDGE");
    	networkType.put(3,"UMTS");
    	networkType.put(4,"CDMA");
    	networkType.put(5,"EVDO_0");
    	networkType.put(6,"EVDO_A");
    	networkType.put(7,"1xRTT");
    	networkType.put(8,"HSDPA");
    	networkType.put(9,"HSUPA");
    	networkType.put(10,"HSPA");
    	networkType.put(11,"IDEN");

        infoMap.put("Cell", "false");
        infoMap.put("Mobile", "false");
        infoMap.put("Wi-Fi", "false");
    }
	
    public static void Update(Context context) {
		// Initialise the network information mapping
        infoMap.put("Cell", "false");            
        TelephonyManager tm = (TelephonyManager) context.getSystemService(Context.TELEPHONY_SERVICE);
		if( tm != null ) {
            infoMap.put("Cell", "true");
            if ( tm.getCellLocation() != null) { 
                infoMap.put("Cell location", tm.getCellLocation().toString());	
            }
			infoMap.put("Cell type", getPhoneType(tm.getPhoneType()));
		}

    	// Find out if we're connected to a network
        infoMap.put("Mobile", "false");
        infoMap.put("Wi-Fi", "false");
        ConnectivityManager cm = (ConnectivityManager) context.getSystemService(Context.CONNECTIVITY_SERVICE);
        NetworkInfo ni = (NetworkInfo) cm.getActiveNetworkInfo();         
       if ( ni != null && ni.isConnected() ) {
            WifiManager wifi = (WifiManager) context.getSystemService(Context.WIFI_SERVICE);
            NetworkInterface intf = getInternetInterface();
            infoMap.put("IP", getIPAddress(intf));
            String type = (String) ni.getTypeName();
            if ( type.equalsIgnoreCase("mobile") ) {
                infoMap.put("Mobile", "true");
                infoMap.put("Mobile type", getNetworkType(tm.getNetworkType())); 
                infoMap.put("Signal", "Good!");
            } else if (  wifi.isWifiEnabled() ) {
                WifiInfo wi = wifi.getConnectionInfo();
                infoMap.put("Wi-Fi", "true");
                infoMap.put("SSID",  wi.getSSID());
                infoMap.put("Signal", "Good!");
            }
        } 
    }

    public static String getInfo(String key) {                                                              
        return infoMap.containsKey(key)? infoMap.get(key): "";
    }

    private static String getPhoneType(Integer key) {
        if( phoneType.containsKey(key) ) {
            return phoneType.get(key);
        } else {
            return "unknown";
        }
    }
                   
    private static String getNetworkType(Integer key) {     
        if( networkType.containsKey(key) ) {
            return networkType.get(key);
        } else {
            return "unknown";
        }
    }

    private static String getIPAddress( NetworkInterface intf) {
        String result = ""; 
        for( Enumeration<InetAddress> enumIpAddr = intf.getInetAddresses(); enumIpAddr.hasMoreElements();) {
            InetAddress inetAddress = enumIpAddr.nextElement();
            result = inetAddress.getHostAddress();
        }
        return result;
    }
    
    private static NetworkInterface getInternetInterface() {
        try {
            for (Enumeration<NetworkInterface> en = NetworkInterface.getNetworkInterfaces(); en.hasMoreElements();) {
                NetworkInterface intf = en.nextElement();
                if( ! intf.equals(NetworkInterface.getByName("lo"))) {
                    return intf;        
                }                   
            }   
        } catch (SocketException ex) {              
        }   
        return null;
    }    

}
