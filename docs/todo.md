# ToDos
## Server
### FixMe
* Detect server is down or unreachable
* Socket reconnect in JS (us the close fn)
* MySQL seems to timeout eventually and return driver bad conn. We should provide to reconnect without rebooting the server.
* Firefox has a click/drag bug. Current dragging causes a click. It should be fixed soon. Currently have ugly workaround with timeout.
* Bootstrap / angular ui doesn't destroy controllers in closed dialogs, this leave my websockets open after closing dialog, currently hacking the bootstrap
* Dragging live tiles is very power hungry and lags
* marker flickr continues
* Maintenance for expired live feeds
* Send alex rter server instructions
* Lagging tab switch (due to opening sockets etc?)
* New a map resize/fix (currently can switch to map mode to resize)
* Naming of taxonomy/tag/tabs
* Image upload, especially for user submitted generic content. 
* Validate stream doesn already exist in the server before handing out a token
* Don't return target heading if it's too old
* Minimize call to db for auth
* Note: GZIP + Websocket + Hijack = fail

### Short Term Wish List
* Perma link to content (location)
* Permissions on tags
* Distance bluring
* Sexy graphs
* Anonymous submissions
* How will we do different universes/instances?
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
* Send messages to the phone user
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
* Port needs to remain constant throughout stream
* Orientation bug

### Android
* Use new Auth system
* Location/orientation issues (can be reproduced easily on Nehil's tablet)