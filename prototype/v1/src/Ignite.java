import processing.core.*;

import panoia.*;

import de.sciss.net.*;
import java.net.SocketAddress;
import java.io.IOException;

import java.io.*;
import java.util.*;
import java.lang.Math;

public class Ignite extends PApplet implements OSCListener {

	Pano pano;
	PanoPov pov;

	float fov = 270;

	ArrayList<CarAccident> carAccidents;
	ArrayList<BikeAccident> bikeAccidents;

	PFont font;
	int roadrand;

	ScrollingPlot temperaturePlot, humidityPlot;

	boolean osc_enable;
	OSCServer osc_in;

	Directions dir;

	public void setup() {
		size(displayWidth, (int)(displayWidth/(6.5f*fov/360)));
		// size(1024*3, 768);

		background(0);
		smooth();
		frameRate(10);

		font = createFont("Helvetica", 14);
		textFont(font, 14);

		pano = new Pano(this, "AIzaSyDklHrdigHKgVYzrDAvSXaaR6Epx1_cygQ");
		pov = pano.getPov();

		// pano.setPano("3qry8ACTZ8Mw6SQ1UaLNMg");
		// pano.setPosition(new LatLng(45.5110809f, -73.5700496f));
		// pano.setPosition(new LatLng(45.52059937f, -73.58165741f));
		// pano.setPano("717wuQJ5lH4xB3Uw5vs4Pw");
		// pano.setPano("FUhF2Lmri2qq6NZErDpn2Q");
		// pano.setPano("wvuAA91CEZ5hP0afgwp_Wg");
		// pano.setPosition(new LatLng(45.506257f,-73.575718f));

		carAccidents = CarAccident.ParseCsv(this);
		bikeAccidents = BikeAccident.ParseCsv(this);

		temperaturePlot = new ScrollingPlot(this, new PVector(100*width/(1024*3), 100*height/768), "Temp.", -100, 100, color(0, 0, 255), color(255, 0, 0), 50*width/(1024*3));
		humidityPlot = new ScrollingPlot(this, new PVector(300*width/(1024*3), 100*height/768), "Humid.", 0, 100, color(255, 128, 0), color(0, 0, 255), 50*width/(1024*3));

		roadrand = round(random(0, 255));

		/* ---- Direction/Interp ---- */

		dir = new Directions(this, pano, new LatLng(45.506903f, -73.570139f), new LatLng(45.506227f, -73.576083f), 60);
		dir.reset();

		/* ---- OSC ---- */

		try {
			osc_in = OSCServer.newUsing(OSCServer.UDP, 8001);
			osc_in.start();
		} catch(IOException e) {
			System.err.println("Error initializing OSCServer ... exiting");
			System.out.println(e);
			System.exit(1);
		}

		osc_in.addOSCListener(this);
	}

	public void draw() {
		try {
			background(0);

			pushMatrix();
			translate(width/2, 0);
			pano.drawThreeFold(width);
			// pano.drawTiles(fov, width);
			popMatrix();

			float alpha = 255;
			for(CarAccident accident : carAccidents) {
				stroke(0, alpha);
				fill(0, 100, 0, alpha);
				project(accident.latLng, accident.toString(), 500);
			}

			alpha = 255;
			for(BikeAccident accident : bikeAccidents) {
				stroke(0, alpha);
				fill(0, 100, 0, alpha);
				project(accident.latLng, accident.toString(), 500);
			}

			stroke(0);
			fill(0, 100);
			rect(20*width/(1024*3), 60*height/768, 420*width/(1024*3), 120*height/768);

			temperaturePlot.update();
			temperaturePlot.draw();
			
			humidityPlot.update();
			humidityPlot.draw();

			drawPanoLinks();
			drawRoads();

			dir.draw();

			//Mouse Ref Lin
			// stroke(255, 0, 0);
			// line(mouseX, 0, mouseX, height);
		} catch(Exception e) {
			
		}
	}

	public void mousePressed() {
		pov.setHeading(pov.heading()+map(mouseX, 0, width, -fov/2, fov/2));
		pano.setPov(pov);

		if (mouseButton == RIGHT) {
			pano.jump();
			roadrand = round(random(0, 255));
		}
	}

