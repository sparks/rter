-App opens (done!)
-Login screen (mostly done, Dan)
-User clicks login:
    -Sent to server for auth (done Dan)
    -Gets cookie back and saves it (done Dan)
-App changes to streaming screen (done Stepan)
-User click start:
    -A new item is created on the server via POST (done Dan)
        -The item should have a Type, a StartTime and a StopTime set before the StartTime (mostly done Dan)
        -If a available this should also include Lat/Lng and Heading (Cameron or Dan?)
    -The returned JSON includes the following fields which are saved: ID, Upload URI, Token (Dan?)
-The streaming beings (done stepan). The UploadURI saved from before is used (not done Dan)
-OPENGL is initialized (Cameron/Stepan)
-During the streaming the following happens at regular intervals (~every second?) 
    -The item is updated on the server via PUT using the ID saved from before. This update changes the Lat/Lng/Heading info. (Cameron?)
    -The TargetHeading for the logged in user is fetched via GET and saved locally (Dan or Cameron?)
-During the streaming the following happens in real time:
    -Video is seen on screen (done Stepan)
    -The OPENGL updates using the most recent TargetHeading and the current actual heading (Cameron/Stepan)
-The user clicks stop or quits or returns to login
    -Optional null frame sent to video server (Stepan)
    -The item is updated on the server via PUT using the ID save from before. The StopTime is set (Dan)

Edge cases to be handled (once everything else is in place or as needed)
    -Video server token is invalid (video server 401s). The app should stop streaming and notify the user. More advanced would be to notify the user while you request a new token and restart the stream
    -rtER server logs you out (401). Cease streaming, bounce to login screen.
    -Internet lost? Probably bounce to login screen? At the very least you will have to restart the stream with a new video token because tokens are bound to your IP

More stuff:
    -If you can easily get accuracy error estimates, add them to the JSON when you send me Lat/Lng/Heading
    -Can we get debug info on screen:
        -Current compass readings
        -Current error accuracy/error estimates