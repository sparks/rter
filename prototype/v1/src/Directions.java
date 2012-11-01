import processing.core.*;

import panoia.*;

import java.io.*;
import java.util.*;

import org.w3c.dom.Document;
import org.w3c.dom.*;

import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.DocumentBuilder;

import java.net.URL;
import java.io.InputStream;

public class Directions {

	PApplet p;
	Pano pano;

	//End points
	LatLng start, stop;

	//Directions
	int current_step;
	LatLng[] steps;

	//Caching
	int cachedex;
	int cachelim;
	int cache_size;

	PImage[][][] tileCache;
	PImage[][] threeFoldCache;
	PanoData[] datacache;
	PanoPov[] povcache;

	//Playback
	int runer = -1;
	boolean anim;
	int anim_runer = -1;

	public Directions(PApplet p, Pano pano, LatLng start, LatLng stop, int cache_size) {
		this.p = p;

		this.start = start;
		this.stop = stop;
		this.pano = pano;

		current_step = 0;
		steps = buildRoute(getXML(start, stop));

		this.cache_size = cache_size;
		cachedex = 0;
		cachelim = 0;

		tileCache = new PImage[cache_size][7][3];
		threeFoldCache = new PImage[cache_size][3];
		datacache = new PanoData[cache_size];
		povcache = new PanoPov[cache_size];
	}

	public void reset() {
		pano.setPosition(start);
		float bearing = (float)pano.getPosition().getInitialBearing(steps[(current_step+1)%steps.length]);
		pano.setPov(new PanoPov(0, bearing, 0));
	}

	public void draw() {
		if(anim) {
			anim_runer++;
			if(anim_runer >= cache_size || anim_runer >= cachelim) {
				anim = false;
			} else {
				pano.tileCache = tileCache[anim_runer];
				pano.threeFoldCache = threeFoldCache[anim_runer];
				pano.data = datacache[anim_runer];
				pano.pov = povcache[anim_runer];
			}
		}
	}

	public boolean step() {
		if(current_step+1 == steps.length) {
			System.out.println("Final step reached");
			return false;
		} else {
			float dist = (float)pano.getPosition().getDistance(steps[(current_step+1)%steps.length]);

			if(dist < 20) {
				current_step = (current_step+1)%steps.length;
				System.out.println("Now Step"+ current_step);

				pano.setPosition(steps[current_step]);
			} else {
				pano.jump();
			}

			if(current_step+1 != steps.length) {
				float bearing = (float)pano.getPosition().getInitialBearing(steps[(current_step+1)%steps.length]);
				pano.setPov(new PanoPov(0, bearing, 0));
			}

			tileCache[cachedex] = pano.tileCache;
			threeFoldCache[cachedex] = pano.threeFoldCache;
			datacache[cachedex] = pano.data;
			povcache[cachedex] = pano.pov;

			cachedex++;
			cachelim++;
			if(cachelim >= cache_size) cachelim = cache_size;
			if(cachedex >= cache_size) cachedex = 0;

			return true;
		}
	}

	public void keyPressed(char key) {
		if(key == 'r') {
			runer++;
			if(runer >= cache_size || runer >= cachelim) runer = 0;

			pano.tileCache = tileCache[runer];
			pano.threeFoldCache = threeFoldCache[runer];
			pano.data = datacache[runer];
			pano.pov = povcache[runer];
		} else if(key == 'g') {
			reset();
			
			cachedex = 0;

			(new Thread(new Runnable() {
				public void run() {
					for(int i = 0;i < cache_size;i++) {
						System.out.println("Caching step "+i+" of "+cache_size);
						step();
						p.draw();
						try {
							Thread.sleep(100);
						} catch(Exception e) {

						}
					}
					System.out.println("Cache Done");
				}
			})).start();

		} else if(key == 's') {
			step();
		} else if( key == 'a') {
			anim_runer = 0;
			anim = true;
		}
	}

	public Document getXML(LatLng start, LatLng stop) {
		try	{
			DocumentBuilderFactory docBuilderFactory = DocumentBuilderFactory.newInstance();
			DocumentBuilder docBuilder = docBuilderFactory.newDocumentBuilder();

			String urlString = "http://maps.googleapis.com/maps/api/directions/xml?origin="+start.toUrlValue()+"&destination="+stop.toUrlValue()+"&sensor=false&units=metric&mode=driving";
			// if(pano.apikey != null && !pano.apikey.equals("")) urlString += "&key="+pano.apikey;

			URL url = new URL(urlString);

			InputStream stream = url.openStream();
			Document xml = docBuilder.parse(stream);
			
			return xml;
		} catch (Exception e) {
			System.err.println(e);
			System.err.println("Panoia: Error in getXML...");
		}
		
		return null;
	}

	public LatLng[] buildRoute(Document xml) {
		try {
			//Parse Location Information and Copyright
			NodeList latTags = xml.getElementsByTagName("lat");
			NodeList lngTags = xml.getElementsByTagName("lng");

			LatLng[] steps = new LatLng[(latTags.getLength()-4)]; //Super hack

			for(int i = 0;i < steps.length;i++) {
				float lat = Float.parseFloat(latTags.item(i).getFirstChild().getNodeValue());
				float lng = Float.parseFloat(lngTags.item(i).getFirstChild().getNodeValue());

				steps[i] = new LatLng(lat, lng);
			}

			return steps;
		} catch(Exception e) {
			//Probably no such street view, so ignore
			System.err.println(e);
			System.err.println("Panoia: Error no such streetview or malformed streetview...");
			return null;
		}

	}

}