	public void keyPressed() {
		if(key == 'n') {
			pano.jump();
			roadrand = round(random(0, 255));
		} else if(key == 'o') {
			osc_enable = !osc_enable;
			System.out.println("OSC_EN = "+osc_enable);
		} else if(key == '1') {
			pano.setPosition(new LatLng(45.506903f, -73.570139f));
			pano.setPov(new PanoPov(0, 0, 0));
		} else if(key == '2') {
			pano.setPosition(new LatLng(45.506257f, -73.575718f));
			pano.setPov(new PanoPov(0, 0, 0));
		}

		dir.keyPressed(key);
	}

	public void drawPanoLinks() {
		stroke(20);
		strokeWeight(2);
		fill(255, 200);
		PanoLink[] links = pano.getLinks();
		for(int i = 0;i < links.length;i++) {
			int x = pano.headingToPixel(links[i].heading, fov, width);
			// line(x, 0, x, height);
			triangle(x, height-20, x-20, height-7, x+20, height-7);
		}
	}

	public void drawRoads() {
		fill(roadrand, 255-roadrand, 0, 50);
		stroke(255, 50);
		PanoLink[] links = pano.getLinks();
		for(int i = 0;i < links.length;i++) {
			int x = pano.headingToPixel(links[i].heading, fov, width);
			float squish = 6f;
			quad(
				x+width/14/squish, height/2+30, 
				x-width/14/squish, height/2+30, 
				x-width/14, height, 
				x+width/14, height
			);

		}
	}

	public boolean project(LatLng point, int size) {
		return project(point, null, size);
	}

	public boolean project(LatLng point, String desc, int size) {
		int x = pano.headingToPixel((float)pano.getPosition().getInitialBearing(point), fov, width);
		// line(x, 0, x, height);

		double tanFactor = 65;
		// int horizonHeight = 2*height/5;
		int horizonHeight = height/2;
		double dis = pano.getPosition().getDistance(point);
		if(dis > 50) return false;
		if(dis < 5) return false;
		double vert = Math.atan(dis/tanFactor*Math.PI);
		int y = (int)(height-vert/Math.PI*2*horizonHeight);

		ellipse(x, y, (int)(size/dis), (int)(size*(Math.PI/2-vert)/Math.PI/dis));

		if(desc != null) {
			int padding = 10;
			int xOff = 30;
			int yOff = -130;
			int xSize = 200;
			int ySize = 75;
			int x2 = constrain(x+xOff, 0, width-xSize);
			int y2 = constrain(y+yOff, ySize, height);
			fill(0, 100);
			stroke(0);
			line(x, y, x2, y2);
			rect(x2, y2, xSize, ySize);
			fill(255);
			stroke(255);
			text(desc, x2+padding, y2+padding, xSize-padding, ySize-padding);
		}

		return true;
	}

	public synchronized void messageReceived(OSCMessage message, SocketAddress sender, long time) {
		if(!osc_enable) return;
		String raw_addr = message.getName();

		if(raw_addr.startsWith("/")) raw_addr = raw_addr.substring(1);

		String[] addr = raw_addr.split("/");

		try {
			if(addr[0].equals("panoid")) {
				if(addr[1].equals("latlong") && message.getArgCount() == 2) {
					float latitude = ((Number)message.getArg(0)).floatValue();
					float longitude = ((Number)message.getArg(1)).floatValue();
					pano.setPosition(new LatLng(latitude, longitude));
					pano.setPov(new PanoPov(0, 0, 0));
				} else if(addr[1].equals("id") && message.getArgCount() == 1) {
					String id = (String)message.getArg(0);
					pano.setPano(id);
				} else if(addr[1].equals("jumpnow") && message.getArgCount() == 0) {
					dir.animate();
				}
			}
		} catch(Exception e) {
			System.out.println("Invalid OSC message contents");
			return;
		}
	}

	public static void main (String [] args) {
		PApplet.main(new String[] { "Ignite" });
	}
}