package teaonly.projects.droidipcam;

import teaonly.projects.droidipcam.NanoHTTPD;

import java.io.*;
import java.util.*;
import android.util.Log;

public class StreamingServer extends NanoHTTPD
{
    public static interface OnRequestListen {
        public abstract InputStream onRequest();
        public abstract void requestDone();
    }

    private OnRequestListen myRequestListen = null;
    private File homeDir;
    private Response streamingResponse = null;

    public StreamingServer(int port, String wwwroot) throws IOException {
        super(port, new File( wwwroot ).getAbsoluteFile() );
        homeDir = new File( wwwroot);
    }
    
    public void setOnRequestListen( OnRequestListen orl) {
        myRequestListen = orl;
    }

    public Response serve( String uri, String method, Properties header, Properties parms, Properties files ) {
        Log.d("TEAONLY", method + " '" + uri + "' " );

        if ( uri.equalsIgnoreCase("/live.flv") ) {
            Response res = new Response( HTTP_NOTFOUND, MIME_PLAINTEXT,  "Error 404, file not found." );
            if ( myRequestListen == null) {
                return res;
            } else {


                InputStream ins;
                ins = myRequestListen.onRequest();
                if ( ins == null)
                    return res;
                
                if ( streamingResponse == null ) {
                    Random rnd = new Random();
                    String etag = Integer.toHexString( rnd.nextInt() );

                    res = new Response( HTTP_OK, "video/x-flv", ins);
                    res.addHeader( "Connection", "Keep-alive");
                    res.addHeader( "ETag", etag);
                    res.isStreaming = true; 
                    streamingResponse = res;
                    Log.d("TEAONLY", "Starting streaming server");
                }
                return res;
            }
        } else {
            return serveFile( uri, header, homeDir, true ); 
        }
    }

    public void serveDone( Response r) {
        if ( r.mimeType.equalsIgnoreCase("video/x-flv") && r == streamingResponse) {
            if ( myRequestListen != null) {
                myRequestListen.requestDone();
                streamingResponse = null;
            }    
        }
    }
    
    
}
