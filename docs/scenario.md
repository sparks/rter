# Round 2
## Scenario
In this short scenario we are going to present a scenario based loosely on real events that occured at McGill. We want to use this scenario to help us think together about the features of our system and how they could be used in an emergency. In particular our data filtering system and how it can be used collaboratively.

Last month, a 48-inch water main broke at the top of the hill that descends into downtown Montreal, right through the main McGill campus. There was massive flooding through all of the streets and into many of the buildings on surrounding campus, with black ice underneath, making it practically impassable in many spots. The campus was almost completely surrounded by water making it challenging to find a way off campus.

Jeff left just after the main break, but before any official news had been sent. He could not get off campus easily since the water was flowing too deep and fast through the streets, and had no way of figuring out how far the situation spread without slogging through knee-deep water. There were many students and others standing around with cameras taking pictures and video all around campus (later dozens of videos were posted online) and posting on social media.

For those in the buildings trying to plan their way home, being able to see these images, videos and tweets might have been helpful. But without some way of sorting and filtering it might also simply have been overwhelming. 

We imagine that our rtER would have helped us solve this problem. By allowing information -- video, image or text -- to be visualised in geographic context (on the rtER map) users could quickly understand what was occuring where. Moreover, filtering using the map would have provided a way for those looking for a way of campus to eliminate all superflous information and see only the information coming from one area on their map. This would mean you could quickly confirm or deny if a route was safe. 

Also since we are trying to build our system with collaboration in mind, repetious search would be avoided and time could be saved. For example in this case most people posting pictures and videos were "ambulance chaser". Showing the drama and devastation, not the less interesting safe area. This means that the "useful" information is being burried by the "shock value" videos. Once a person had found a safe passage using our system, maybe that one tweet suggesting a safe route, maybe a picture or a video showing an unflooded area, they could promote that information to higher priority for others to see. As a results the next person wouldn't need to repeat that search. They might instead simply add a confirmation after using the route, or perhaps provide an alternative. In this way we hope our system may make a more efficient use of human resources.

