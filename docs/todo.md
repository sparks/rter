# ToDos
## Server
### FixMe
* Append logfile don't overwrite (this should be a flag option)
* Hijack bug (don't crash the whole server when websockets are hijacked)
* Twitter items
	* Single tweet type should be lower case
	* Set start-stop time and other fields for tweets
* Finish videoStreamingV1 item implemenation (use tiles from alex)
* Caching services
	* Update tag-cloud on the fly
	* Update Map-Dir on the fly (Userdirection)
	* Update comments on the fly (it's a chat system!)
* Image upload, especially for user submitted generic content. 
* Validate stream doesn already exist in the server before handing out a token
* Don't return target heading if it's too old
* Minimize call to db for auth

### Short Term Wish List
* Tile view improvements
	* Live badge
	* Content type badge
	* Viewed/Dealt with check box
	* Generaly make more informative
* Stream of live items
* Provide a lock on the mobile user control
* Tour of the UI (tool tips, animations ...)
* Structural
	* Callback/Hooks for CRUD instead of mass of switch statements
	* Store foreign item tokens in the DB such as the video server token?
* Self rebooting server

### Long Term Wish List
* Timeline with scrubber
* Logging
* More filters for the sidebar
	* by trust levels
	* By content type
	* By associated "task" or response group?
* Standard auto tagging format (maybe TOML/JSON/XML spec)

### Even More Features!
* Auto window tiling mechanism for multi screen. Breakout UI/UX elements accross windows
* Free roam mode for phone

## VideoServer
### Bugs and Feature Requests
* Check permissions on directories

## Mobile

### Both
* More debug information
	* Compas and Lat/Lng readout
	* Show uncertainty metrics onscreen

### iOS
* Stream doesnt work a second time without logout
* Orientation bug

### Android
* Use new Auth system
* Location/orientation issues (can be reproduced easily on Nehil's tablet)