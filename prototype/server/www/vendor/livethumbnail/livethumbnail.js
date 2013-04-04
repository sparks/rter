// A live video thumbnail widget that continuously fetches images.
//
// usage:
//   <div live-thumbnail="config"></div>
//
//   config: {
//      showtitle: bool expr,   // display title, (disabled)
//      autoplay: bool expr,    // toggle auto-update, (enabled)
//      selectable: bool expr,  // selectability, (disabled)
//      skimmable: bool expr,   // enable skim mode, (enabled)
//      clickable: bool expr,   // emit clicked signal, (disabled)
//      interval: int expr,     // seconds between updates
//      width: int expression   // width of single thumbnail in pixels (160)
//      height: int expression  // height of single thumbnail in pixels (90)
//      debuglive: bool expr,   // assume stream starts 'now' regarless of start time
//      video: video object     // the video object
//  }
//
//  Requires a video object with the following parameters:
//    .title
//    .thumbnailUrl
//    .StartTime [ISO time string in UTC]
//    .EndTime (optional) [ISO time string in UTC]
//
//  Emitted events: catch with scope.$on(name, function(event, args){})
//
//  Event Name     Args           Description
//  ---------------------------------------------------------------
//  selected       video object   fired when thumbnail is selected
//  deselected     video object   fired when thumbnail is unselected
//  clicked        video object   fired when thumbnail is clicked
//  playing        video object   fired when live playback (re)starts
//  paused         video object   fired when live playback pauses (not skimming)
//  skimming       video object   fired when playback is paused for skimming
//  eos            video object   fired when stream ends
//
//  TODO
//  - test download timer tuning to match server timing (avoiding wrong EOS detection)
//  - image url template as config option?
//  - display live metadata (start/end/current time, duration, skim position)
//  - touch events
//
angular.module('tsunamijs.livethumbnail', [])