(Possible blurb about interactive video, e.g. direct field users to look away from drama for safe route. But I think it's too specific and should be left out for this scenario)

Of course in this context the data would have been most useful for the McGill comunity, and we are not a group trained responders. However we think it provides a simple illustration of how a group of people in an emergency could be using our system to share and manage information to help acheive a common goal: finding safe routes off campus. We could imagine replacing the McGill members with emergency responders are looking for safe staging areas, clear routes to move injured people in and out, maybe looking for vantage point access the broken water main, etc. Collaborating using information both from the responders and from outside source to answer a common question or understand a situation better.

# Round 3
## What's New

* Mobile App
	* Moved to iPhone to match the reality of emergency responders
	* Now includes real low latency streaming including meta data (Lat, Lng, Heading)
	* Still includes interaction with the web user (can be directed, etc)
	* Users now log in on their phones with a personal account, this is reflected web interface
	* Fixed many UI issues (more informative of state, error messages etc)
* Server infratructure
	* The new videoserver outline by alex
	* New rtER Server
		* RESTful, documented api for handling our data. This means everything that goes in our system is easily accessible in a standardized API format. Could be reused within other frameworks or applications. Similarly we could easily take in data from anywhere
		* Designed around our core features so that queries are fast/natural. Support Taxonomy, Users, Ranking within taxonomies
* Web client
	* Built top to bottom in HTML5 components (angularjs, html5 video, google maps v3 with HTML5 canvas/drawing)
	* User accounts. We can log actions by user. Who did what when. This is important since trust and confidence in sources is key. Understanding who's taking an action or submitting information implies a level of trust. 
	* Everything you submit is an Item. A generic container. As of now an item can be a web page, twitter search/post, a youtube video, a live image stream (old android app), a live video stream (new iphone app).
	* Taxonomy as a first class object
		* Items can be tagged
		* Tags can be automated: by content type, by trust level, be location, etc
		* Each tag has it's own view/workspace. e.g. each tag has ranking for all the items it contains which users can collaboratively work with. 
		* Thus tags are much more powerful tools than normal:
			* Tag with your user name or group name to create a personal or shared workspace
			* Tag with your colleague's user name or group to send to their workspace
			* Enable auto tagging by area, content type to create a real time feed for certain condition
			* Permission can be applied to a tag to prevent external interference in your workspace.
		* Web client now supports more filtering techniques: fuzzy searching, map filtering, filtering by user trust level, filter by user, filter by user role
		* Alpha support for HTML5 live stream video with the same interactive control of the remote user (directing the remote user) using a map.
		* Twitter integration
		* Various UI improvement (new interactive map elements, intelligent tag fields, global alerts, improved drag and drop grid)


## Scenario and Pitch

In round three of development on the rtER project, we began with a proof of concept and finished with an alpha product. Key in this round was building up the system so that it was viable for in the field testing by partner emergency response groups: Quebec City (...) and Red Wing (...). From the get go we knew some of features that we would need. Integration with twitter and YouTube were added. We continued our work improving real-time video streaming from mobile device. However early feedback from both groups guided many new key developments. 

In Red Wing responders were concerned about creating group and individual workspaces so responders at different hierarchical levels could interact without interfering with each other. In Quebec City they asked about having multiple filtered information streams for parallel events of unequal importance. We attacked both these by making Taxonomy a first class citizen. An items - a video stream, a twitter search, a youtube video - can be tagged either by through the web interface or by automated filters configured on the server. Each tag has it's own workspace. This means that each tag has a grid view which can be opened. This grid view can be sorted collaboratively by drag and drop. As a result items can also be sent and shared between tag workspaces by adding and removing tags from items. (List example here? Red Wing make a "city EOC" and "state EOC" tags during a major crisis to share back and forth. Quebec has a general feed via the "quebec all" tag, but during a flood an auto tag for live video in a certain area "video-zone1" is created giving only video for the affected area).

The Quebec city responders who have spent more time in the field with our system made numerous user interface suggestions and bug reports that were brought to on the design. In particular their field work underlined the necessity of low latency video streaming. They experienced difficulty interacting with remote users due to both low frame rate of the prototype and the latency, as much as 10 seconds, between action and response. The alpha prototype brings high frame rate and truly low latency streaming to mobile devices, effectively address these issues. Only in browser playback remains a challenge at very low latency.

Finally we want to ensure going forward that we can easily not only bring data into our system but also make our data available to other services making us part of the emergency response ecosystem. To this end our data is available through a standard RESTful API. This means data aggregations tools can be extended to leverage our filtering tools and our live streaming platform.

## Video Pitch

### 1

The rtER web client is allows collaborative filtering and organization of information. Drag key information up and it's instantly promoted for all your collaborator. However, in Red Wing emergency responders often perform join operation with local and state emergency response groups. They asked us, how multiple groups share and collaborate while loosing control of their workspace. 

In the new alpha release of rtER we tackle this problem with our taxonomy system. Like in a traditional taxonomy system content can be tagged, by users or automatically, providing searchable meta data. However, in rtER each tag is a workspace in which content is independently organized and ranked. 

This makes it easy create independent workspace and transfer information between then. In red wing local responders are tag their content as "redwing-eoc", meanwhile state responders tag their content "minn-state-eoc". Transfering content from the redwing workspace to the state workspace is as easy as tagging it with the "minn-state-eoc" tag. Additionally with auto-tagging content matching certain criteria, geolocation, trustworthiness, etc can be automatically added to active workspaces. For example video stream from within the Red Wing municipality area and streamed from trusted volunteers might be autotagged as "redwing-eoc" such that they automatically appear in th redwing workspace.

### 2

rtER we are making interactive live video a first class citizen. In the rtER web interface, not only can you view live video from mobile users, you can also interactively direct their gaze to points of interest in the field. Emergency responders in Quebec City performed in the field testing with the prototype mobile app. In addition to suggesting many UI improvements, they underscored the need for lower latency to create a more responsive interactions. To this end the new version of our mobile app for iOS provides much higher framerate video delivered to the servers, along with location and heading information, at very low latency. 

While we can now boast subsecond end-to end latency for location and heading data, some challenges still remains for video data. Streaming video arrives at our servers with subsecond latency but no native streaming video players support such low latency. As a result one of goals in the next round is to develop an ultra low latency HTML5 video player.

(Blurb about bandwidth an power limits?)

### 3

RESTful?

## Dump from e-mails

"Try to tell us the story."

"Tell us what you can't do now (without the technology you're building)..."

"It's not technology focused... lead with the use (case)."

"What are you going to show in Chicago June 24-26 (approx) to the folks in the broader US Ignite community?  Are you going to have a prototype that works and how are you getting to that point?"

"What is the simplest question you can answer with this tool that is still interesting?"

"It's easier for us (the developers on a particular project) to understand things at the 30,000 foot level, but talking about all the cool tech doesn't register to people listening... "

"It makes more sense to start with a scenario..." (hey, isn't that what we did? ;-)

for example, "Red Wing has a nuclear power plant 5 miles out of town and everybody's worried about what could happen if the plant melts down... Our system can cut the response time down from 7 minutes to 5..."  (well, that's probably not the way we'd go in a pitch, but it helps give the idea of what Will's suggesting).

What does rtER bring that is really new, vs. existing mapping and information filtering tools like those based on Ushahidi? First, it makes real-time video a first class citizen for emergency response. Second, it offers real-time interaction with the people contributing video content from mobile devices. Last, it provides a novel UI for managing content while collaborating with others to analyze and organize it.

The third point probably needs help from Sparky. For me, the main problem is that I'm not confident enough in my knowledge of how it is different from something like SwiftRiver. It doesn't help that their website demo doesn't work for me at all:

http://swiftly.org

Anyway, text like this could possibly serve as the transition into the demonstration, where we show all three of these things and talk about our progress in making them real, and how they were elicited by requests from real potential users, who know emergency management.


[VISUALS: IMAGES OF FLOOD AT MCGILL + BIGGER FLOODS, maybe under Jeremy talking?]
Here in Quebec, as well as places like Red Wing, Minnesota, floods are a major problem, impacting the lives of thousands of people at a time. Closer to home for our Mozilla Ignite development team, earlier this year our campus here at McGill University was inundated with water, resulting in millions of dollars in damages, when a construction crew broke a 48-inch water main directly connected to a major reservoir on a hill over the main campus. As water coursed through the streets and flooded buildings, information was scattered, stale and hard to come by. Some of the best and most current information came from social media. Videos of the damage and dangerous areas were circulated widely, eventually ending up in media reports.

[VISUALS: SNAPSHOTS, ETC OF WHAT WE DID IN PREVIOUS ROUNDS]
Over the previous two development rounds, we told you about the VOST volunteer organizations who are activated by official emergency response teams, and their difficulty in managing such information in real-time to help both emergency responders and the larger public in a catastrophe. We've also shown early prototypes of the three technology pieces of our Mozilla Ignite system: First, our immersive CAVE environment showing a street view with embedded real-time video and sensor content, designed to support emergency operation centres directly. Second, a web tool for VOST-type volunteers that aggregates information from social media, as other existing tools also attempt to do, but that also promotes live video to a first-class citizen, which requires lower latency and higher bandwidth networks than usually available. Third, we demonstrated an early version of our mobile video application, running on an Android smartphone, that not only sends live video to our web app, but also allows someone observing the video to intelligently direct it to more useful viewpoints without having to resort to a live phone call and verbal commands.

Today, we're going to use a flood scenario to demonstrate our progress, then paint a picture of where we expect to be by June of this year with further development.

Most important, we want to tell you about our inroads with emergency response teams who are beginning to test and evaluate our system for practical use. The feedback we've received from them has driven the bulk of our development this milestone, and has come from two primary sources. Thomas Poirier is a <title> for the province of Quebec, and works out of the Emergency Operation Centre (EOC) in Quebec City. He has promoted interest in our system within <their organization>; some of our team members have visited their EOC to demonstrate our Ignite tools to their response management and technology team members. They have already begun using the prototype system, providing detailed feedback on key changes and additional features that would be required to deploy such a tool in practice. Their technical team has received the go-ahead from their management to explore integrating rtER into their existing EOC tool set.

[VIDEO OF THOMAS SAYING SOMETHING ABOUT HOW THEY WOULD ANTICIPATE USING IT]

[something from, or at least about, Red Wing is important here, as they are in the USA...]

[ACTUAL DEMO OF FEATURES TBD, MAKING SURE TO KEEP TYING IT BACK TO THE FEEDBACK FROM QUEBEC AND MINNESOTA USERS]

As with the previous milestones, our code is publicly available on github, including the server and smartphone applications, for anyone to use. Anyone with an Android phone or tablet can download our app from the website, run it, and see their video streaming live on our website. This is the system that the Quebec and Minnesota emergency responders have been evaluating on their own tablets and phones.

In the final development round running through June, we anticipate...

[Is this too much?]  We would welcome the opportunity to demonstrate our system live to US Ignite participants in Chicago in June. By then, with continued development and testing with the responders in Minnesota and Quebec, we expect that we would have a quite compelling demonstration, likely with members of the participating emergency response organizations providing testimonials, based on their ongoing evaluation, of how the technology would enhance their response capabilities.