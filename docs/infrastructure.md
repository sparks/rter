# Infrastucture Document

This document is a sketch of the information flow and technologies being used or proposed for use in the rtER project.

## System Structure

                   ----------------------------       ---------------
                   | Immersive Command Center |       | ISAS System |
                   |                          |       |             |
                   ----------------------------       ---------------
             (selected media/content) |  (2)           (3) | (I/O user content: sound files)
                                      |                    |
    -------------------    (1)     -------------------------------------   (4)   --------------------------
    | VOST Web Client | ---------  | rtER Server | Spatialized User    | ------- | Video Streaming Server |
    |                 | (content)  |             | Meta-Content Server | (URIs)  |                        |
    -------------------            -------------------------------------         --------------------------
                                            | (5)       (meta-data) | (6)       (7) | (video stream)
                                            |                       |               |
                                            |                      ------------------- 
                                            ---------------------- | rtER Mobile App | 
                                                                   |                 | 
                                                                   ------------------- 

## Details and Technologies

### Immersive Command Center
CAVE environment with top down projected maps and immersive street view (developped for round 1). Includes various visualized/spatialized data.

* Queries for most relevant and highly ranked data from the rtER server via (2).
* (2) not currently implemented like uses HTTP/GET with xml/json. Queries returned include URIs for content.
* [Video](http://vimeo.com/52631497)

### VOST Web Client
Collaborative application for VOST volunteers (developped for round 2). Shows user content which can be collaboratively manipulated discussed and promoted.

 * Queries for new content over (1) via HTTP/GET usually via AJAX, data returned as JSON with URI references to content.
 * Submit content ranking via (1) HTTP/GET, usually AJAX, data JSON.
 * Submit new content via HTTP/POST, mime-multipart for images.
 * Submit other content manipulation via (1) again HTTP/POST via AJAX with JSON.
 * [Video](http://vimeo.com/57946497)


### rtER Server / Spatialized User Meta-Content Server
Generic system the "Spatialized User Meta-Content Server" used to store content submitted by users which is tagged with spatial data. Availble via API to submit, retrieve and modify this data.

Some special tools and feature for the VOST Web client (such as servicing AJAX request for grid layout). Also some special tools for the rtER Mobile App such handling the interactive heading adjustements between VOST and field users.

 * (1) Content queries and VOST content
 * (2) ....
 * (3) ISAS will use our system for user submitted spatialized audio clips.
 * (5) Desired heading information relayed from VOST volunteers.
 * (6) Heading and location information from phone. Video stream information.
 * (7) Video stream.

### rtER Mobile App
Mobile application to stream video. Allows users to be directed where to film by VOST from the VOST Mobile App.

* (5) Desired heading information relayed from VOST volunteers.
* (6) Heading and location information from phone. Video stream information.
* (7) Video stream.
* [Video](http://vimeo.com/57946497)

### Video Stream Server
?