.directive('liveThumbnail', ['$timeout', '$log', function($timeout, $log) {
    'use strict';
    return {
        restrict: 'A',
        replace: true,
        templateUrl: 'vendor/livethumbnail/livethumbnail.html',
        controller: 'LiveThumbnailCtrl',
        link: function(scope, element, attrs, controller) {

            var DisplayState = {
                playing : 0,
                skimming : 1,
                paused : 2
            };

            var StreamState = {
                live : 0,
                eos : 1
            };

            // default config
            scope.c = {
                    showtitle: false,
                    autoplay: true,
                    selectable: false,
                    skimmable: true,
                    clickable: false,
                    interval: 2,
                    width: 160,
                    height: 90,
                    debuglive: false,
                    video: {}
            };

            // cache Angular elements to query positions on mousemove
            var outer = angular.element(element.children()[1]);
            var thumb = angular.element(outer.children()[0]);

            // local variables private to each instance of this widget
            var dstate = DisplayState.skimming;
            var sstate = StreamState.live;
            var timeoutId = 0;
            var playhead = 0;
            // image streaming
            var timeRange = 0;
            var timeRangeError = 0;
            var failedCount = 0;
            var updateInProgress = false;
            var looping = false;
            var baseUrl = "";
            var lastUrl = "";
            var img = new Image();
            var startTimeUTC = 0;
            var endTimeUTC = 0;

            // live variables used in DOM update
            scope.skimmerWidth = 0;
            scope.selected = false;
            scope.bgimageurl = "";
            scope.status = "";

            scope.isPlaying = function() {
                return dstate === DisplayState.playing;
            };

            scope.isSkimming = function() {
                return dstate === DisplayState.skimming;
            };

            scope.isPaused = function() {
                return dstate === DisplayState.paused;
            };

            scope.isLive = function() {
                return sstate === StreamState.live;
            };

            scope.isEOS = function() {
                return sstate === StreamState.eos;
            };

            function setDisplayState(s) {
                dstate = s;
                switch (s) {
                    case DisplayState.skimming:
                        scope.status = "Skimming";
                        break;
                    case DisplayState.playing:
                        switch (sstate) {
                        case StreamState.live:
                            scope.status = "Live";
                            break;
                        case StreamState.eos:
                            scope.status = "Playback";
                            break;
                        }
                        break;
                    case DisplayState.paused:
                        switch (sstate) {
                        case StreamState.live:
                            scope.status = "Paused";
                            break;
                        case StreamState.eos:
                            scope.status = "Playback";
                            break;
                        }
                        break;
                }
            }

            function setStreamState(s) {
                sstate = s;
                setDisplayState(dstate);
            }

            function logConfig() {
                $log.info("Live Thumbnail Config:" +
                    " showtitle=" + scope.c.showtitle +
                    " selectable=" + scope.c.selectable +
                    " clickable=" + scope.c.clickable +
                    " autoplay=" + scope.c.autoplay +
                    " skimmable=" + scope.c.skimmable +
                    " size=" + scope.c.width + "x" + scope.c.height +
                    " debuglive=" + scope.c.debuglive +
                    " interval=" + scope.c.interval +
                    " video=" + scope.c.video.title );
            }

            function resetState() {
                dstate = DisplayState.skimming;
                sstate = StreamState.live;
                $timeout.cancel(timeoutId);
                timeoutId = 0;
                playhead = 0;
                timeRange = 0;
                timeRangeError = 0;
                failedCount = 0;
                updateInProgress = false;
                looping = false;
                baseUrl = "";
                lastUrl = "";
                img = new Image();
                startTimeUTC = 0;
                endTimeUTC = 0;

                scope.bgimageurl = "";
                scope.skimmerWidth = 0;
                scope.selected = false;
                scope.bgimageurl = "";
                scope.status = "";
            }

            function updateConfig(newVal, oldVal) {
                if (newVal === undefined) { return; }

                // save some state
                var prevVideo = scope.c.video;
                var prevSelected = scope.selected;

                if (scope.selected) {
                    scope.$emit("deselected", prevVideo);
                }

                // reset live state
                resetState();

                // re-read config
                var extConfig;
                if (attrs.liveThumbnail) {
                    extConfig = scope.$eval(attrs.liveThumbnail);
                } else {
                    extConfig = {};
                }
                angular.extend(scope.c, extConfig);

                // update data from video object
                baseUrl = scope.c.video.thumbnailUrl.slice(0, scope.c.video.thumbnailUrl.lastIndexOf('/') + 1);

                // init time values
                var now = new Date();
                var nowutc = new Date(now.toUTCString());
                var utc = nowutc.getTime();
                var s, e;
                if (scope.c.video.StartTime !== undefined) {
                    if (scope.c.debuglive) {
                        startTimeUTC = utc;
                    } else {
                        s = new Date(scope.c.video.StartTime);
                        startTimeUTC = s.getTime();
                    }
                } else {
                    startTimeUTC = utc;
                }

                if (scope.c.video.StopTime !== undefined) {
                    if (scope.c.debuglive) {
                        s = new Date(scope.c.video.StopTime);
                        e = new Date(scope.c.video.StartTime);
                        endTimeUTC = s.getTime() - e.getTime() + startTimeUTC;
                    } else {
                        s = new Date(scope.c.video.StopTime);
                        endTimeUTC = s.getTime();
                    }
                } else {
                    endTimeUTC = 0;
                }

                // infer live state from values of start and end time
                if (scope.c.debuglive) {
                    setStreamState(StreamState.live);
                } else {
                    setStreamState(startTimeUTC < endTimeUTC ? StreamState.eos : StreamState.live);
                }

                if (endTimeUTC !== 0 && !scope.c.debuglive) {
                    timeRange = ~~((endTimeUTC - startTimeUTC)/1000);
                } else {
                    timeRange = ~~((nowutc - startTimeUTC)/1000);
                }

                // select widget if it was selected before
                scope.selected = scope.c.selectable && prevSelected;

                // set initial image to the first (timeRange=0) or last (timeRange=1)
                // available image
                lastUrl = makeUrl(timeRange, 0); //HACK: Length sometimes load wrong, last is safer for now

                // play state
                if (scope.c.autoplay) {
                    play();
                } else {
                    play(); //HACK: Load initial image
                    if (scope.c.skimmable) { skim(); } else { pause(); }
                }

                // start playloop to fetch new images in the background
                if (!looping) { playLoop(); }
            }

            // initially call update right away
            updateConfig();

            // watch configuration changes
            scope.$watch(attrs.liveThumbnail, updateConfig, true);

            function toggleSelect() {
                if (scope.c.selectable) {
                    // set new state
                    if (scope.selected) { scope.selected = false; }
                    else {
                        if (!scope.isSkimming()) { scope.selected = true; }
                    }
                    // emit signals
                    if (scope.selected) {
                        scope.$emit("selected", scope.c.video);
                    } else {
                        scope.$emit("deselected", scope.c.video);
                    }
                }
            }

            scope.mouseClicked = function() {

                // forward click event
                if (scope.c.clickable) { scope.$emit("clicked", scope.c.video); }

                // toggle selection
                toggleSelect();

                if(scope.c.clickable) { //HACK: Click tmp click solution
                    // toggle play state
                    if (scope.c.skimmable) {
                        if (!scope.isSkimming()) { skim(); } else { play(); }
                    } else {
                        if (!scope.isPaused()) { pause(); } else { play(); }
                    }
                }
            };

            // start live play
            function play() {
                if (sstate === StreamState.eos) {
                    scope.$emit("eos", scope.c.video);
                } else {
                    scope.$emit("playing", scope.c.video);
                }
                setDisplayState(DisplayState.playing);

                if (lastUrl !== "") { scope.bgimageurl = lastUrl; }
                else { scheduleUpdate(); }
            }

            // skim through history
            function skim() {
                setDisplayState(DisplayState.skimming);
                scope.$emit("skimming", scope.c.video);
            }

            function pause() {
                if (sstate !== StreamState.eos) {
                    setDisplayState(DisplayState.paused);
                    scope.$emit("paused", scope.c.video);
                }
            }

            // timer callback to adjust playhead and continue auto-skimming
            function playLoop() {

                // stop live playback at EOS
                if (scope.isEOS()) {
                    looping = false;
                    return;
                }

                looping = true;

                // save the timeoutId for canceling
                timeoutId = $timeout(function() {
                        scheduleUpdate(); // update DOM asynchronously
                        playLoop();  // schedule another update
                    }, scope.c.interval * 1000);
            }

            // manual skimming called from mouse callback
            scope.skipTo = function(e) {
                if (!scope.c.skimmable || !scope.isSkimming()) { return; }

                // clamp and normalise to available number of pixels
                // requires jQuery for .position()
                playhead = Math.max(0, Math.min(~~(e.offsetX - thumb.position().left), scope.c.width));

                // propagate to DOM
                scope.skimmerWidth = playhead / scope.c.width * 100;

                // pick the image and display it
                var reachTime = timeRange * playhead / scope.c.width;
                fetchImageAsync(makeUrl(timeRange, playhead / scope.c.width),
                    function() {
                        // set the fetched image as new background on success
                        scope.bgimageurl = img.src;
                    },
                    function() {
                        timeRange = reachTime; //HACK: Truncate down if length is off w/r/thumbs
                    });
            };

            // listen on DOM destroy (removal) event, and cancel timer
            // to prevent updating after the DOM element was removed.
            scope.$on("$destroy",function() {
                $timeout.cancel(timeoutId);
            });

            // construct the image url and pad with leading zeros
            // http://stackoverflow.com/questions/2998784/how-to-output-integers-with-leading-zeros-in-javascript
            function makeUrl(range, pos) {
                return baseUrl +
                       ('000000000' + (~~(range*pos/scope.c.interval + 1))).substr(-9) +
                       ".jpg";
            }

            // fetch a new image from server
            function fetchImageAsync(url, onloadfunc, onerrorfunc) {
                // browser will load the image asynchronously and call functions
                img = new Image();
                img.onload = onloadfunc;
                img.onerror = onerrorfunc;
                img.src = url;
            }

            // fetch newest images during live playback
            function scheduleUpdate() {


                // server is slow
                if (updateInProgress) {
                    // cancel the download
                    $log.info("LiveThumb: slow server, cancelling download");
                    img.src = "";
                }
                updateInProgress = true;

                // calculate new thumbnail id
                // http://praveenlobo.com/techblog/how-to-convert-javascript-local-date-to-utc-and-utc-to-local-date/
                var now = new Date(); // this is in local time
                var nowutc = new Date(now.toUTCString());
                var utc = nowutc.getTime();

                // check for EOS state by end time (if end time was set)
                if (endTimeUTC > 0 && endTimeUTC <= utc) {
                    setStreamState(StreamState.eos);
                    scope.$emit("eos", scope.c.video);
                    return;
                }

                // update time range for live streams
                timeRange = ~~((utc - startTimeUTC)/1000) - timeRangeError; // sub and convert to sec

                var url = makeUrl(timeRange, 1);
                var s = new Date(startTimeUTC);

                // fetch new thumbnail images in background to be able to check
                // for EOS state (after several consecutive fetch errors we assume EOS)
                fetchImageAsync(url,
                    // on success
                    function() {
                        updateInProgress = false;
                        failedCount = 0;

                        // set the fetched image as new background
                        if (!scope.isSkimming() && !scope.isPaused()) {
                            scope.bgimageurl = img.src;
                        }
                        lastUrl = img.src;
                    },
                    // on error
                    function() {
                        updateInProgress = false;
                        failedCount++;

                        // accumulate range on error (don't reset on success)
                        timeRangeError += scope.c.interval;

                        // also limit range to avoid skimming beyond end
                        // note: will be reset by next call to scheduleUpdate()
                        timeRange -= scope.c.interval;
                    });

                // could not fetch multiple thumbnails in a row (assume stream has ended)
                if (failedCount > 5) {
                    // store apparent end time
                    endTimeUTC = startTimeUTC + timeRange - timeRangeError;
                    setStreamState(StreamState.eos);
                    scope.$emit("eos", scope.c.video);
                }

            }
        }
    };
}])
.controller('LiveThumbnailCtrl', ['$scope', '$log', function($scope, $log) {

}]);
